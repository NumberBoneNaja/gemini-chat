package entity

import "gorm.io/gorm"
type Conversation struct {
	gorm.Model
	Message string
	ChatRoomID uint
	ChatRoom ChatRoom `gorm:"foreignKey:ChatRoomID"`
    SendTypeID uint
	SendType SendType `gorm:"foreignKey:SendTypeID"`
}