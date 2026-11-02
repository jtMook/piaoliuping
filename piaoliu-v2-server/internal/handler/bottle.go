package handler

import (
	"drift-bottle-server/internal/models"
	"drift-bottle-server/pkg/response"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateBottleRequest struct {
	Content      string   `json:"content" binding:"required"`
	MediaUrls    []string `json:"mediaUrls"`
	Tags         []string `json:"tags"`
	BottleType   int8     `json:"bottleType"`
	VisibleScope int8     `json:"visibleScope"`
}

func CreateBottle(c *gin.Context) {
	userID := c.GetInt64("userId")

	var req CreateBottleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	// 校验内容长度
	if len(req.Content) > 500 {
		response.Fail(c, 400, "内容不能超过500字")
		return
	}

	// 校验图片数量
	if len(req.MediaUrls) > 3 {
		response.Fail(c, 400, "最多只能上传3张图片")
		return
	}

	// 获取用户信息
	var user models.User
	models.DB.First(&user, userID)

	// 检查瓶子余额
	if user.BottleCount <= 0 {
		response.Fail(c, 400, "瓶子余额不足，请先获取瓶子")
		return
	}

	bottleType := req.BottleType
	if bottleType == 0 {
		bottleType = 1
	}
	visibleScope := req.VisibleScope
	if visibleScope == 0 {
		visibleScope = 1
	}

	bottle := models.Bottle{
		UserID:       userID,
		Content:      req.Content,
		MediaUrls:    models.StringArray(req.MediaUrls),
		Tags:         models.StringArray(req.Tags),
		BottleType:   bottleType,
		ProvinceCode: user.ProvinceCode,
		CityCode:     user.CityCode,
		DistrictCode: user.DistrictCode,
		Status:       1, // 已通过（开发环境自动审核）
		VisibleScope: visibleScope,
	}

	// 秘密瓶设置过期时间 24h
	if bottleType == 3 {
		expireAt := time.Now().Add(24 * time.Hour)
		bottle.ExpireAt = &expireAt
	}

	if err := models.DB.Create(&bottle).Error; err != nil {
		response.Fail(c, 500, "创建瓶子失败")
		return
	}

	// 扣减瓶子余额
	models.DB.Model(&models.User{}).Where("id = ?", userID).UpdateColumn("bottle_count", user.BottleCount-1)

	response.SuccessWithMessage(c, "瓶子已投入海中，审核通过后即可被捡到", gin.H{
		"id":     bottle.ID,
		"status": bottle.Status,
	})
}

func PickBottle(c *gin.Context) {
	userID := c.GetInt64("userId")

	// 获取筛选参数
	genderFilter := c.Query("gender")

	// 获取用户信息（含兴趣标签）
	var user models.User
	models.DB.First(&user, userID)

	// 检查瓶子余额
	if user.BottleCount <= 0 {
		response.Fail(c, 400, "瓶子余额不足，请先获取瓶子")
		return
	}

	// 先查出用户已捡过的瓶子 ID
	var pickedIDs []int64
	models.DB.Table("bottle_picked_records").Where("user_id = ?", userID).Pluck("bottle_id", &pickedIDs)

	// 7天过期时间
	expireTime := time.Now().AddDate(0, 0, -7)

	// 基础查询：已通过、非自己扔的、未捡过的、7天内
	baseQuery := func() *gorm.DB {
		q := models.DB.Where("status = ? AND user_id != ? AND created_at > ?", 1, userID, expireTime)
		if len(pickedIDs) > 0 {
			q = q.Where("id NOT IN ?", pickedIDs)
		}
		if genderFilter != "" {
			genderVal, _ := strconv.ParseInt(genderFilter, 10, 64)
			q = q.Where("user_id IN (?)", models.DB.Table("users").Select("id").Where("gender = ?", genderVal))
		}
		return q
	}

	// 第一阶段：收集候选瓶子（按地域优先级）
	var candidates []models.Bottle

	// 同城优先
	baseQuery().Where("city_code = ?", user.CityCode).
		Order("created_at DESC").Limit(50).Find(&candidates)

	// 同城不够，补同省
	if len(candidates) < 50 {
		var provinceBottles []models.Bottle
		baseQuery().Where("city_code != ? AND province_code = ?", user.CityCode, user.ProvinceCode).
			Order("created_at DESC").Limit(50 - len(candidates)).Find(&provinceBottles)
		candidates = append(candidates, provinceBottles...)
	}

	// 还不够，全国兜底
	if len(candidates) < 50 {
		var nationalBottles []models.Bottle
		baseQuery().Where("province_code != ?", user.ProvinceCode).
			Order("created_at DESC").Limit(50 - len(candidates)).Find(&nationalBottles)
		candidates = append(candidates, nationalBottles...)
	}

	if len(candidates) == 0 {
		response.Fail(c, 404, "附近没有瓶子了，稍后再试")
		return
	}

	// 第二阶段：基于兴趣标签 + 时间随机算法打分
	userTags := make(map[string]bool)
	for _, tag := range user.InterestTags {
		userTags[tag] = true
	}

	now := time.Now()
	type scoredBottle struct {
		bottle models.Bottle
		score  float64
	}
	var scored []scoredBottle

	for _, b := range candidates {
		score := 0.0

		// 1. 兴趣标签匹配：每个匹配标签 +30 分
		for _, tag := range b.Tags {
			if userTags[tag] {
				score += 30
			}
		}

		// 2. 时间权重：越新的瓶子分数越高（24小时内满分，之后递减）
		hoursSince := now.Sub(b.CreatedAt).Hours()
		timeScore := 20.0
		if hoursSince < 24 {
			timeScore = 20.0 * (1 - hoursSince/24)
		} else if hoursSince < 72 {
			timeScore = 5.0
		} else {
			timeScore = 1.0
		}
		score += timeScore

		// 3. 随机扰动：±10 分，避免频繁捞到重复内容
		score += rand.Float64()*20 - 10

		scored = append(scored, scoredBottle{bottle: b, score: score})
	}

	// 按分数排序，选最高分
	var best models.Bottle
	bestScore := -999.0
	for _, s := range scored {
		if s.score > bestScore {
			bestScore = s.score
			best = s.bottle
		}
	}

	bottle := best

	// 创建捡瓶子记录
	record := models.BottlePickedRecord{
		BottleID: bottle.ID,
		UserID:   userID,
	}
	models.DB.Create(&record)

	// 更新瓶子被捡次数（不再设置 is_picked，允许多人捡）
	models.DB.Model(&bottle).UpdateColumn("picked_count", bottle.PickedCount+1)

	// 扣减捡瓶人的瓶子余额
	models.DB.Model(&models.User{}).Where("id = ?", userID).UpdateColumn("bottle_count", gorm.Expr("bottle_count - 1"))

	// 获取瓶子发布者信息
	var bottleOwner models.User
	models.DB.Select("id, nickname, avatar_url, province_name, city_name, district_name, gender").First(&bottleOwner, bottle.UserID)

	nickname := bottleOwner.Nickname
	if nickname == "" {
		nickname = "匿名用户"
	}
	provinceName := bottleOwner.ProvinceName
	cityName := bottleOwner.CityName
	if provinceName == "" && cityName == "" {
		provinceName = "未知"
	}
	districtName := bottleOwner.DistrictName
	gender := bottleOwner.Gender

	response.Success(c, gin.H{
		"id":           bottle.ID,
		"content":      bottle.Content,
		"mediaUrls":    bottle.MediaUrls,
		"tags":         bottle.Tags,
		"bottleType":   bottle.BottleType,
		"provinceName": provinceName,
		"cityName":     cityName,
		"districtName": districtName,
		"nickname":     nickname,
		"gender":       gender,
		"createdAt":    bottle.CreatedAt,
		"canReply":     true,
	})
}

func GetBottleDetail(c *gin.Context) {
	userID := c.GetInt64("userId")
	bottleID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var bottle models.Bottle
	if err := models.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, nickname, avatar_url, gender, birthday, city_name, district_name")
	}).First(&bottle, bottleID).Error; err != nil {
		response.Fail(c, 404, "瓶子不存在")
		return
	}

	// 已删除的瓶子对所有人不可见
	if bottle.Status == 3 {
		response.Fail(c, 404, "瓶子不存在或已删除")
		return
	}

	// 从 User 获取 cityName
	if bottle.User != nil {
		bottle.CityName = bottle.User.CityName
	}

	// 检查是否是自己扔的瓶子
	isOwner := bottle.UserID == userID

	// 检查是否有他人回复过（即有其他用户捡过并回复了）
	var replyCount int64
	models.DB.Table("bottle_picked_records").
		Where("bottle_id = ? AND user_id != ? AND is_replied = ?", bottleID, userID, true).
		Count(&replyCount)

	// 检查是否已捡过
	var record models.BottlePickedRecord
	hasPicked := models.DB.Where("bottle_id = ? AND user_id = ?", bottleID, userID).First(&record).Error == nil

	// 获取回复者的会话信息
	var conversationID int64 = 0
	if replyCount > 0 {
		var conv models.Conversation
		models.DB.Where("bottle_id = ? AND user_a_id = ? AND user_b_id != ?", bottleID, userID, userID).
			Or("bottle_id = ? AND user_b_id = ? AND user_a_id != ?", bottleID, userID, userID).
			First(&conv)
		conversationID = conv.ID
	}

	response.Success(c, gin.H{
		"bottle":         bottle,
		"hasPicked":      hasPicked,
		"isOwner":        isOwner,
		"hasOtherReply":  replyCount > 0,
		"conversationId": conversationID,
		"canReply":       hasPicked,
	})
}

