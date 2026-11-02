package handler

import (
	"drift-bottle-server/pkg/response"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, 400, "请选择文件")
		return
	}

	// 校验文件大小（10MB）
	if file.Size > 10*1024*1024 {
		response.Fail(c, 400, "文件大小不能超过10MB")
		return
	}

	// 校验文件类型
	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".mp3":  true,
		".wav":  true,
		".m4a":  true,
	}
	if !allowedExts[ext] {
		response.Fail(c, 400, "不支持的文件类型")
		return
	}

	// 生成唯一文件名
	datePath := time.Now().Format("2006/01/02")
	filename := fmt.Sprintf("%s/%s%s",
		datePath,
		uuid.New().String()[:8]+time.Now().Format("150405"),
		ext,
	)

	// 确保上传目录存在
	uploadDir := "./uploads/" + datePath
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.Fail(c, 500, "上传目录创建失败")
		return
	}

	// 本地保存文件
	dst := "uploads/" + filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		response.Fail(c, 500, "文件保存失败")
		return
	}

	// 返回可访问的URL（本地开发用）
	fileURL := fmt.Sprintf("http://localhost:8080/uploads/%s", filename)

	response.Success(c, gin.H{
		"url": fileURL,
	})
}
