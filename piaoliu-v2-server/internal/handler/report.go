package handler

import (
	"drift-bottle-server/internal/models"
	"drift-bottle-server/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateReportRequest struct {
	TargetUserID    int64  `json:"targetUserId" binding:"required"`
	TargetBottleID  *int64 `json:"targetBottleId"`
	TargetMessageID *int64 `json:"targetMessageId"`
	Reason          int8   `json:"reason" binding:"required,min=1,max=5"`
	Description     string `json:"description"`
}

func CreateReport(c *gin.Context) {
	userID := c.GetInt64("userId")

	var req CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	report := models.Report{
		ReporterID:      userID,
		TargetUserID:    req.TargetUserID,
		TargetBottleID:  req.TargetBottleID,
		TargetMessageID: req.TargetMessageID,
		Reason:          req.Reason,
		Description:     req.Description,
		Status:          0,
	}

	if err := models.DB.Create(&report).Error; err != nil {
		response.Fail(c, 500, "提交举报失败")
		return
	}

	response.Success(c, gin.H{"id": report.ID})
}

type AddBlacklistRequest struct {
	BlockedUserID int64 `json:"blockedUserId" binding:"required"`
}

func AddBlacklist(c *gin.Context) {
	userID := c.GetInt64("userId")

	var req AddBlacklistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	if req.BlockedUserID == userID {
		response.Fail(c, 400, "不能拉黑自己")
		return
	}

	blacklist := models.Blacklist{
		UserID:        userID,
		BlockedUserID: req.BlockedUserID,
	}

	// 使用 FirstOrCreate 避免重复
	result := models.DB.Where("user_id = ? AND blocked_user_id = ?", userID, req.BlockedUserID).First(&blacklist)
	if result.Error == nil {
		response.Fail(c, 400, "已在黑名单中")
		return
	}

	if err := models.DB.Create(&blacklist).Error; err != nil {
		response.Fail(c, 500, "操作失败")
		return
	}

	response.Success(c, nil)
}

func RemoveBlacklist(c *gin.Context) {
	userID := c.GetInt64("userId")
	blockedID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := models.DB.Where("user_id = ? AND blocked_user_id = ?", userID, blockedID).Delete(&models.Blacklist{}).Error; err != nil {
		response.Fail(c, 500, "操作失败")
		return
	}

	response.Success(c, nil)
}

func GetBlacklist(c *gin.Context) {
	userID := c.GetInt64("userId")

	var blacklists []models.Blacklist
	models.DB.Where("user_id = ?", userID).Find(&blacklists)

	// 获取被拉黑用户信息
	type BlacklistItem struct {
		models.Blacklist
		BlockedUser *models.User `json:"blockedUser"`
	}

	var result []BlacklistItem
	for _, bl := range blacklists {
		var user models.User
		models.DB.Select("id, nickname, avatar_url").First(&user, bl.BlockedUserID)
		result = append(result, BlacklistItem{
			Blacklist:   bl,
			BlockedUser: &user,
		})
	}

	response.Success(c, result)
}
