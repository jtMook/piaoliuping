package handler

import (
	"drift-bottle-server/internal/models"
	"drift-bottle-server/pkg/response"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	userID := c.GetInt64("userId")

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		response.Fail(c, 404, "用户不存在")
		return
	}

	response.Success(c, user)
}

type UpdateProfileRequest struct {
	Nickname     string   `json:"nickname"`
	AvatarURL    string   `json:"avatarUrl"`
	Gender       int8     `json:"gender"`
	Birthday     string   `json:"birthday"`
	InterestTags []string `json:"interestTags"`
}

func UpdateProfile(c *gin.Context) {
	userID := c.GetInt64("userId")

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.AvatarURL != "" {
		updates["avatar_url"] = req.AvatarURL
	}
	if req.Gender > 0 {
		updates["gender"] = req.Gender
	}
	if req.Birthday != "" {
		updates["birthday"] = req.Birthday
	}
	if req.InterestTags != nil {
		updates["interest_tags"] = models.StringArray(req.InterestTags)
	}

	if len(updates) == 0 {
		response.Fail(c, 400, "没有要更新的字段")
		return
	}

	if err := models.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		response.Fail(c, 500, "更新失败")
		return
	}

	var user models.User
	models.DB.First(&user, userID)
	response.Success(c, user)
}

type UpdateAddressRequest struct {
	ProvinceCode string `json:"provinceCode" binding:"required"`
	ProvinceName string `json:"provinceName" binding:"required"`
	CityCode     string `json:"cityCode" binding:"required"`
	CityName     string `json:"cityName" binding:"required"`
	DistrictCode string `json:"districtCode"`
	DistrictName string `json:"districtName"`
}

func UpdateAddress(c *gin.Context) {
	userID := c.GetInt64("userId")
	fmt.Printf("[UpdateAddress] userID: %d\n", userID)

	var req UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("[UpdateAddress] 参数绑定失败: %v\n", err)
		response.Fail(c, 400, "参数错误")
		return
	}
	fmt.Printf("[UpdateAddress] 收到请求: province=%s, city=%s\n", req.ProvinceName, req.CityName)

	now := time.Now()
	updates := map[string]interface{}{
		"province_code":      req.ProvinceCode,
		"province_name":      req.ProvinceName,
		"city_code":          req.CityCode,
		"city_name":          req.CityName,
		"district_code":      req.DistrictCode,
		"district_name":      req.DistrictName,
		"address_updated_at": now,
	}

	if err := models.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		fmt.Printf("[UpdateAddress] 数据库更新失败: %v\n", err)
		response.Fail(c, 500, "更新失败")
		return
	}
	fmt.Printf("[UpdateAddress] 数据库更新成功\n")

	var updatedUser models.User
	models.DB.First(&updatedUser, userID)
	fmt.Printf("[UpdateAddress] 返回用户数据: cityName=%s\n", updatedUser.CityName)
	response.Success(c, updatedUser)
}
