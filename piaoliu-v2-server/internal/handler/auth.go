package handler

import (
	"drift-bottle-server/internal/middleware"
	"drift-bottle-server/internal/models"
	"drift-bottle-server/pkg/response"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Code string `json:"code" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 调用微信 code2Session 接口获取 openid
	// 注意：实际项目需要配置正确的 appid 和 secret
	// 这里使用 mock 方式
	var openID string
	unionID := ""

	// 测试环境：将真实微信 code 映射到测试账号
	// 生产环境需要替换为真实的微信 code2Session 调用
	if req.Code == "dev_test_code_a" {
		openID = "mock_openid_dev_test_code_a"
	} else if req.Code == "dev_test_code_b" {
		openID = "mock_openid_dev_test_code_b"
	} else {
		// 真实微信 code，生成固定的 openid（实际应该调用微信接口）
		openID = fmt.Sprintf("mock_openid_%s", req.Code)
	}

	// 实际调用微信接口（取消注释使用）
	/*
		wxURL := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
			config.LoadConfig().WxAppID, config.LoadConfig().WxAppSecret, req.Code)
		wxResp, err := http.Get(wxURL)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "微信登录失败")
			return
		}
		defer wxResp.Body.Close()
		body, _ := io.ReadAll(wxResp.Body)
		var wxSession WxSessionResponse
		json.Unmarshal(body, &wxSession)
		if wxSession.ErrCode != 0 {
			response.Fail(c, http.StatusInternalServerError, fmt.Sprintf("微信登录失败: %s", wxSession.ErrMsg))
			return
		}
		openID = wxSession.OpenID
		unionID = wxSession.UnionID
	*/

	// 查找或创建用户
	var user models.User
	result := models.DB.Where("openid = ?", openID).First(&user)

	if result.Error != nil {
		// 新用户，创建临时记录（未完善资料）
		user = models.User{
			OpenID:   openID,
			UnionID:  unionID,
			Nickname: fmt.Sprintf("海边旅人#%d", generateRandomID()),
			Gender:   0, // 未知，待完善
			Birthday: "",
			CityCode: "",
			CityName: "",
			Status:   1,
		}
		if err := models.DB.Create(&user).Error; err != nil {
			response.Fail(c, http.StatusInternalServerError, "创建用户失败")
			return
		}
	}

	// 检查用户状态
	if user.Status == 3 {
		response.Fail(c, http.StatusForbidden, "账号已被封禁")
		return
	}

	// 生成 token
	token, err := middleware.GenerateToken(user.ID, user.OpenID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "生成 token 失败")
		return
	}

	// 检查资料是否完善
	profileComplete := user.Gender != 0 && user.Birthday != "" && user.CityCode != ""

	response.Success(c, gin.H{
		"token":           token,
		"user":            user,
		"profileComplete": profileComplete,
	})
}

func generateRandomID() int {
	var count int64
	models.DB.Model(&models.User{}).Count(&count)
	return int(count) + 1000
}
