package main

import (
	"log"
	"subscription-aggregator/internal/handler"
	"subscription-aggregator/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file")
	}

	repository.InitAndMigrateDB()

	r := gin.Default()

	r.POST("/create", handler.CreateSubscription)
	r.GET("/read/:id", handler.ReadSubscription)
	r.PUT("/update/:id", handler.UpdateSubscription)
	r.DELETE("/delete/:id", handler.DeleteSubscription)
	r.GET("/list", handler.ListSubscriptions)
	r.GET("/sum", handler.SumSubscriptionsPrice)
}
