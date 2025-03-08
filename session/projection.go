package session

import (
	"campaign/dto"
	"campaign/logger"
	"context"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type sessionService struct {
	db *gorm.DB
}

func (s *sessionService) GetGroupName() string {
	return "session_worker"
}

func (s *sessionService) SubscribedTo() []string {
	return []string{"session"}
}

func New() (s *sessionService, err error) {
	client, err := gorm.Open(sqlite.Open("session.db"), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		logger.Println(err)
		return
	}

	s = &sessionService{
		db: client,
	}

	err = s.db.AutoMigrate(&Cursor{})
	if err != nil {
		err = errors.Wrap(err, "failed to migrate cursor")
		logger.Println(err)
		return
	}

	err = s.db.AutoMigrate(&Session{})
	if err != nil {
		err = errors.Wrap(err, "failed to migrate session")
		logger.Println(err)
		return
	}

	cursor := Cursor{}

	s.db.First(&cursor)
	if cursor.EventID == "" {
		s.db.Create(&Cursor{EventID: "0"})
	}

	return
}

func (s *sessionService) Reset() {
	s.db.Unscoped().Where("1 = 1").Delete(&Session{})
	s.db.Unscoped().Where("1 = 1").Delete(&Cursor{})
}

type Cursor struct {
	EventID string
}

type Session struct {
	gorm.Model
	LoginID string `json:"login_id" gorm:"unique"`
	Email   string `json:"email"`
	UserID  string `json:"user_id"`
	Status  string `json:"status"`
}

func (s *sessionService) Project(ctx context.Context, eventID string, event dto.Event, dateTime time.Time) (err error) {
	entity, eventName := dto.ExtractEvent(event)

	if entity != "session" {
		return
	}

	switch eventName {
	case "login_succeeded":
		err = s.loginSucceeded(ctx, eventID, event, dateTime)
	}
	return
}

func (s *sessionService) loginSucceeded(ctx context.Context, eventID string, event dto.Event, dateTime time.Time) (err error) {
	err = s.db.Transaction(func(tx *gorm.DB) (err error) {
		data := Session{
			LoginID: event.Session.LoginSucceeded.LoginID,
			Email:   event.Session.LoginSucceeded.Email,
			UserID:  event.Session.LoginSucceeded.UserID,
			Status:  "active",
		}

		err = tx.Create(&data).Error
		if err != nil {
			err = errors.Wrap(err, "failed to create session")
			logger.Println(err)
			return
		}

		resp := tx.Model(&Cursor{}).Where("1 = 1").Update("EventID", eventID)
		if resp.Error != nil {
			err = errors.Wrap(resp.Error, "failed to update cursor")
			logger.Println(err)
			return
		}

		return
	})

	return
}

func (s *sessionService) GetSession(loginID string) (sess *Session, err error) {
	sess = &Session{}
	err = s.db.Where("login_id = ?", loginID).First(sess).Error
	if err != nil {
		err = errors.Wrap(err, "failed to get session")
		logger.Println(err)
		return
	}
	return
}

func (s *sessionService) GetCursor() (eventID string, err error) {
	cursor := Cursor{}
	err = s.db.Last(&cursor).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(err, "failed to get cursor")
			logger.Println(err)
			return
		}
		cursor = Cursor{EventID: "0"}
		s.db.Create(&cursor)
	}

	eventID = cursor.EventID
	return
}
