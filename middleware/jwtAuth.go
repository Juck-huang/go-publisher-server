package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"hy.juck.com/go-publisher-server/config"
	"hy.juck.com/go-publisher-server/utils"
)

var (
	G = config.G
)

// AuthJwtToken gin框架进行token认证
func AuthJwtToken() func(c *gin.Context) {
	return func(c *gin.Context) {
		fmt.Println("request", c.Request.URL.Path)
		requestUrl := c.Request.URL.Path
		// 判断请求的url是不是在白名单列表中，在，则放行
		for _, url := range G.C.Jwt.WhiteUrlList {
			split := strings.Split(requestUrl, G.C.Server.Path)
			if split[1] == url {
				c.Next()
				return
			}
		}
		token := c.Request.Header.Get("x-token")
		// token为空不通过
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
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
			G.Logger.Errorf("token解析错误或已经失效:[%s]", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "token解析错误或token已失效",
				"result":  []string{},
			})
			c.Abort()
			return
		}
		c.Set("username", claims.Username)
		// 以下解析该用户是否已经登出且token已为失效
		err = CheckLogoutRedis(claims.Username, token)
		if err != nil {
			c.Abort()
			return
		}
		c.Next()
	}
}

func CheckLogoutRedis(username string, token string) error {
	// 以用户名+过期时间+token作为key
	redisKey := fmt.Sprintf("%s_%d_%s", username, G.C.Jwt.Token.Expire, token)
	// 从redis获取登出信息，如果有，则说明已经登出不让继续使用，如果没有则继续使用
	tokenRedis, err := G.RedisClient.Get(redisKey).Result()
	switch {
	case err == redis.Nil:
		// 从redis中未获取到token，身份通过
		//G.Logger.Infof("未从redis中获取到登出登记token,key为[%s]，信息：[%v]", redisKey, err)
		G.Logger.Infof("[%s]身份验证通过", username)
		return nil
	case err != nil:
		// 获取token错误处理
		G.Logger.Errorf("校验token失败,key为[%s]，错误信息：[%v]", redisKey, err)
		return errors.New("token解析错误或token已失效")
	case tokenRedis != "":
		// 从redis中获取到token，则说明已经登出不让继续使用
		G.Logger.Infof("已从redis中获取到登出登记token:[%s]", tokenRedis)
		return errors.New("token解析错误或token已失效")
	}
	return nil
}
