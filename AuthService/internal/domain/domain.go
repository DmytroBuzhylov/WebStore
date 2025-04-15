package domain

import (
	"context"
	"time"
)

type User struct {
	ID        string    `json:"ID" gorm:"type:uuid;primaryKey;not null;index"`
	FullName  string    `binding:"required,min=3" json:"userName" gorm:"not null;size:30"`
	Email     string    `json:"email" binding:"required,email" gorm:"unique;not null;index"`
	Password  string    `json:"password" binding:"required,min=8" gorm:"not null"`
	Role      string    `json:"role" gorm:"default:user;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP;not null"`
	Verify    bool      `json:"verify" gorm:"default:false;not null"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	Save(ctx context.Context, user *User) error
}

type VerificationCodeRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
}

type Notifier interface {
	SendVerificationCode(ctx context.Context, email string, code string) error
}
