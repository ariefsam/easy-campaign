package tiktokapi_test

import (
	"campaign/logger"

	"campaign/rapidapi/tiktokapi"
	"context"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestGetUser(t *testing.T) {
	godotenv.Load()

	service := tiktokapi.New()
	username := "ustadz.khalidbasalamah"
	got, err := service.GetUser(context.TODO(), username)
	require.NoError(t, err)
	require.NotNil(t, got)
	logger.PrintJSON(got)
}
