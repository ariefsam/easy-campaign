package dto

import "reflect"

type Event struct {
	User struct {
		UserCreated struct {
			UserID    string `json:"user_id,omitzero"`
			Email     string `json:"email,omitzero"`
			FirstName string `json:"first_name,omitzero"`
			LastName  string `json:"last_name,omitzero"`
		} `json:"user_created,omitzero"`
		ProfileUpdatedRequested struct {
			RequestID       string `json:"request_id,omitzero"`
			UserID          string `json:"user_id,omitzero"`
			FirstName       string `json:"first_name,omitzero"`
			LastName        string `json:"last_name,omitzero"`
			RequesterUserID string `json:"requester_user_id,omitzero"`
		} `json:"profile_updated_requested,omitzero"`
	} `json:"user,omitzero"`
	Auth struct {
		PasswordCreated struct {
			UserID         string `json:"user_id,omitzero"`
			BcryptPassword string `json:"bcrypt_password,omitzero"`
		} `json:"password_created,omitzero"`
		PasswordUpdated struct {
			UserID         string `json:"user_id,omitzero"`
			BcryptPassword string `json:"bcrypt_password,omitzero"`
		} `json:"password_updated,omitzero"`
	} `json:"auth,omitzero"`
	Session struct {
		LoginFailed struct {
			LoginID string `json:"login_id,omitzero"`
			Email   string `json:"email,omitzero"`
		} `json:"login_failed,omitzero"`
		LoginSucceeded struct {
			UserID  string `json:"user_id,omitzero"`
			LoginID string `json:"login_id,omitzero"`
			Email   string `json:"email,omitzero"`
		} `json:"login_succeeded,omitzero"`
	} `json:"session,omitzero"`
	InstagramAccount struct {
		AccountCreated struct {
			InstagramID       string `json:"instagram_id,omitzero"`
			InstagramUsername string `json:"username,omitzero"`
			APIRawResponse    string `json:"api_response,omitzero"`
		} `json:"account_created,omitzero"`
		AccountUpdated struct {
			InstagramID       string `json:"instagram_id,omitzero"`
			InstagramUsername string `json:"username,omitzero"`
			APIRawResponse    string `json:"api_response,omitzero"`
		} `json:"account_updated,omitzero"`
	} `json:"instagram_account,omitzero"`
}

func ExtractEvent(event any) (entityName string, eventName string) {
	v := reflect.ValueOf(event)
	for i := 0; i < v.NumField(); i++ {
		entity := v.Field(i)
		entityType := v.Type().Field(i).Name
		for j := 0; j < entity.NumField(); j++ {
			event := entity.Field(j)
			if !event.IsZero() {
				eventType := entity.Type().Field(j).Name
				return entityType, eventType
			}
		}
	}
	return "", ""
}

func ExtractListEntities(event any) (entities []string) {
	v := reflect.ValueOf(event)
	for i := 0; i < v.NumField(); i++ {
		entity := v.Type().Field(i).Name
		entities = append(entities, entity)
	}
	return entities
}
