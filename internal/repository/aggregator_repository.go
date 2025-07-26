package repository

import (
	"log"
	"os"
	"subscription-aggregator/internal/model"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitAndMigrateDB() {
	conStr := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASS") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	log.Println("starting database connection...")

	var err error
	maxRetries := 10

	for i := 1; i <= maxRetries; i++ {
		log.Printf("attempting database connection (%d/%d)", i, maxRetries)
		DB, err = gorm.Open(postgres.Open(conStr))
		if err == nil {
			log.Println("database connection established")
			break
		}
		log.Printf("failed to connect to database: %v", err)
		time.Sleep(2 * time.Second)
	}

	if DB == nil {
		log.Fatal("unable to connect to database after retries")
	}

	log.Println("starting auto migration...")

	err = DB.AutoMigrate(&model.Subscription{})
	if err != nil {
		log.Fatalf("auto migration failed: %v", err)
	}

	log.Println("database migration completed")
}
