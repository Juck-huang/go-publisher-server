package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"hy.juck.com/go-publisher-server/model/auth"
)

// WhiteAuth 白名单认证
func WhiteAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		if !G.C.White.Status {
			c.Next()
			return
		}
		// 获取客户端公网ip
		host := c.ClientIP()
		var exist bool
		// 1.首先校验配置文件中的whiteIpList中的ip
		for _, ip := range G.C.White.WhiteIpList {
			if ip == host {
				exist = true
				break
			}
		}

		// 2.再校验数据库中的access_ip_white中的ip列表
		var accessIpWhite []auth.AccessIpWhite
		G.DB.Debug().Find(&accessIpWhite)
		var ipList []string
		for _, white := range accessIpWhite {
			ipList = append(ipList, strings.Split(white.IpList, ",")...)
		}
		for _, ip := range ipList {
			if ip == host {
				exist = true
				break
			}
		}
		if exist {
			G.Logger.Infof("Host[%s]白名单校验通过", host)
			c.Next()
		} else {
			fmt.Println("iplist", ipList)
			G.Logger.Errorf("Host[%s]未在白名单列表授权", host)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "非法访问，请联系管理员授权",
			})
			c.Abort()
			return
		}
	}
}
