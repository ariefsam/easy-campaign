package report

import (
	"campaign/dto"
	"campaign/logger"
	"context"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type reportService struct {
	db *gorm.DB
}

type Cursor struct {
	EventID string
}

func New() *reportService {
	client, err := gorm.Open(sqlite.Open("report.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	report := &reportService{
		db: client,
	}

	err = report.db.AutoMigrate(&Cursor{})
	if err != nil {
		logger.Println(err)
	}

	err = report.db.AutoMigrate(&Influencer{})
	if err != nil {
		logger.Println(err)
	}

	err = report.db.AutoMigrate(&Plan{})
	if err != nil {
		logger.Println(err)
	}

	return report
}

func (r *reportService) Reset() {
	r.db.Unscoped().Where("1 = 1").Delete(&Influencer{})
	r.db.Unscoped().Where("1 = 1").Delete(&Cursor{})
}

func (r *reportService) GetGroupName() string {
	return "report_worker"
}

func (r *reportService) SubscribedTo() []string {
	allEvent := dto.Event{}
	events := allEvent.GetEntityList()
	return events
}

func (s *reportService) GetCursor() (eventID string, err error) {
	cursor := Cursor{}
	err = s.db.Last(&cursor).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			eventID = "0"
			s.db.Create(&Cursor{EventID: "0"})
			return
		}
		err = errors.Wrap(err, "failed to get cursor")
		logger.Println(err)
		return
	}

	eventID = cursor.EventID
	return
}

func (s *reportService) Project(ctx context.Context, eventID string, event dto.Event, dateTime time.Time) (err error) {
	entity, eventName := dto.ExtractEvent(event)
	switch entity {
	case "campaign":
		s.handleCampaignEvent(ctx, eventID, eventName, event, dateTime)
	case "influencer":
		s.handleInfluencerEvent(ctx, eventID, eventName, event, dateTime)
	case "plan":
		s.handlePlanEvent(ctx, eventID, eventName, event, dateTime)
	}
	return
}

func (s *reportService) handleCampaignEvent(ctx context.Context, eventID, eventName string, event dto.Event, dateTime time.Time) {
	campaignHandler := &campaignHandler{
		db: s.db,
	}
	switch eventName {
	case "campaign_created":
		campaignHandler.handleCampaignCreated(ctx, eventID, event, dateTime)
	case "campaign_updated":
		campaignHandler.handleCampaignUpdated(ctx, eventID, event, dateTime)
	case "campaign_deleted":
		campaignHandler.handleCampaignDeleted(ctx, eventID, event, dateTime)
	}
}

type campaignHandler struct {
	db *gorm.DB
}

type Influencer struct {
	gorm.Model
	InfluencerID               string    `json:"influencer_id" gorm:"unique"`
	Name                       string    `json:"name"`
	InstagramUsername          string    `json:"instagram_username"`
	TiktokUsername             string    `json:"tiktok_username"`
	IsInstagramUsernameValid   bool      `json:"is_instagram_username_valid"`
	IsTiktokUsernameValid      bool      `json:"is_tiktok_username_valid"`
	LastCheckInstagramUsername time.Time `json:"last_check_instagram_username"`
	LastCheckTiktokUsername    time.Time `json:"last_check_tiktok_username"`
}

func (c *campaignHandler) handleCampaignCreated(ctx context.Context, eventID string, event dto.Event, dateTime time.Time) {

}

func (c *campaignHandler) handleCampaignUpdated(ctx context.Context, eventID string, event dto.Event, dateTime time.Time) {
}

func (c *campaignHandler) handleCampaignDeleted(ctx context.Context, eventID string, event dto.Event, dateTime time.Time) {
}

type influencerHandler struct {
	db *gorm.DB
}

func (r *reportService) handleInfluencerEvent(ctx context.Context, eventID, eventName string, event dto.Event, dateTime time.Time) {
	i := &influencerHandler{
		db: r.db,
	}
	switch eventName {
	case "influencer_created":
		i.handleInfluencerCreated(ctx, eventID, event, dateTime)
	case "influencer_updated":
		i.handleInfluencerUpdated(ctx, eventID, event, dateTime)
	case "influencer_deleted":
		i.handleInfluencerDeleted(ctx, eventID, event, dateTime)
	}
}

func (i *influencerHandler) handleInfluencerCreated(ctx context.Context, eventID string, event dto.Event, dateTime time.Time) {
	i.db.Transaction(func(db *gorm.DB) (err error) {
		influencer := Influencer{
			InfluencerID:               event.Influencer.InfluencerCreated.InfluencerID,
			Name:                       event.Influencer.InfluencerCreated.Name,
			InstagramUsername:          event.Influencer.InfluencerCreated.InstagramUsername,
			TiktokUsername:             event.Influencer.InfluencerCreated.TiktokUsername,
			IsInstagramUsernameValid:   false,
			IsTiktokUsernameValid:      false,
			LastCheckInstagramUsername: time.Time{},
			LastCheckTiktokUsername:    time.Time{},
		}
		err = db.Create(&influencer).Error
		if err != nil {
			err = errors.Wrap(err, "failed to create influencer")
			logger.Println(err)
			return
		}

		resp := db.Model(&Cursor{}).Where("1 = 1").Update("EventID", eventID)
		if resp.Error != nil {
			err = errors.Wrap(resp.Error, "failed to update cursor")
			logger.Println(err)
			return
		}
		return
	})
}

func (i *influencerHandler) handleInfluencerUpdated(ctx context.Context, eventID string, event dto.Event, dateTime time.Time) {
	i.db.Transaction(func(db *gorm.DB) (err error) {
		influencer := Influencer{
			InfluencerID: event.Influencer.InfluencerUpdated.InfluencerID,
			Name:         event.Influencer.InfluencerUpdated.Name,
		}
		err = db.Model(&Influencer{}).Where("influencer_id = ?", event.Influencer.InfluencerUpdated.InfluencerID).Updates(&influencer).Error
		if err != nil {
			err = errors.Wrap(err, "failed to update influencer")
			logger.Println(err)
			return
		}

		resp := db.Model(&Cursor{}).Where("1 = 1").Update("EventID", eventID)
		if resp.Error != nil {
			err = errors.Wrap(resp.Error, "failed to update cursor")
			logger.Println(err)
			return
		}
		return
	})

}

func (i *influencerHandler) handleInfluencerDeleted(ctx context.Context, eventID string, event dto.Event, dateTime time.Time) {
	i.db.Transaction(func(db *gorm.DB) (err error) {
		err = db.Where("influencer_id = ?", event.Influencer.InfluencerDeleted.InfluencerID).Delete(&Influencer{}).Error
		if err != nil {
			err = errors.Wrap(err, "failed to delete influencer")
			logger.Println(err)
			return
		}

		resp := db.Model(&Cursor{}).Where("1 = 1").Update("EventID", eventID)
		if resp.Error != nil {
			err = errors.Wrap(resp.Error, "failed to update cursor")
			logger.Println(err)
			return
		}
		return
	})
}

func (r *reportService) FetchInfluencers() (influencers []Influencer, err error) {
	err = r.db.Order("id DESC").Find(&influencers).Error
	if err != nil {
		err = errors.Wrap(err, "failed to fetch influencers")
		logger.Println(err)
		return
	}
	return
}

func (r *reportService) GetInfluencer(influencerID string) (influencer *Influencer, err error) {
	influencer = &Influencer{}
	err = r.db.Where("influencer_id = ?", influencerID).First(influencer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		err = errors.Wrap(err, "failed to get influencer")
		logger.Println(err)
		return
	}
	return
}
