package db

import (
	"fmt"
	"github.com/alirezamastery/graph_task/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func SetupDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASS"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
		// Logger:          logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalln("error in connecting to database:", err)
		return nil
	} else {
		log.Println("Connected to Database")
	}

	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db
}

func MigrateDB(db *gorm.DB) {
	err := db.Debug().AutoMigrate(
		&models.TodoItem{},
	)

	if err != nil {
		log.Fatalln(fmt.Errorf("error migrating users: %v", err))
	}
}
