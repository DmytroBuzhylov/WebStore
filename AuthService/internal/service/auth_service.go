package service

import (
	"AuthService/internal/domain"
	"AuthService/pkg/config"
	"context"
	"crypto/rand"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"math/big"
	rand2 "math/rand"
	"time"
)

type AuthService struct {
	userRepo         domain.UserRepository
	verificationRepo domain.VerificationCodeRepository
	notifier         domain.Notifier
}

func NewAuthService(ur domain.UserRepository, vr domain.VerificationCodeRepository, ntf domain.Notifier) *AuthService {
	return &AuthService{
		userRepo:         ur,
		verificationRepo: vr,
		notifier:         ntf,
	}
}

func (s *AuthService) createNewUser(ctx context.Context, user *domain.User) error {
	var err error
	user.ID = uuid.New().String()
	user.CreatedAt, user.UpdatedAt = time.Now(), time.Now()
	user.Password, err = hashPassword(user.Password)

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return err
	}

	if err = s.sendVerifyCode(ctx, user.Email); err != nil {
		return err
	}
	return nil
}

type Register interface {
	Register(ctx context.Context, user domain.User) error
}

func (s *AuthService) Register(ctx context.Context, user *domain.User) error {
	findData, err := s.userRepo.FindByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println(err)
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = s.createNewUser(ctx, user)
		if err != nil {
			return err
		}
		return nil
	}

	return s.newPasswordForNonVerifyUser(ctx, findData)
}

func (s *AuthService) Login(ctx context.Context, user *domain.User) error {
	findUser, err := s.userRepo.FindByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if !checkPassword(user.Password, findUser.Password) {
		return errors.New("Invalid password")
	}

	if err = s.sendVerifyCode(ctx, findUser.Email); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Verify(ctx context.Context, email, code string) (string, error) {
	cfg := config.AppConfig
	findUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	redisCode, err := s.verificationRepo.Get(ctx, findUser.Email)
	if err != nil {
		return "", err
	}

	if redisCode != code || code == "" || redisCode == "" {
		return "", errors.New("code invalid")
	}

	accessToken, err := generateToken(&AccessTokenGenerator{}, findUser.ID, findUser.Role, cfg.JWTSecret)
	refreshToken, err := generateToken(&RefreshTokenGenerator{}, findUser.ID, findUser.Role, cfg.JWTSecret)
	if err != nil {
		return "", err
	}

	err = s.verificationRepo.Set(ctx, refreshToken, "active", 7*24*time.Hour)
	if err != nil {
		return "", err
	}

	findUser.Verify = true
	err = s.userRepo.Save(ctx, findUser)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *AuthService) sendVerifyCode(ctx context.Context, email string) error {

	if err := s.deleteLastCode(ctx, email); err != nil {
		return err
	}

	if err := s.sendCode(ctx, email); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) deleteLastCode(ctx context.Context, email string) error {
	err := s.verificationRepo.Del(ctx, email)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) sendCode(ctx context.Context, email string) error {
	newInt := big.NewInt(int64(rand2.Intn(1000000)))
	code, _ := rand.Int(rand.Reader, newInt)

	return s.verificationRepo.Set(ctx, email, code.String(), 15*time.Minute)
}

func (s *AuthService) updatePassword(ctx context.Context, user *domain.User) error {
	var err error
	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return err
	}

	return s.userRepo.Save(ctx, user)
}

func (s *AuthService) newPasswordForNonVerifyUser(ctx context.Context, user *domain.User) error {
	err := s.updatePassword(ctx, user)
	if err != nil {
		return err
	}

	err = s.sendVerifyCode(ctx, user.Email)
	if err != nil {
		return err
	}
	return nil
}
