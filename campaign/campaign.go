package campaign

import (
	"campaign/dto"
	"campaign/logger"
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type campaignService struct {
	db *gorm.DB
}

type Campaign struct {
	gorm.Model
	UserID          string    `json:"user_id"`
	Name            string    `gorm:"type:varchar(100);not null" json:"name"`
	Description     string    `gorm:"type:text" json:"description"`
	StartDate       time.Time `gorm:"not null" json:"start_date"`
	EndDate         time.Time `gorm:"not null" json:"end_date"`
	Budget          int64     `gorm:"not null" json:"budget"`
	Status          string    `gorm:"type:varchar(20);default:'active'" json:"status"`
	ChangeStartDate time.Time `gorm:"not null" json:"change_start_date"`
	ChangeEndDate   time.Time `gorm:"not null" json:"change_end_date"`
	ChangeBudget    int64     `gorm:"not null" json:"change_budget"`
}

type Cursor struct {
	EventID string
}

func New() (s *campaignService, err error) {
	client, err := gorm.Open(sqlite.Open("campaign.db"), &gorm.Config{})
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		logger.Println(err)
		return
	}

	s = &campaignService{
		db: client,
	}

	err = s.db.AutoMigrate(&Cursor{})
	if err != nil {
		err = errors.Wrap(err, "failed to migrate cursor")
		logger.Println(err)
		return
	}

	err = s.db.AutoMigrate(&Campaign{})
	if err != nil {
		err = errors.Wrap(err, "failed to migrate campaign")
		logger.Println(err)
		return
	}

	newCampaign := Campaign{
		Name:            "Campaign A",
		Description:     "Deskripsi campaign A",
		StartDate:       time.Now(),
		EndDate:         time.Now().AddDate(0, 1, 0), // 1 bulan kemudian
		Budget:          1000,
		Status:          "active",
		ChangeStartDate: time.Now(),
		ChangeEndDate:   time.Now().AddDate(0, 0, 15), // 15 hari kemudian
		ChangeBudget:    500,
	}

	if err := s.db.Create(&newCampaign).Error; err != nil {
		fmt.Println("Error creating campaign:", err)
	} else {
		fmt.Println("Campaign created:", newCampaign)
	}

	cursor := Cursor{}

	s.db.First(&cursor)
	if cursor.EventID == "" {
		s.db.Create(&Cursor{EventID: "0"})
	}

	return
}

func (c *campaignService) Project(ctx context.Context, eventID string, event dto.Event, dateTime time.Time) (err error) {
	entity, eventName := dto.ExtractEvent(event)
	if entity != "Campaign" {
		return
	}

	switch eventName {
	case "CampaignCreated":
		err = c.campaignCreated(ctx, eventID, event, dateTime)
	case "CampaignUpdated":
		err = c.campaignUpdated(ctx, eventID, event, dateTime)
	case "CampaignDeleted":
		err = c.campaignDeleted(ctx, eventID, event, dateTime)
	}
	return
}

func (c *campaignService) campaignDeleted(ctx context.Context, eventID string, event dto.Event, dateTime time.Time) (err error) {
	err = c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
		campaignID := event.Campaign.CampaignDeleted.CampaignID

		err = tx.Delete(&Campaign{}, campaignID).Error
		if err != nil {
			err = errors.Wrap(err, "failed to delete campaign")
			logger.Println(err)
			return
		}

		resp := tx.Model(&Cursor{}).Where("1 = 1").Update("EventID", eventID)
		if resp.Error != nil {
			err = errors.Wrap(resp.Error, "failed to update cursor")
			logger.Println(err)
			return
		}

		logger.Println(resp.RowsAffected)

		return
	})

	return
}

