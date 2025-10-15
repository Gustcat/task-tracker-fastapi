package auth

import (
	descAuth "auth/pkg/auth_v1"
	"context"
)

func (i *Implementation) GetAccessToken(ctx context.Context, req *descAuth.GetAccessTokenRequest) (*descAuth.GetAccessTokenResponse, error) {
	accessToken, err := i.authService.GetAccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &descAuth.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}
