package models

import (
	"errors"

	"github.com/TuanLe53/Go-HTMX-URL-Shortener/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email     string    `gorm:"type:varchar(100);not null;unique"`
	FirstName string    `gorm:"type:varchar(100);not null"`
	LastName  string    `gorm:"type:varchar(100);not null"`
	Password  string    `gorm:"type:varchar(255);not null"`
}

func FindUserWithEmail(email string) (*User, error) {
	db := db.DB()

	var user User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.New("error looking for user")
	} else {
		return &user, nil
	}
}

func CreateUser(email, firstName, lastName, password string) *User {
	db := db.DB()

	user := User{Email: email, FirstName: firstName, LastName: lastName, Password: password}

	db.Create(&user)

	return &user
}
