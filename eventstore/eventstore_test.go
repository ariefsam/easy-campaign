package eventstore

import (
	"campaign/dto"
	"campaign/idgenerator"
	"context"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	godotenv.Load()
	ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)
	es, err := New(ctx)
	require.NoError(t, err)
	require.NotNil(t, es)

	dataEvent := dto.Event{}
	dataEvent.User.UserCreated.Email = idgenerator.New().Generate(ctx) + "john@doe.com"

	err = es.Save(ctx, dataEvent)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)
	dataEvent.User.UserCreated.Email = idgenerator.New().Generate(ctx) + "x@xx.com"
	err = es.Save(ctx, dataEvent)
	require.NoError(t, err)

	<-ctx.Done()

}
