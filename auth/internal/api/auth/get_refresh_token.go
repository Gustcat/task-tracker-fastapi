package auth

import (
	descAuth "auth/pkg/auth_v1"
	"context"
)

func (i *Implementation) GetRefreshToken(ctx context.Context, req *descAuth.GetRefreshTokenRequest) (*descAuth.GetRefreshTokenResponse, error) {
	refreshToken, err := i.authService.GetRefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &descAuth.GetRefreshTokenResponse{
		RefreshToken: refreshToken,
	}, nil
}
