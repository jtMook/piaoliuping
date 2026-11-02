package models

import (
	"drift-bottle-server/internal/config"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移
	if err := DB.AutoMigrate(
		&User{},
		&Bottle{},
		&BottlePickedRecord{},
		&BottleReply{},
		&Conversation{},
		&Message{},
		&Report{},
		&Blacklist{},
		&CheckinRecord{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connected and migrated successfully")
}
