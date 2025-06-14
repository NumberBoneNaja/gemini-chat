package entity

import "gorm.io/gorm"

type SendType struct {
	gorm.Model
	Type string

	Conversations []Conversation `gorm:"foreignKey:SendTypeID"`
}