package user

import (
	"auth/internal/logger"
	"auth/internal/model"
	"context"
	"go.uber.org/zap"
)

func (s *serv) Create(ctx context.Context, info *model.UserInfo, pwd string) (int64, error) {
	logger.Info("Creating user...",
		zap.String("name", info.Name),
		zap.String("email", info.Email),
		zap.String("password", pwd),
		zap.Int32("role", info.Role))

	id, err := s.userRepository.Create(ctx, info, pwd)
	if err != nil {
		return 0, err
	}

	logger.Info("User created", zap.Int64("id", id))

	return id, nil
}
