package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// WhiteAuth 白名单认证
func WhiteAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		if !G.C.White.Status {
			c.Next()
			return
		}
		host := strings.Split(c.Request.RemoteAddr, ":")[0]
		var exist bool
		for _, ip := range G.C.White.WhiteIpList {
			if ip == host {
				exist = true
				break
			}
		}
		if exist {
			G.Logger.Infof("Host[%s]白名单校验通过", host)
			c.Next()
		} else {
			G.Logger.Errorf("Host[%s]未在白名单列表授权", host)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "系统认证失败",
			})
			c.Abort()
			return
		}
	}
}
