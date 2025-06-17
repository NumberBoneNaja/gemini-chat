package main

import (
	"fmt"
	"four/config"
	"four/controller"

	"net/http"

	"github.com/gin-gonic/gin"
)

const PORT = "8000"

func main() {
	fmt.Println("Hello World!")
	config.ConnectionDB()
	config.SetupDatabase()
	r := gin.Default()

	r.Use(CORSMiddleware())
	r.GET("/", func(c *gin.Context) {

		c.String(http.StatusOK, "API RUNNING... PORT: %s", PORT)

	})

	r.POST("/gemini", controller.GeminiHistory)
	r.GET("/conversation/:id", controller.GetConversationHistory)
    

	r.Run("localhost:" + PORT)
}

func CORSMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")

		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {

			c.AbortWithStatus(204)

			return

		}

		c.Next()

	}

}
