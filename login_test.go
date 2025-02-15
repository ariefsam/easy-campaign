package campaign_test

import (
	"campaign"
	"campaign/dto"
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockEventStore struct {
	mock.Mock
}

func (m *mockEventStore) Save(ctx context.Context, event dto.Event) (err error) {
	args := m.Called(event)
	return args.Error(0)
}

type mockIDGenerator struct {
	mock.Mock
}

func (m *mockIDGenerator) Generate(ctx context.Context) string {
	args := m.Called()
	return args.String(0)
}

func (m *mockIDGenerator) GenerateUUID(ctx context.Context) string {
	args := m.Called()
	return args.String(0)
}

type mockTokenGenerator struct {
	mock.Mock
}

func (m *mockTokenGenerator) Generate(ctx context.Context, session dto.Session) string {
	args := m.Called(session)
	return args.String(0)
}

type mockTokenParser struct {
	mock.Mock
}

func (m *mockTokenParser) Parse(ctx context.Context, tokenString string) (session *dto.Session, err error) {
	args := m.Called(tokenString)
	session, _ = args.Get(0).(*dto.Session)
	return session, args.Error(1)
}

func TestLogin(t *testing.T) {
	ctx := context.TODO()
	os.Setenv("ADMIN_EMAIL", "admin@gmail.com")
	os.Setenv("ADMIN_BCRYPT_PASSWORD", "$2a$12$vQaZu3QdooA4MgKySPUHfuQn/QoQIqEktoIu6hDMjpFOR2SxYHhGO") //abc123
	payload := campaign.Request{}
	state := campaign.InternalState{}
	resp := campaign.Response{}
	es := &mockEventStore{}
	idGenerator := &mockIDGenerator{}
	tokenGenerator := &mockTokenGenerator{}
	tokenParser := &mockTokenParser{}

	authService, err := campaign.NewAuthService()
	require.NoError(t, err)
	authService.SetEventStore(es)
	authService.SetTokenParser(tokenParser)
	authService.SetTokenGenerator(tokenGenerator)
	authService.SetIDGenerator(idGenerator)

	t.Run("failed login", func(t *testing.T) {
		payload.Login.Email = "admin@gmail.com"
		payload.Login.Password = "abc12"
		err := authService.Login(ctx, &payload, &state, &resp)
		require.NoError(t, err)
		require.Equal(t, "", resp.Auth.Token)
	})

	t.Run("success login admin", func(t *testing.T) {
		payload.Login.Email = "admin@gmail.com"
		payload.Login.Password = "abc123"

		idGenerator.On("Generate").Return("loginID123")

		expectEvent := dto.Event{}
		expectEvent.Session.LoginSucceeded.Email = payload.Login.Email
		expectEvent.Session.LoginSucceeded.LoginID = "loginID123"
		tokenGenerator.On("Generate", mock.Anything).Return("token123", nil)

		es.On("Save", mock.Anything).Return(nil)

		err := authService.Login(ctx, &payload, &state, &resp)
		require.NoError(t, err)

		es.AssertExpectations(t)
		assert.Equal(t, "token123", resp.Auth.Token)
		es.AssertCalled(t, "Save", mock.MatchedBy(func(event dto.Event) bool {
			return event.Session.LoginSucceeded.Email == "admin@gmail.com" &&
				event.Session.LoginSucceeded.LoginID == "loginID123" &&
				event.Session.LoginSucceeded.UserID == "root"

		}))

		require.Equal(t, "token123", resp.Auth.Token)

	})

	t.Run("parser token failed", func(t *testing.T) {
		token := "tokenxxx"
		tokenParser.On("Parse", token).Return(mock.Anything, errors.New("invalid token"))

		session, err := authService.ParseToken(ctx, token)
		require.Error(t, err)
		require.Nil(t, session)
	})

	t.Run("parser token success", func(t *testing.T) {
		token := "tokenabc"
		session := &dto.Session{}
		session.Email = "admin@gmail.com"
		tokenParser.On("Parse", token).Return(session, nil)

		result, err := authService.ParseToken(ctx, token)
		require.NoError(t, err)
		require.Equal(t, session, result)

	})

}
