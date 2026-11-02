package handler

import (
	"drift-bottle-server/internal/middleware"
	"drift-bottle-server/internal/models"
	"drift-bottle-server/pkg/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Checkin(c *gin.Context) {
	userID := middleware.GetUserID(c)

	today := time.Now().Format("2006-01-02")

	// 检查今天是否已签到
	var count int64
	models.DB.Model(&models.CheckinRecord{}).
		Where("user_id = ? AND DATE(created_at) = ?", userID, today).
		Count(&count)

	if count > 0 {
		response.Fail(c, http.StatusBadRequest, "今天已经签到过了")
		return
	}

	// 计算连续签到天数
	var lastCheckin models.CheckinRecord
	err := models.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&lastCheckin).Error

	consecutiveDays := 1
	if err == nil {
		lastDate := lastCheckin.CreatedAt.Format("2006-01-02")
		yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		if lastDate == yesterday {
			// 7天后重置为第1天，循环签到
			if lastCheckin.ConsecutiveDays >= 7 {
				consecutiveDays = 1
			} else {
				consecutiveDays = int(lastCheckin.ConsecutiveDays) + 1
			}
		}
	}

	// 计算奖励瓶子数（7天循环：10, 10, 30, 10, 10, 10, 100）
	rewards := []int{10, 10, 30, 10, 10, 10, 100}
	dayIndex := consecutiveDays - 1
	bottleReward := rewards[dayIndex]

	// 创建签到记录
	record := models.CheckinRecord{
		UserID:          userID,
		ConsecutiveDays: int64(consecutiveDays),
		BottleReward:    int64(bottleReward),
	}
	if err := models.DB.Create(&record).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "签到失败")
		return
	}

	// 增加瓶子余额
	models.DB.Model(&models.User{}).Where("id = ?", userID).UpdateColumn("bottle_count", gorm.Expr("bottle_count + ?", bottleReward))

	response.Success(c, gin.H{
		"consecutiveDays": consecutiveDays,
		"bottleReward":    bottleReward,
	})
}

func GetCheckinStatus(c *gin.Context) {
	userID := middleware.GetUserID(c)

	today := time.Now().Format("2006-01-02")

	var count int64
	models.DB.Model(&models.CheckinRecord{}).
		Where("user_id = ? AND DATE(created_at) = ?", userID, today).
		Count(&count)

	// 获取最近7天的签到记录
	var records []models.CheckinRecord
	models.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(7).
		Find(&records)

	var lastConsecutiveDays int64
	if len(records) > 0 {
		lastConsecutiveDays = records[0].ConsecutiveDays
	}

	response.Success(c, gin.H{
		"checkedInToday":  count > 0,
		"consecutiveDays": lastConsecutiveDays,
	})
}
