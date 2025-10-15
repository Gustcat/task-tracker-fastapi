package repository

import (
	"auth/internal/model"
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

type UserRepository interface {
	Create(ctx context.Context, info *model.UserInfo, pwd string) (int64, error)
	Get(ctx context.Context, id int64) (int64, *model.UserInfo, time.Time, sql.NullTime, error)
	Update(ctx context.Context, id int64, name string, email string) error
	Delete(ctx context.Context, id int64) error
	Login(ctx context.Context, username string) (string, *model.UserToken, error)
}

type AccessRepository interface {
	Check(ctx context.Context, role int32, endpoint string) error
}
