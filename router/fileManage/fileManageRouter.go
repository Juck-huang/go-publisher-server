package fileManage

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// GetFileContent 获取文件内容
func GetFileContent(c *gin.Context) {
	file, _ := os.ReadFile("temp/config.json")
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"success": true,
		"msg":     string(file),
	})
}
