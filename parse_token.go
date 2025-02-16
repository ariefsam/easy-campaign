package campaign

import (
	"campaign/dto"
	"context"
)

func (a *AuthService) ParseToken(ctx context.Context, tokenString string) (session *dto.Session, err error) {
	session, err = a.tokenParser.Parse(ctx, tokenString)
	return
}
