package main

import (
	"log"
	"os"

	"wms-backend-go/config"
	"wms-backend-go/routes"

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
	routes.RegisterRoutes(router)

	log.Printf("Server running on http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
