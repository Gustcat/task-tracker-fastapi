package auth

import (
	"auth/internal/service"
	descAuth "auth/pkg/auth_v1"
)

type Implementation struct {
	descAuth.UnimplementedAuthV1Server
	authService service.AuthService
}

func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{authService: authService}
}
