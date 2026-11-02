package middleware

import (
	"drift-bottle-server/internal/config"
	"drift-bottle-server/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var cfg = config.LoadConfig()

type Claims struct {
	UserID int64  `json:"userId"`
	OpenID string `json:"openId"`
	jwt.RegisteredClaims
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未提供认证token")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			response.Unauthorized(c, "token格式错误")
			c.Abort()
			return
		}

		claims, err := ParseToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "token无效或已过期")
			c.Abort()
			return
		}

		c.Set("userId", claims.UserID)
		c.Set("openId", claims.OpenID)
		c.Next()
	}
}

func GenerateToken(userID int64, openID string) (string, error) {
	claims := Claims{
		UserID: userID,
		OpenID: openID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenMalformed
}

func GetUserID(c *gin.Context) int64 {
	userID, _ := c.Get("userId")
	return userID.(int64)
}