func (c *campaignService) campaignCreated(ctx context.Context, eventID string, event dto.Event, dateTime time.Time) (err error) {
	err = c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
		data := Campaign{
			UserID:      event.Campaign.CampaignCreated.UserID,
			Name:        event.Campaign.CampaignCreated.Name,
			Description: event.Campaign.CampaignCreated.Description,
			StartDate:   event.Campaign.CampaignCreated.StartDate,
			EndDate:     event.Campaign.CampaignCreated.EndDate,
			Budget:      event.Campaign.CampaignCreated.Budget,
			Status:      "active",
		}

		err = tx.Create(&data).Error
		if err != nil {
			err = errors.Wrap(err, "failed to create campaign")
			logger.Println(err)
			return
		}

		resp := tx.Model(&Cursor{}).Where("1 = 1").Update("EventID", eventID)
		if resp.Error != nil {
			err = errors.Wrap(resp.Error, "failed to update cursor")
			logger.Println(err)
			return
		}

		logger.Println(resp.RowsAffected)

		return
	})

	return
}

func (c *campaignService) campaignUpdated(ctx context.Context, eventID string, event dto.Event, dateTime time.Time) (err error) {
	err = c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
		campaignID := event.Campaign.CampaignUpdated.CampaignID

		updatedData := Campaign{
			UserID:          event.Campaign.CampaignUpdated.UserID,
			Name:            event.Campaign.CampaignUpdated.Name,
			Description:     event.Campaign.CampaignUpdated.Description,
			StartDate:       event.Campaign.CampaignUpdated.StartDate,
			EndDate:         event.Campaign.CampaignUpdated.EndDate,
			Budget:          event.Campaign.CampaignUpdated.Budget,
			ChangeStartDate: event.Campaign.CampaignUpdated.ChangeStartDate,
			ChangeEndDate:   event.Campaign.CampaignUpdated.ChangeEndDate,
			ChangeBudget:    event.Campaign.CampaignUpdated.ChangeBudget,
			Status:          event.Campaign.CampaignUpdated.Status,
		}

		err = tx.Model(&Campaign{}).Where("id = ?", campaignID).Updates(updatedData).Error
		if err != nil {
			err = errors.Wrap(err, "failed to update campaign")
			logger.Println(err)
			return
		}

		resp := tx.Model(&Cursor{}).Where("1 = 1").Update("EventID", eventID)
		if resp.Error != nil {
			err = errors.Wrap(resp.Error, "failed to update cursor")
			logger.Println(err)
			return
		}

		logger.Println(resp.RowsAffected)

		return
	})

	return err
}

func (c *campaignService) GetCampaign(ctx context.Context, id uint) (campaign *Campaign, err error) {
	campaign = &Campaign{}
	err = c.db.WithContext(ctx).First(campaign, id).Error
	if err != nil {
		err = errors.Wrap(err, "failed to get campaign")
		logger.Println(err)
		return
	}
	return
}

func (c *campaignService) GetCursor() (cursor Cursor, err error) {
	cursor = Cursor{}
	err = c.db.Last(&cursor).Error
	return
}

// --- CREATE ---

// 	// --- READ ---
// 	var campaign Campaign
// 	// Mencari campaign berdasarkan primary key (ID)
// 	if err := db.First(&campaign, newCampaign.ID).Error; err != nil {
// 		fmt.Println("Error reading campaign:", err)
// 	} else {
// 		fmt.Println("Campaign found:", campaign)
// 	}

// 	// --- UPDATE ---
// 	// Mengupdate nama dan budget campaign
// 	if err := db.Model(&campaign).Updates(Campaign{Name: "Updated Campaign A", Budget: 2000}).Error; err != nil {
// 		fmt.Println("Error updating campaign:", err)
// 	} else {
// 		// Ambil data terbaru
// 		db.First(&campaign, campaign.ID)
// 		fmt.Println("Campaign updated:", campaign)
// 	}

// 	// --- DELETE ---
// 	if err := db.Delete(&campaign).Error; err != nil {
// 		fmt.Println("Error deleting campaign:", err)
// 	} else {
// 		fmt.Println("Campaign deleted")
// 	}
// }
