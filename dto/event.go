package dto

import (
	"reflect"
	"time"
)

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
	Influencer struct {
		InfluencerCreated struct {
			InfluencerID string `json:"influencer_id,omitzero"`
			Name         string `json:"name,omitzero"`
			CreatedBy    string `json:"created_by,omitzero"`
		} `json:"influencer_created,omitzero"`
	} `json:"influencer,omitzero"`
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
	Campaign struct {
		CampaignCreated struct {
			ID          uint      `json:"id,omitzero"`
			UserID      string    `json:"user_id,omitzero"`
			Name        string    `json:"name,omitzero"`
			Description string    `json:"description,omitzero"`
			StartDate   time.Time `json:"start_date,omitzero"`
			EndDate     time.Time `json:"end_date,omitzero"`
			Budget      int64     `json:"budget,omitzero"`
		} `json:"campaign_created,omitzero"`
		CampaignUpdated struct {
			ID              uint      `json:"id,omitzero"`
			UserID          string    `json:"user_id,omitzero"`
			Name            string    `json:"name,omitzero"`
			Description     string    `json:"description,omitzero"`
			StartDate       time.Time `json:"start_date,omitzero"`
			EndDate         time.Time `json:"end_date,omitzero"`
			Budget          int64     `json:"budget,omitzero"`
			ChangeStartDate time.Time `json:"change_start_date,omitzero"`
			ChangeEndDate   time.Time `json:"change_end_date,omitzero"`
			ChangeBudget    int64     `json:"change_budget,omitzero"`
			Status          string    `json:"status,omitzero"`
		} `json:"campaign_updated,omitzero"`
		CampaignDeleted struct {
			ID uint `json:"id,omitzero"`
		} `json:"campaign_deleted,omitzero"`
	} `json:"campaign,omitzero"`
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
