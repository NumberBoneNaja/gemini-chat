package controller

import (
	"context"
	"fmt"
	"four/config"
	"four/entity"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
)
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Failed to load .env file")
	}
}



func Gemini() {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Println("❌ Error creating Gemini client:", err)
		return
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text("นายชื่ออะไร"),
		nil,
	)
	if err != nil {
		log.Println("❌ Error generating content:", err)
		return
	}

	fmt.Println("✅ Gemini response:", result.Text())
}

type Role string

const (
	RoleUser  Role = "user"
	RoleModel Role = "model"
)

func GeminiHistory(c *gin.Context) {
	// TODO: implement history controller
  var question entity.Conversation

  if err := c.ShouldBindJSON(&question); err != nil {

	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	return

}
  ctx := context.Background()
  client, err := genai.NewClient(ctx, &genai.ClientConfig{
      APIKey:  os.Getenv("GEMINI_API_KEY"),
      Backend: genai.BackendGeminiAPI,
  })
  if err != nil {
	log.Println("❌ Error creating Gemini client:", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Gemini client"})
	return
	}

	var historyConversations []entity.Conversation
	if err := config.DB().Where("chat_room_id = ?", question.ChatRoomID).Order("created_at asc").Find(&historyConversations).Error; err != nil {
		log.Println("❌ Error fetching history:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch conversation history"})
		return
	}

	// ✅ แปลงเป็น []*genai.Content
	var history []*genai.Content
	roleMap := map[uint]genai.Role{
		1: genai.RoleUser,
		2: genai.RoleModel,
	}
	
	for _, conv := range historyConversations {
		role := roleMap[conv.SendTypeID]
		history = append(history, genai.NewContentFromText(conv.Message, role))
	}

  chat, err := client.Chats.Create(ctx, "gemini-2.0-flash", nil, history)
  if err != nil {
	  log.Println("❌ Error creating chat:", err)
	  c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
	  return
  }

  res,err := chat.SendMessage(ctx, genai.Part{Text: question.Message})
  if err != nil {
	  log.Println("❌ Error sending message:", err)
	  c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
	  return
  }

  if len(res.Candidates) > 0 {
	answer := res.Candidates[0].Content.Parts[0].Text
	config.DB().Create(&entity.Conversation{
		Message:     question.Message,
		ChatRoomID:  question.ChatRoomID,
		SendTypeID:  1, // user
	})

	// บันทึกคำตอบ
	config.DB().Create(&entity.Conversation{
		Message:     answer,
		ChatRoomID:  question.ChatRoomID,
		SendTypeID:  2, // bot
	})

	  c.JSON(http.StatusOK, gin.H{"message": res.Candidates[0].Content.Parts[0].Text})
	
  }
}

func CreateChatRoom(c *gin.Context) {
	db := config.DB()
	resault := db.Create(&entity.ChatRoom{
		StartDate: time.Now(),
	})
	if resault.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": resault.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "create chat room success"})

	
}