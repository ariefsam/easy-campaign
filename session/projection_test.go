package session_test

import (
	"campaign/dto"
	"campaign/idgenerator"
	"campaign/session"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_sessionService_Project(t *testing.T) {
	ctx := context.TODO()
	sessionService, err := session.New()
	require.NoError(t, err)
	require.NotNil(t, sessionService)

	dataEvent := dto.Event{}
	dataEvent.Session.LoginSucceeded.Email = "admin@gmail.com"
	dataEvent.Session.LoginSucceeded.LoginID = "loginXXX" + idgenerator.New().Generate(ctx)
	dataEvent.Session.LoginSucceeded.UserID = "userXXX"
	now := time.Now()
	err = sessionService.Project(ctx, "e1", dataEvent, now)
	require.NoError(t, err)

	sess, err := sessionService.GetSession(dataEvent.Session.LoginSucceeded.LoginID)
	require.NoError(t, err)
	require.NotNil(t, sess)

	require.Equal(t, dataEvent.Session.LoginSucceeded.LoginID, sess.LoginID)
	require.Equal(t, "userXXX", sess.UserID)
	require.Equal(t, "admin@gmail.com", sess.Email)

	cursor, err := sessionService.GetCursor()
	require.NoError(t, err)
	require.NotNil(t, cursor)
	require.Equal(t, "e1", cursor.EventID)
}
