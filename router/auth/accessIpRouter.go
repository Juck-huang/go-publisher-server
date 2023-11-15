package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"hy.juck.com/go-publisher-server/config"
	"hy.juck.com/go-publisher-server/dto/auth"
	"hy.juck.com/go-publisher-server/service"
)

var (
	G = config.G
)

// ReceiveAuthIp 接收认证ip
func ReceiveAuthIp(c *gin.Context) {
	// 接收到认证的ip信息后，先校验认证信息是否正确，正确再存储到数据库
	var authDto auth.AuthRequest
	err := c.ShouldBindJSON(&authDto)
	if err != nil {
		G.Logger.Errorf("IP认证失败:[%s]", err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    http.StatusOK,
			"message": "IP认证失败,参数解析错误",
		})
		return
	}
	accessIpService := service.NewAccessIpService()
	err = accessIpService.CheckAuthInfo(authDto)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    http.StatusOK,
			"message": err.Error(),
		})
		return
	}
	err = accessIpService.SaveAuthIp(authDto)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    http.StatusOK,
			"message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": "IP认证成功",
		"result":  authDto,
	})
}
