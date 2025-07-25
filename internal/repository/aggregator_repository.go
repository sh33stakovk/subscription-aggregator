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

	var err error
	maxRetries := 10

	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(postgres.Open(conStr))
		if err == nil {
			break
		}

		log.Printf("failed to connect to database, attempt: %v/%v", i, maxRetries)
		time.Sleep(time.Second * 2)
	}

	if DB == nil {
		log.Fatal("unable to connect to database")
	}

	log.Println("connected to database")

	err = DB.AutoMigrate(model.Subscription{})
	if err != nil {
		log.Fatalf("unable to migrate to database: %v", err.Error())
	}

	log.Println("migrated to database")
}
