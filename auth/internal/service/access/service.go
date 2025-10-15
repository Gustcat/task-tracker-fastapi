package access

import (
	"auth/internal/config"
	"auth/internal/repository"
	"auth/internal/service"
)

type serv struct {
	accessRepository repository.AccessRepository
	tokenConfig      config.TokenConfig
}

func NewService(accessRepository repository.AccessRepository, tokenConfig config.TokenConfig) service.AccessService {
	return &serv{
		accessRepository: accessRepository,
		tokenConfig:      tokenConfig,
	}
}
