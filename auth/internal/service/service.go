package service

import (
	"auth/internal/model"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService interface {
	Create(ctx context.Context, info *model.UserInfo, pwd string) (int64, error)
	Get(ctx context.Context, id int64) (int64, *model.UserInfo, *timestamppb.Timestamp, *timestamppb.Timestamp, error)
	Update(ctx context.Context, id int64, name string, email string) error
	Delete(ctx context.Context, id int64) error
}

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
}

type AccessService interface {
	Check(ctx context.Context, endpointAddress string, accessToken string) error
}
