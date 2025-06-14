package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username    string
	Email       string
	Password    string
	Facebook    string
	Line        string
	PhoneNumber string
	Role        string
	Gender      string
	Age         int

	//fk
	ChatRooms []ChatRoom `gorm:"foreignKey:UsersID"`
}