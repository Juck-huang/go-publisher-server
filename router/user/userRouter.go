package user

import (
	"github.com/gin-gonic/gin"
	"hy.juck.com/go-publisher-server/dto/login"
	"hy.juck.com/go-publisher-server/service"
	"hy.juck.com/go-publisher-server/utils"
	"net/http"
)

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "获取用户数据成功",
		"success": true,
		"result":  map[string]any{},
	})
}

// Logout 登出
func Logout(c *gin.Context) {
	// 登出逻辑
	// 1.登出后把该用户的token加入redis，并设置过期时间和token的有效时间一致
	// 2.请求接口时，校验token通过后，先去redis中读取该token，如果存在，则证明该用户已经登出，直接返回没有权限
	//（存入redis的token有效时间永远比生成的token过期时间长），需要重新登陆获取token
	token := c.GetHeader("x-token")
	username, exist := c.Get("username")
	if !exist {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "登出失败",
			"success": false,
		})
		return
	}
	userService := service.NewUserService()
	err := userService.Logout(token, username.(string))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "登出失败:" + err.Error(),
			"success": false,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "登出成功",
		"success": true,
	})
}

// Login 登录
func Login(c *gin.Context) {
	var requestDto login.RequestDto
	err := c.ShouldBindJSON(&requestDto)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"result":  map[string]any{},
			"message": err.Error(),
		})
		return
	}
	var userService = service.NewUserService()
	err = userService.CheckUsernameAndPassword(requestDto.Username, requestDto.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"result":  map[string]any{},
			"message": err.Error(),
		})
		return
	}
	token, err := utils.GenToken(requestDto.Username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"result":  map[string]any{},
			"message": err.Error(),
		})
		return
	}
	// 登录接口
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"success": true,
		"result": map[string]any{
			"token": token,
		},
		"message": "登录成功",
	})
}
