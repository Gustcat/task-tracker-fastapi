package auth

import (
	"auth/internal/model"
	"auth/internal/utils"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serv) GetRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, []byte(s.tokenConfig.RefreshTokenSecretKey()))
	if err != nil {
		return "", status.Errorf(codes.Aborted, "invalid refresh token")
	}

	newRefreshToken, err := utils.GenerateToken(model.UserToken{
		Name: claims.Name,
		Role: claims.Role,
		ID:   claims.ID,
	},
		[]byte(s.tokenConfig.RefreshTokenSecretKey()),
		s.tokenConfig.RefreshTokenExpiration(),
	)
	if err != nil {
		return "", err
	}

	return newRefreshToken, nil
}
