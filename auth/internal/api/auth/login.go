package auth

import (
	descAuth "auth/pkg/auth_v1"
	"context"
)

func (i *Implementation) Login(ctx context.Context, req *descAuth.LoginRequest) (*descAuth.LoginResponse, error) {
	refreshToken, err := i.authService.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &descAuth.LoginResponse{
		RefreshToken: refreshToken,
	}, nil
}
