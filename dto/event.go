package dto

import (
	"encoding/json"
	"time"
)

/*
expectEvent.Influencer.InfluencerUpdated.Name = "Ucup Edit"
	expectEvent.Influencer.InfluencerUpdated.InstagramUsername = "ucupgramx"
	expectEvent.Influencer.InfluencerUpdated.TiktokUsername = "ucuptoky"
	expectEvent.Influencer.InfluencerUpdated.InfluencerID = expectID
	expectEvent.Influencer.InfluencerUpdated.UpdatedBy = state.Session.UserID
*/

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
			InfluencerID      string `json:"influencer_id,omitzero"`
			Name              string `json:"name,omitzero"`
			TiktokUsername    string `json:"tiktok_username,omitzero"`
			InstagramUsername string `json:"instagram_username,omitzero"`
			CreatedBy         string `json:"created_by,omitzero"`
		} `json:"influencer_created,omitzero"`
		InfluencerUpdated struct {
			InfluencerID string `json:"influencer_id,omitzero"`
			Name         string `json:"name,omitzero"`
			UpdatedBy    string `json:"updated_by,omitzero"`
		} `json:"influencer_updated,omitzero"`
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
	Plan Plan `json:"plan,omitzero"`
}

type Plan struct {
	PlanCreated PlanCreated `json:"plan_created,omitzero"`
}

type PlanCreated struct {
	PlanID    string    `json:"plan_id,omitzero"`
	Name      string    `json:"name,omitzero"`
	StartDate time.Time `json:"start_date,omitzero"`
	EndDate   time.Time `json:"end_date,omitzero"`
	CreatedBy string    `json:"created_by,omitzero"`
}

func (e *Event) GetEntityList() (entityList []string) {
	entityList = []string{
		"user", "auth", "session", "influencer", "instagram_account", "campaign",
	}

	return
}

func ExtractEvent(event any) (entityName string, eventName string) {

	data, err := json.Marshal(event)
	if err != nil {
		return "", ""
	}

	var result map[string]map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return "", ""
	}

	for entityName, events := range result {
		for eventName := range events {
			return entityName, eventName
		}
	}
	return
}
