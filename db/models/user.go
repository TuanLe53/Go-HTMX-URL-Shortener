package models

import (
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
