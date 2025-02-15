package dto

import "github.com/golang-jwt/jwt"

type Session struct {
	jwt.StandardClaims
	Email string `json:"email"`
}