type ReplyBottleRequest struct {
	Content string `json:"content" binding:"required"`
}

func ReplyBottle(c *gin.Context) {
	userID := c.GetInt64("userId")
	bottleID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var req ReplyBottleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	// 获取瓶子
	var bottle models.Bottle
	if err := models.DB.First(&bottle, bottleID).Error; err != nil {
		response.Fail(c, 404, "瓶子不存在")
		return
	}

	// 创建或获取捡瓶子记录
	var record models.BottlePickedRecord
	models.DB.Where("bottle_id = ? AND user_id = ?", bottleID, userID).FirstOrCreate(&record, models.BottlePickedRecord{
		BottleID: bottleID,
		UserID:   userID,
	})

	// 标记已回复
	models.DB.Model(&record).Update("is_replied", true)

	// 查找或创建会话
	var conversation models.Conversation
	userAID := bottle.UserID // 扔瓶人
	userBID := userID        // 回复人

	result := models.DB.Where(
		"bottle_id = ? AND user_a_id = ? AND user_b_id = ?",
		bottleID, userAID, userBID,
	).First(&conversation)

	if result.Error != nil {
		conversation = models.Conversation{
			BottleID: bottleID,
			UserAID:  userAID,
			UserBID:  userBID,
		}
		models.DB.Create(&conversation)
	}

	// 发送消息
	message := models.Message{
		ConversationID: conversation.ID,
		SenderID:       userID,
		Content:        req.Content,
		MsgType:        1,
	}
	models.DB.Create(&message)

	// 更新会话最后消息
	now := time.Now()
	models.DB.Model(&conversation).Updates(map[string]interface{}{
		"last_message_id": message.ID,
		"last_message_at": now,
	})

	response.Success(c, gin.H{
		"conversationId": conversation.ID,
		"messageId":      message.ID,
	})
}

