package models

import (
	"errors"
	"fmt"
	"log"
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

func GetURLDetail(short_code string) (*URL, error) {
	db := db.DB()

	var url URL
	result := db.Where("short_code = ?", short_code).First(&url)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, errors.New("error looking for short URL")
	}

	return &url, nil
}

func CreateURLClick(url *URL, ip_address string, user_agent string) (*URLClick, error) {
	db := db.DB()

	click := URLClick{
		URL_ID:     url.ID,
		IP_Address: ip_address,
		User_Agent: user_agent,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&click).Error; err != nil {
			log.Printf("Failed to create URLClick record: %v", err)
			return fmt.Errorf("could not create URLClick record: %v", err)
		}

		url.Clicks++
		if err := tx.Save(&url).Error; err != nil {
			log.Printf("Failed to save URL record: %v", err)
			return fmt.Errorf("could not save URL record: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &click, nil
}
