package models

import (
	"time"

	"github.com/TuanLe53/Go-HTMX-URL-Shortener/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type URL struct {
	gorm.Model
	ID         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Long_URL   string     `gorm:"type:text;not null"`
	Short_Code string     `gorm:"type:text;not null;unique"`
	Expires_At *time.Time `gorm:"null"`
	User_ID    *uuid.UUID `gorm:"type:uuid;not null"`
	Clicks     int        `gorm:"default:0"`
}

type URLClick struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	URL_ID     uuid.UUID `gorm:"type:uuid;not null"`
	Clicked_At time.Time `gorm:"default:now()"`
	IP_Address string    `gorm:"type:varchar(45);null"`
	User_Agent string    `gorm:"type:text;null"`
}

func CreateShortURL(long_url string, short_code string, created_by uuid.UUID, expired_at time.Time) (*URL, error) {
	db := db.DB()

	url := URL{
		Long_URL:   long_url,
		Short_Code: short_code,
		User_ID:    &created_by,
		Expires_At: &expired_at,
	}

	if err := db.Create(&url).Error; err != nil {
		return nil, err
	}

	return &url, nil
}