type CommentBottleRequest struct {
	Content string `json:"content" binding:"required"`
}

func CommentBottle(c *gin.Context) {
	userID := c.GetInt64("userId")
	bottleID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var req CommentBottleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	// 获取瓶子
	var bottle models.Bottle
	if err := models.DB.First(&bottle, bottleID).Error; err != nil {
		response.Fail(c, 404, "瓶子不存在")
		return
	}

	// 创建评论
	comment := models.BottleReply{
		BottleID: bottleID,
		UserID:   userID,
		Content:  req.Content,
	}
	if err := models.DB.Create(&comment).Error; err != nil {
		response.Fail(c, 500, "评论失败")
		return
	}

	// 更新捡瓶子记录为已回复
	var record models.BottlePickedRecord
	models.DB.Where("bottle_id = ? AND user_id = ?", bottleID, userID).FirstOrCreate(&record, models.BottlePickedRecord{
		BottleID: bottleID,
		UserID:   userID,
	})
	models.DB.Model(&record).Update("is_replied", true)

	response.Success(c, gin.H{
		"id": comment.ID,
	})
}

func GetBottleComments(c *gin.Context) {
	bottleID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var comments []models.BottleReply
	models.DB.Preload("User").Where("bottle_id = ?", bottleID).Order("created_at DESC").Find(&comments)

	response.Success(c, comments)
}

func GetMyBottles(c *gin.Context) {
	userID := c.GetInt64("userId")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 查出用户已捡过的瓶子 ID，从扔出列表中排除
	var pickedIDs []int64
	models.DB.Table("bottle_picked_records").Where("user_id = ?", userID).Pluck("bottle_id", &pickedIDs)

	var total int64
	query := models.DB.Model(&models.Bottle{}).Where("user_id = ? AND is_picked = ? AND status != ?", userID, false, 3)
	if len(pickedIDs) > 0 {
		query = query.Where("id NOT IN ?", pickedIDs)
	}
	query.Count(&total)

	var bottles []models.Bottle
	bQuery := models.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, nickname, avatar_url, gender, birthday, city_name")
	}).Where("user_id = ? AND is_picked = ? AND status != ?", userID, false, 3)
	if len(pickedIDs) > 0 {
		bQuery = bQuery.Where("id NOT IN ?", pickedIDs)
	}
	bQuery.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&bottles)

	response.Success(c, gin.H{
		"list":  bottles,
		"total": total,
		"page":  page,
	})
}

