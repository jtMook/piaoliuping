package main

import (
	"drift-bottle-server/internal/config"
	"drift-bottle-server/internal/models"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 初始化配置
	cfg := config.LoadConfig()

	// 初始化数据库
	models.InitDB(cfg)

	// 初始化现有用户的瓶子余额
	var users []models.User
	models.DB.Find(&users)

	for _, user := range users {
		var bottleCount int64
		models.DB.Model(&models.Bottle{}).Where("user_id = ?", user.ID).Count(&bottleCount)

		if bottleCount > 0 {
			models.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("bottle_count", bottleCount)
			fmt.Printf("用户 %d (%s) 初始化瓶子余额: %d\n", user.ID, user.Nickname, bottleCount)
		}
	}

	fmt.Println("瓶子余额初始化完成")
}
