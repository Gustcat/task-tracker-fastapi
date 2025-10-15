package auth

import (
	"auth/internal/config"
	"auth/internal/repository"
	"auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	tokenConfig    config.TokenConfig
}

func NewService(userRepository repository.UserRepository, tokenConfig config.TokenConfig) service.AuthService {
	return &serv{
		userRepository: userRepository,
		tokenConfig:    tokenConfig,
	}
}
