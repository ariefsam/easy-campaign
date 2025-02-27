package dto

import (
	"campaign/logger"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractEvent(t *testing.T) {
	t.Run("UserCreated", func(t *testing.T) {
		data := []byte(`{
			"user": {
				"user_created": {
					"user_id": "123"
				}
			}
		}`)
		var event Event
		err := json.Unmarshal(data, &event)
		require.NoError(t, err)

		entityName, eventName := ExtractEvent(event)
		require.Equal(t, "user", entityName)
		require.Equal(t, "user_created", eventName)
	})

	t.Run("ProfileUpdatedRequested", func(t *testing.T) {
		data := []byte(`{
			"user": {
				"profile_updated_requested": {
					"request_id": "123"
				}
			}
		}`)
		var event Event
		err := json.Unmarshal(data, &event)
		require.NoError(t, err)

		entityName, eventName := ExtractEvent(event)
		require.Equal(t, "user", entityName)
		require.Equal(t, "profile_updated_requested", eventName)
	})

	t.Run("Extract List Entities", func(t *testing.T) {
		event := Event{}
		entities := event.GetEntityList()
		require.NotEmpty(t, entities)

		logger.PrintJSON(entities)
	})
}
