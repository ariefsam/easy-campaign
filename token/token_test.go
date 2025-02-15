package token_test

import (
	"campaign/dto"
	"campaign/logger"
	"campaign/token"
	"context"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func Test_token_Generate(t *testing.T) {
	ctx := context.TODO()
	godotenv.Load()
	tokenService := token.New()
	cl := dto.Session{
		Email: "admin@gmail.com",
	}
	cl.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
	tokenString := tokenService.Generate(ctx, cl)
	t.Log(tokenString)
	require.NotEmpty(t, tokenString)

	claim, err := tokenService.Parse(ctx, tokenString)
	require.NoError(t, err)
	require.Equal(t, cl, *claim)

	logger.PrintJSON(claim.Valid())

}
