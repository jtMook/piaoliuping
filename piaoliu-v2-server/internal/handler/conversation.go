package handler

import (
	"drift-bottle-server/internal/models"
	"drift-bottle-server/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetConversations(c *gin.Context) {
	userID := c.GetInt64("userId")

	var conversations []models.Conversation
	models.DB.
		Preload("UserA", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, nickname, avatar_url, city_name")
		}).
		Preload("UserB", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, nickname, avatar_url, city_name")
		}).
		Preload("Bottle", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, content, bottle_type")
		}).
		Preload("LastMsg").
		// 排除本人已删除（隐藏）的会话
	Where("(user_a_id = ? AND deleted_by_a = false) OR (user_b_id = ? AND deleted_by_b = false)", userID, userID).
		Order("last_message_at DESC").
		Find(&conversations)

	// 计算未读数并组装返回数据
	type ConvResponse struct {
		models.Conversation
		UnreadCount int          `json:"unreadCount"`
		OtherUser   *models.User `json:"otherUser"`
	}

	var result []ConvResponse
	for _, conv := range conversations {
		var unreadCount int64
		models.DB.Model(&models.Message{}).
			Where("conversation_id = ? AND sender_id != ? AND is_read = ?", conv.ID, userID, false).
			Count(&unreadCount)

		otherUser := conv.UserB
		if conv.UserBID == userID {
			otherUser = conv.UserA
		}

		result = append(result, ConvResponse{
			Conversation: conv,
			UnreadCount:  int(unreadCount),
			OtherUser:    otherUser,
		})
	}

	response.Success(c, result)
}

// DeleteConversation 删除会话：仅对当前用户隐藏，对方仍可见；对方再来新消息时恢复显示
func DeleteConversation(c *gin.Context) {
	userID := c.GetInt64("userId")
	convID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var conv models.Conversation
	if err := models.DB.First(&conv, convID).Error; err != nil {
		response.Fail(c, 404, "会话不存在")
		return
	}
	if conv.UserAID != userID && conv.UserBID != userID {
		response.Fail(c, 403, "无权访问该会话")
		return
	}

	// 只标记当前用户一侧
	if conv.UserAID == userID {
		models.DB.Model(&conv).UpdateColumn("deleted_by_a", true)
	} else {
		models.DB.Model(&conv).UpdateColumn("deleted_by_b", true)
	}

	response.Success(c, gin.H{"message": "已删除"})
}

func GetConversationDetail(c *gin.Context) {
	userID := c.GetInt64("userId")
	convID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var conv models.Conversation
	if err := models.DB.
		Preload("UserA", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, nickname, avatar_url, city_name, gender")
		}).
		Preload("UserB", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, nickname, avatar_url, city_name, gender")
		}).
		Preload("Bottle", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, content, bottle_type")
		}).
		First(&conv, convID).Error; err != nil {
		response.Fail(c, 404, "会话不存在")
		return
	}

	if conv.UserAID != userID && conv.UserBID != userID {
		response.Fail(c, 403, "无权访问该会话")
		return
	}

	otherUser := conv.UserB
	if conv.UserBID == userID {
		otherUser = conv.UserA
	}

	response.Success(c, gin.H{
		"id":        conv.ID,
		"bottleId":  conv.BottleID,
		"otherUser": otherUser,
		"bottle":    conv.Bottle,
		"isClosed":  conv.IsClosed,
	})
}

func GetMessages(c *gin.Context) {
	userID := c.GetInt64("userId")
	convID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	// 验证用户是否属于该会话
	var conv models.Conversation
	if err := models.DB.First(&conv, convID).Error; err != nil {
		response.Fail(c, 404, "会话不存在")
		return
	}
	if conv.UserAID != userID && conv.UserBID != userID {
		response.Fail(c, 403, "无权访问该会话")
		return
	}

	var messages []models.Message
	models.DB.
		Preload("Sender", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, nickname, avatar_url")
		}).
		Where("conversation_id = ?", convID).
		Order("created_at ASC").
		Find(&messages)

	// 标记已读
	models.DB.Model(&models.Message{}).
		Where("conversation_id = ? AND sender_id != ?", convID, userID).
		Update("is_read", true)

	response.Success(c, messages)
}

type SendMessageRequest struct {
	Content  string `json:"content"`
	MediaURL string `json:"mediaUrl"`
	MsgType  int8   `json:"msgType"`
}

func SendMessage(c *gin.Context) {
	userID := c.GetInt64("userId")
	convID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	// 验证用户是否属于该会话
	var conv models.Conversation
	if err := models.DB.First(&conv, convID).Error; err != nil {
		response.Fail(c, 404, "会话不存在")
		return
	}
	if conv.UserAID != userID && conv.UserBID != userID {
		response.Fail(c, 403, "无权访问该会话")
		return
	}

	msgType := req.MsgType
	if msgType == 0 {
		msgType = 1
	}

	message := models.Message{
		ConversationID: convID,
		SenderID:       userID,
		Content:        req.Content,
		MediaURL:       req.MediaURL,
		MsgType:        msgType,
	}

	if err := models.DB.Create(&message).Error; err != nil {
		response.Fail(c, 500, "发送失败")
		return
	}

	// 更新会话最后消息；若任一方曾删除（隐藏），新消息到来时恢复显示
	models.DB.Model(&conv).Updates(map[string]interface{}{
		"last_message_id": message.ID,
		"last_message_at": message.CreatedAt,
		"deleted_by_a":    false,
		"deleted_by_b":    false,
	})

	// TODO: 通过 WebSocket 推送消息给接收方

	response.Success(c, message)
}
