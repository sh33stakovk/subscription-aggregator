package main

import (
	"log"
	"os"
	"subscription-aggregator/internal/handler"
	"subscription-aggregator/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("loading .env file...")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}
	log.Println(".env file loaded")

	log.Println("initializing database...")
	repository.InitAndMigrateDB()

	r := gin.Default()

	log.Println("registering routes...")

	r.POST("/create", handler.CreateSubscription)
	r.GET("/read/:id", handler.ReadSubscription)
	r.PUT("/update/:id", handler.UpdateSubscription)
	r.DELETE("/delete/:id", handler.DeleteSubscription)
	r.GET("/list", handler.ListSubscriptions)
	r.GET("/sum", handler.SumSubscriptionsPrice)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("PORT not set in env, using default: %s", port)
	} else {
		log.Printf("using port from env: %s", port)
	}

	log.Printf("starting server on port %s...", port)
	err = r.Run(":" + port)
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
