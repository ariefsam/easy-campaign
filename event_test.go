package campaign

import (
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
		require.Equal(t, "User", entityName)
		require.Equal(t, "UserCreated", eventName)
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
		require.Equal(t, "User", entityName)
		require.Equal(t, "ProfileUpdatedRequested", eventName)
	})
}