func GetPickedBottles(c *gin.Context) {
	userID := c.GetInt64("userId")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	type PickedBottle struct {
		models.Bottle
		PickedAt time.Time `json:"pickedAt"`
	}

	var total int64
	models.DB.Table("bottle_picked_records").
		Joins("JOIN bottles ON bottles.id = bottle_picked_records.bottle_id").
		Where("bottle_picked_records.user_id = ? AND bottles.status != 3", userID).
		Count(&total)

	var results []PickedBottle
	models.DB.Table("bottles").
		Select("bottles.*, bottle_picked_records.created_at as picked_at").
		Joins("JOIN bottle_picked_records ON bottle_picked_records.bottle_id = bottles.id").
		Where("bottle_picked_records.user_id = ? AND bottles.status != 3", userID).
		Order("bottle_picked_records.created_at DESC").
		Offset(offset).Limit(pageSize).
		Scan(&results)

	// Preload user info
	userIDs := make([]int64, 0, len(results))
	for _, r := range results {
		userIDs = append(userIDs, r.UserID)
	}
	var users []models.User
	if len(userIDs) > 0 {
		models.DB.Select("id, nickname, avatar_url, gender, birthday, city_name").Where("id IN ?", userIDs).Find(&users)
		userMap := make(map[int64]*models.User)
		for i := range users {
			userMap[users[i].ID] = &users[i]
		}
		for i := range results {
			if u, ok := userMap[results[i].UserID]; ok {
				results[i].User = u
			}
		}
	}

	response.Success(c, gin.H{
		"list":  results,
		"total": total,
		"page":  page,
	})
}

// ThrowBackBottle 把捡到的瓶子扔回海里：删除本人捡瓶记录，使其从“捞到的”列表消失
func ThrowBackBottle(c *gin.Context) {
	userID := c.GetInt64("userId")
	bottleID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	result := models.DB.Where("bottle_id = ? AND user_id = ?", bottleID, userID).Delete(&models.BottlePickedRecord{})
	if result.Error != nil {
		response.Fail(c, 500, "操作失败")
		return
	}
	if result.RowsAffected == 0 {
		response.Fail(c, 400, "该瓶子不在你捞到的列表中")
		return
	}

	response.Success(c, gin.H{"message": "已扔回海里"})
}

// DeleteBottle 删除自己扔出的瓶子（软删除 status=3 已删除）
func DeleteBottle(c *gin.Context) {
	userID := c.GetInt64("userId")
	bottleID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var bottle models.Bottle
	if err := models.DB.First(&bottle, bottleID).Error; err != nil {
		response.Fail(c, 404, "瓶子不存在")
		return
	}
	if bottle.UserID != userID {
		response.Fail(c, 403, "只能删除自己扔出的瓶子")
		return
	}

	models.DB.Model(&bottle).UpdateColumn("status", 3)

	response.Success(c, gin.H{"message": "已删除"})
}

func GetBottleCount(c *gin.Context) {
	userID := c.GetInt64("userId")

	// 获取筛选参数
	genderFilter := c.Query("gender")
	// 仅允许 1(男) 或 2(女)，防止非法/注入参数
	if genderFilter != "1" && genderFilter != "2" {
		genderFilter = ""
	}

	// 与“捡瓶子”逻辑保持一致：排除本人已捡过的瓶子
	var pickedIDs []int64
	models.DB.Table("bottle_picked_records").Where("user_id = ?", userID).Pluck("bottle_id", &pickedIDs)

	var count int64
	expireTime := time.Now().AddDate(0, 0, -7)
	// 口径与 PickBottle 完全一致：已通过、非本人扔、未捡过(本人)、7 天内、全国范围
	query := models.DB.Model(&models.Bottle{}).
		Where("status = 1 AND user_id != ? AND created_at > ?", userID, expireTime)
	if len(pickedIDs) > 0 {
		query = query.Where("id NOT IN ?", pickedIDs)
	}
	if genderFilter != "" {
		query = query.Where("user_id IN (SELECT id FROM users WHERE gender = ?)", genderFilter)
	}

	query.Count(&count)

	response.Success(c, gin.H{
		"nearbyCount": count,
	})
}

func WatchAd(c *gin.Context) {
	userID := c.GetInt64("userId")

	// 增加瓶子余额
	models.DB.Model(&models.User{}).Where("id = ?", userID).UpdateColumn("bottle_count", gorm.Expr("bottle_count + 10"))

	response.Success(c, gin.H{
		"reward": 10,
	})
}
