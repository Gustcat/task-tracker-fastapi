package auth

import (
	"auth/internal/utils"
	"context"
	"errors"
)

func (s *serv) Login(ctx context.Context, username, password string) (string, error) {
	hashedPassword, user, err := s.userRepository.Login(ctx, username)
	if err != nil {
		return "", err
	}

	if !utils.VerifyPassword(hashedPassword, password) {
		return "", errors.New("invalid password")
	}

	refreshToken, err := utils.GenerateToken(*user,
		[]byte(s.tokenConfig.RefreshTokenSecretKey()),
		s.tokenConfig.RefreshTokenExpiration(),
	)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
