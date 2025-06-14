package entity

import (
	"time"

	"gorm.io/gorm"
)

type ChatRoom struct {      
	gorm.Model
	StartDate time.Time
	EndDate   time.Time

	UsersID uint
	Users   User `gorm:"foreignKey:UsersID"`

	//fk
}