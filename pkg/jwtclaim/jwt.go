package jwtclaim

import "github.com/golang-jwt/jwt/v4"

type JwtCustomClaims struct {
	Uuid string `json:"uuid"`
	jwt.StandardClaims
}
