package main

import (
	"log"
	"os"
	"time"

	"wms-backend-go/config"
	"wms-backend-go/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}

	config.ConnectDB()

	router := gin.Default()

	// ✅ CORS ミドルウェアを適用
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // 許可するオリジン
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // プリフライトリクエストの結果をキャッシュする時間
	}))

	// ルート登録
	routes.RegisterRoutes(router)

	log.Printf("Server running on http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
