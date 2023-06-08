package middleware

import (
	"github.com/gin-gonic/gin"
	"hy.juck.com/go-publisher-server/utils"
	"net/http"
)

// AuthJwtToken gin框架进行token认证
func AuthJwtToken() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("x-token")
		// token为空不通过
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "token不能为空",
				"result":  []string{},
			})
			c.Abort()
			return
		}
		// token解析错误不通过
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "token解析错误或token不正确",
				"result":  []string{},
			})
			c.Abort()
			return
		}
		c.Set("username", claims.Username)
		c.Next()
	}
}
