package campaign

import (
	"campaign/dto"
	"campaign/eventstore"
	"campaign/idgenerator"
	"campaign/logger"
	"campaign/token"
	"context"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	eventStore     eventStore
	idGenerator    idGenerator
	tokenGenerator tokenGenerator
	tokenParser    tokenParser
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type idGenerator interface {
	Generate(ctx context.Context) string
	GenerateUUID(ctx context.Context) string
}

type eventStore interface {
	Save(ctx context.Context, event dto.Event) (err error)
}

type tokenGenerator interface {
	Generate(ctx context.Context, session dto.Session) string
}

type tokenParser interface {
	Parse(ctx context.Context, tokenString string) (session *dto.Session, err error)
}

func NewAuthService() (a *AuthService, err error) {
	ctx := context.TODO()

	idGenerator := idgenerator.New()
	tokenService := token.New()
	eventStore, _ := eventstore.New(ctx)

	return &AuthService{
		eventStore:     eventStore,
		idGenerator:    idGenerator,
		tokenGenerator: tokenService,
		tokenParser:    tokenService,
	}, nil
}

func (auth *AuthService) SetEventStore(eventStore eventStore) {
	auth.eventStore = eventStore
}

func (auth *AuthService) SetTokenParser(tokenParser tokenParser) {
	auth.tokenParser = tokenParser
}

func (auth *AuthService) SetTokenGenerator(tokenGenerator tokenGenerator) {
	auth.tokenGenerator = tokenGenerator
}

func (auth *AuthService) SetIDGenerator(idGenerator idGenerator) {
	auth.idGenerator = idGenerator
}

func (auth *AuthService) Login(ctx context.Context, payload *Request, state *InternalState, resp *Response) (err error) {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminBcryptPassword := os.Getenv("ADMIN_BCRYPT_PASSWORD")

	bcryptPassword := ""
	userID := ""

	if payload.Login.Email == adminEmail {
		bcryptPassword = adminBcryptPassword
		userID = "root"
	}

	now := time.Now()

	if CheckPasswordHash(payload.Login.Password, bcryptPassword) {
		loginID := auth.idGenerator.Generate(ctx)
		session := dto.Session{}

		session.Email = payload.Login.Email
		session.Id = loginID
		session.Issuer = "easy-campaign"
		session.Audience = "easy-campaign"
		session.ExpiresAt = now.Add(24 * time.Hour).Unix()
		session.IssuedAt = now.Unix()
		session.NotBefore = now.Unix()
		session.Subject = userID

		resp.Auth.Token = auth.tokenGenerator.Generate(ctx, session)

		event := dto.Event{}
		event.Session.LoginSucceeded.Email = payload.Login.Email
		event.Session.LoginSucceeded.LoginID = loginID
		event.Session.LoginSucceeded.UserID = userID
		err = auth.eventStore.Save(ctx, event)
		if err != nil {
			logger.Error(err)
			return
		}

		return
	}

	resp.StatusCode = http.StatusUnauthorized
	resp.Error = "Invalid username or password"

	return
}

func CheckPasswordHash(password, hash string) (valid bool) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
