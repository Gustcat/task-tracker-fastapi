package model

import (
	"github.com/dgrijalva/jwt-go"
)

type UserClaims struct {
	jwt.StandardClaims
	Name string `json:"name"`
	Role int32  `json:"role"`
	ID   int64  `json:"id"`
}
