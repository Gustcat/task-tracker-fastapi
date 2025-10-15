package access

import (
	"auth/internal/config"
	"auth/internal/service"
	descAccess "auth/pkg/access_v1"
)

type Implementation struct {
	descAccess.UnimplementedAccessV1Server
	accessService service.AccessService
	tokenConfig   config.TokenConfig
}

func NewImplementation(accessService service.AccessService, tokenConfig config.TokenConfig) *Implementation {
	return &Implementation{
		accessService: accessService,
		tokenConfig:   tokenConfig,
	}
}
