package token

import (
	"campaign/dto"
	"context"
	"os"

	"github.com/golang-jwt/jwt"
)

type token struct {
	secret string
}

func New() *token {
	secret := os.Getenv("JWT_SECRET")
	return &token{
		secret: secret,
	}
}

func (t *token) Generate(ctx context.Context, claim dto.Session) (tokenString string) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, _ = token.SignedString([]byte(t.secret))

	return
}

func (t *token) Parse(ctx context.Context, tokenString string) (claim *dto.Session, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &dto.Session{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(t.secret), nil
	})

	if err != nil {
		return
	}

	if claims, ok := token.Claims.(*dto.Session); ok && token.Valid {
		claim = claims
	} else {
		err = jwt.ErrSignatureInvalid
	}
	return
}
