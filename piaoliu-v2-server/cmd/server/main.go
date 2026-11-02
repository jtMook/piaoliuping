package main

import (
	"drift-bottle-server/internal/config"
	"drift-bottle-server/internal/handler"
	"drift-bottle-server/internal/middleware"
	"drift-bottle-server/internal/models"
	"drift-bottle-server/internal/websocket"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// 定时清理超过7天的瓶子
func startCleanupJob() {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			expireTime := time.Now().AddDate(0, 0, -7)
			result := models.DB.Model(&models.Bottle{}).
				Where("status = 1 AND created_at < ?", expireTime).
				Update("status", 3)
			if result.Error != nil {
				log.Printf("清理过期瓶子失败: %v", result.Error)
			} else if result.RowsAffected > 0 {
				log.Printf("已清理 %d 个过期瓶子", result.RowsAffected)
			}
		}
	}()
	log.Println("已启动过期瓶子定时清理任务（每小时检查一次，清理超过7天的瓶子）")
}

func main() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 初始化配置
	cfg := config.LoadConfig()

	// 初始化数据库
	models.InitDB(cfg)

	// 初始化 Redis
	models.InitRedis(cfg)

	// 启动定时清理任务
	startCleanupJob()

	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// 静态文件服务（上传的图片等）
	r.Static("/uploads", "./uploads")

	// WebSocket
	r.GET("/ws", websocket.HandleWebSocket)

	// 公开接口
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", handler.Login)
	}

	// 需要认证的接口
	api := r.Group("/api")
	api.Use(middleware.JWTAuth())
	{
		// 用户
		users := api.Group("/users")
		{
			users.GET("/profile", handler.GetProfile)
			users.PUT("/profile", handler.UpdateProfile)
			users.PUT("/address", handler.UpdateAddress)
		}

		// 瓶子
		bottles := api.Group("/bottles")
		{
			bottles.POST("", handler.CreateBottle)
			bottles.POST("/pick", handler.PickBottle)
			bottles.GET("/:id", handler.GetBottleDetail)
			bottles.POST("/:id/reply", handler.ReplyBottle)
			bottles.POST("/:id/comment", handler.CommentBottle)
			bottles.GET("/:id/comments", handler.GetBottleComments)
			bottles.GET("/my", handler.GetMyBottles)
			bottles.GET("/picked", handler.GetPickedBottles)
			bottles.GET("/count", handler.GetBottleCount)
			bottles.POST("/watch-ad", handler.WatchAd)
			bottles.POST("/:id/throw-back", handler.ThrowBackBottle)
			bottles.DELETE("/:id", handler.DeleteBottle)
		}

		// 会话
		conversations := api.Group("/conversations")
		{
			conversations.GET("", handler.GetConversations)
			conversations.DELETE("/:id", handler.DeleteConversation)
			conversations.GET("/:id", handler.GetConversationDetail)
			conversations.GET("/:id/messages", handler.GetMessages)
			conversations.POST("/:id/messages", handler.SendMessage)
		}

		// 举报
		reports := api.Group("/reports")
		{
			reports.POST("", handler.CreateReport)
		}

		// 黑名单
		blacklists := api.Group("/blacklists")
		{
			blacklists.POST("", handler.AddBlacklist)
			blacklists.DELETE("/:id", handler.RemoveBlacklist)
			blacklists.GET("", handler.GetBlacklist)
		}

		// 签到
		checkin := api.Group("/checkin")
		{
			checkin.POST("", handler.Checkin)
			checkin.GET("/status", handler.GetCheckinStatus)
		}

		// 上传
		api.POST("/upload", handler.UploadFile)
	}

	r.Run(":8080")
}
