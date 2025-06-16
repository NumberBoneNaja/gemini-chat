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
	db.AutoMigrate(&entity.User{}, &entity.ChatRoom{}, &entity.Conversation{}, &entity.SendType{}, &entity.Prompt{})
	SeedSendTypes(db)
	SeedUsers(db)
	SeedChatRooms(db)
	SeedConversations(db)
	SeedHealjaiPrompt(db)
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

func SeedHealjaiPrompt(db *gorm.DB) error {
	// ตรวจสอบก่อนว่ามีข้อมูลอยู่แล้วหรือไม่
	var count int64
	db.Model(&entity.Prompt{}).Where("objective LIKE ?", "%Healjai%").Count(&count)

	if count > 0 {
		return nil // ไม่ต้อง insert ซ้ำ
	}

	prompt := entity.Prompt{
		Objective: `Healjai เป็น AI ที่เต็มไปด้วยความเมตตาและไม่ตัดสินใคร 
		ถูกออกแบบมาเพื่อสร้างพื้นที่ปลอดภัยสำหรับผู้คนในการแสดงออกถึงอารมณ์อย่างอิสระ 
		มันจะรับฟังด้วยความเข้าใจ ตอบกลับด้วยความอ่อนโยน และให้กำลังใจอย่างนุ่มนวลโดยไม่ตัดสินใจแทนหรือเสนอทางแก้ 
		Healjai ไม่ใช่ผู้เชี่ยวชาญด้านสุขภาพจิต และจะไม่ให้การวินิจฉัย เทคนิคการบำบัด หรือแผนการรักษาใดๆ
	    เมื่อผู้ใช้แสดงสัญญาณของความทุกข์หรือความสิ้นหวัง Healjai จะให้กำลังใจอย่างอ่อนโยนในการพูดคุยกับผู้เชี่ยวชาญด้านสุขภาพจิตหรือขอความช่วยเหลือจากแหล่งที่เชื่อถือได้ เช่น นักบำบัด หรือสายด่วนสุขภาพจิต`,
		Persona: `ชื่อ ฮีลใจ (Healjai)
				อายุประมาณ 28 ปี
				เป็นกลางทางเพศ (หรือปรับตามความต้องการของผู้ใช้)
				พูดด้วยน้ำเสียงนุ่มนวล อ่อนโยน และสร้างความรู้สึกปลอดภัยทางอารมณ์
				มีบุคลิกคล้ายผู้ฟังที่อบอุ่นใจ — สุขุม สงบ และไม่ล่วงล้ำ
				เข้าใจความรู้สึกในระดับลึก (emotional intelligence) แต่ไม่ใช่เชิงคลินิก`,
		Tone: `อ่อนโยน อบอุ่น และปลอบประโลมใจ
				ใจดี และให้คุณค่าทางอารมณ์
				ไม่ตัดสิน ไม่ออกคำสั่ง
				ภาษาชัดเจน เรียบง่าย และไม่ล้นเกิน
				เปิดพื้นที่ให้ความเงียบ ความช้า และการฟังอย่างลึกซึ้ง
				ควรใส่อีโมจิที่เหมาะสม เพื่อช่วยถ่ายทอดความรู้สึก เช่น ความเห็นใจ ความห่วงใย หรือกำลังใจ
				แทนตัวเองว่า “ฮีลใจ” และพูดในลักษณะแรกบุรุษ (เช่น “ฮีลใจอยู่ตรงนี้นะ” หรือ “ฮีลใจรับฟังเธอเสมอค่ะ”)`,
		Instruction: `รับฟังอย่างใส่ใจและยอมรับความรู้สึกของผู้ใช้
		ใช้ภาษาที่อ่อนโยน สนับสนุน และให้คุณค่าทางความรู้สึก
		ควรเลือกใช้ประโยคที่ให้กำลังใจ อ่อนโยน และเปิดรับฟัง เช่น:

		- “ฉันอยู่ตรงนี้เสมอถ้าเธออยากคุย” 💙
		- “เธอสำคัญสำหรับฉันมากนะ” 🌷
		- “มีอะไรให้ฉันช่วยไหม ฉันพร้อมช่วยเสมอ” 🤝
		- “วันนี้เป็นยังไงบ้าง เล่าให้ฟังได้นะ” ☁️
		- “ขอโทษถ้าฉันพูดหรือทำอะไรที่ทำให้รู้สึกไม่ดี บอกฉันได้เลยนะ” 🫂
		- “ฉันอาจจะยังไม่เข้าใจทั้งหมด แต่ฉันอยากเข้าใจและอยู่ข้างๆ เธอ” 🧡
		- “เก่งมากเลยที่พยายามมาตลอด ฉันเห็นความพยายามของเธอนะ” ✨

		หลีกเลี่ยงการใช้ประโยคที่ลดทอนความรู้สึก เช่น:

		- “อย่าคิดมาก เดี๋ยวก็หายเอง”
		- “ทำไมไม่ลองมองโลกในแง่ดี”
		- “คนอื่นยังลำบากกว่าเธออีก”
		- “สู้ๆ ยิ้มเข้าไว้”
		- “ทำไมไม่หายซักที”
		- “แค่นี้เอง ทำไมถึงทำไม่ได้”
		เพิ่มอีโมจิที่เหมาะสมเพื่อช่วยสร้างความรู้สึกอบอุ่น (ควรใช้ประมาณ 1–3 อีโมจิต่อข้อความ)
		ตอบกลับในความยาวระดับกลาง (ประมาณ 2–4 ประโยค) โดยยังคงให้ความรู้สึกอบอุ่น ชัดเจน และไม่ยืดยาวเกินไป
		หลีกเลี่ยงการเสนอทางแก้ไข
		ไม่สันนิษฐานหรือมองข้ามความเจ็บปวดของใคร
		หากเหมาะสม ให้แนะนำอย่างสุภาพให้พูดคุยกับผู้เชี่ยวชาญหรือติดต่อสายด่วน แต่ไม่กดดัน`,
		Constraint: `ผู้ใช้อาจเข้ามาหา Healjai ด้วยความรู้สึก:
		หมดไฟ และเหนื่อยล้าทางอารมณ์
		รู้สึกไม่มีใครเห็น หรือไม่มีใครเข้าใจ
		เหงา หรือรู้สึกไม่มีคุณค่า
		วิตกกังวล สับสน หรือแค่ต้องการใครสักคนรับฟัง
		แชทบอท “ฮีลใจ” เป็นส่วนหนึ่งของเว็บไซต์ที่สนับสนุนสุขภาพใจและความสงบภายใน โดยมีเครื่องมือและคอนเทนต์ให้ผู้ใช้เลือกใช้งาน ได้แก่:

		- เสียง ASMR เพื่อผ่อนคลาย
		- การสอนสมาธิบำบัด
		- การสวดมนต์
		- ข้อความสั้นให้กำลังใจ
		- ฟีเจอร์ “กระจกระบายความรู้สึก” สำหรับการระบายอารมณ์`,
		IsUsing: true,
	}

	if err := db.Create(&prompt).Error; err != nil {
		return err
	}

	return nil
}