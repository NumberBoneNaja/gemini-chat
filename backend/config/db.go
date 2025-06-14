package config

import (
	"fmt"
	"four/entity"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"  // or the Docker service name if running in another container
	port     = 5432         // default PostgreSQL port
	user     = "myuser"     // as defined in docker-compose.yml
	password = "mypassword" // as defined in docker-compose.yml
	dbname   = "mydatabase" // as defined in docker-compose.yml
)

var db *gorm.DB


func DB() *gorm.DB {

   return db

}
func ConnectionDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open a connection
	var err error
	db, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	fmt.Println("Successfully connected!")

	
}

func SetupDatabase() {
	db.AutoMigrate(&entity.User{}, &entity.ChatRoom{}, &entity.Conversation{}, &entity.SendType{})
	SeedSendTypes(db)
	SeedUsers(db)
	SeedChatRooms(db)
	SeedConversations(db)
}

func SeedSendTypes(db *gorm.DB) {
	var count int64
	db.Model(&entity.SendType{}).Count(&count)
	if count == 0 {
		sendTypes := []entity.SendType{
			{Type: "user"},
			{Type: "model"},
			
		}
		db.Create(&sendTypes)
		fmt.Println("✅ Seeded SendTypes")
	}
}

// ✅ 2. Seed User
func SeedUsers(db *gorm.DB) {
	var count int64
	db.Model(&entity.User{}).Where("email = ?", "admin@example.com").Count(&count)
	if count == 0 {
		user := entity.User{
			Username:    "admin",
			Email:       "admin@example.com",
			Password:    "admin123",
			Facebook:    "admin.fb",
			Line:        "adminline",
			PhoneNumber: "0812345678",
			Role:        "admin",
			Gender:      "male",
			Age:         30,
		}
		db.Create(&user)
		fmt.Println("✅ Seeded User")
	}
}

// ✅ 3. Seed ChatRoom
func SeedChatRooms(db *gorm.DB) {
	var user entity.User
	if err := db.First(&user, "email = ?", "admin@example.com").Error; err != nil {
		fmt.Println("❌ Cannot find user for chatroom")
		return
	}

	var count int64
	db.Model(&entity.ChatRoom{}).Count(&count)
	if count == 0 {
		room := entity.ChatRoom{
			StartDate: time.Now(),
			EndDate:   time.Now().Add(30 * time.Minute),
			UsersID:   user.ID,
		}
		db.Create(&room)
		fmt.Println("✅ Seeded ChatRoom")
	}
}

// ✅ 4. Seed Conversation
func SeedConversations(db *gorm.DB) {
	var chatRoom entity.ChatRoom
	var sendTypeUser entity.SendType
	var sendTypeBot entity.SendType

	// ดึง chatroom ที่มีอยู่
	db.First(&chatRoom)

	// ดึง SendType สำหรับผู้ใช้และบอท
	db.First(&sendTypeUser, "type = ?", "user")      // สมมุติว่า user ใช้ type = "user"
	db.First(&sendTypeBot, "type = ?", "model")   // สมมุติว่า bot ใช้ type = "model"

	var count int64
	db.Model(&entity.Conversation{}).Count(&count)
	if count == 0 {
		conversations := []entity.Conversation{
			{
				Message:    "สวัสดีครับ ชอบดูบอลมากๆ และ ผู้หญิงสวยๆด้วย",
				ChatRoomID: chatRoom.ID,
				SendTypeID: sendTypeUser.ID,
			},
			{
				Message:    "สวัสดีครับ มีอะไรให้ช่วยไหมครับ?",
				ChatRoomID: chatRoom.ID,
				SendTypeID: sendTypeBot.ID,
			},
			{
				Message:    "ผมอยากรู้ว่าวันนี้อากาศเป็นยังไง",
				ChatRoomID: chatRoom.ID,
				SendTypeID: sendTypeUser.ID,
			},
			{
				Message:    "วันนี้อากาศแจ่มใส อุณหภูมิประมาณ 32 องศาเซลเซียส",
				ChatRoomID: chatRoom.ID,
				SendTypeID: sendTypeBot.ID,
			},
		}

		if err := db.Create(&conversations).Error; err != nil {
			fmt.Println("❌ Failed to seed conversations:", err)
		} else {
			fmt.Println("✅ Seeded multiple Conversations")
		}
	}
}
