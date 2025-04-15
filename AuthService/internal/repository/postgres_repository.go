package repository

import (
	"AuthService/internal/domain"
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Create(ctx context.Context, user *domain.User) error {
	var err error
	for i := 0; i < 3; i++ {
		err = r.db.WithContext(ctx).Create(user).Error
		if err == nil {
			return nil
		}
		time.Sleep(time.Second)
	}
	return err
}

func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var (
		user domain.User
		err  error
	)

	for i := 0; i < 3; i++ {
		err = r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
		if err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
			return &user, err
		}
		time.Sleep(time.Second)
	}

	return &domain.User{}, err
}

func (r *AuthRepository) Save(ctx context.Context, user *domain.User) error {
	var err error
	for i := 0; i < 3; i++ {
		err = r.db.WithContext(ctx).Save(&user).Error
		if err == nil {
			return nil
		}
	}
	return err
}
