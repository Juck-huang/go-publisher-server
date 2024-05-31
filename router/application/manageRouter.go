package application

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"hy.juck.com/go-publisher-server/service"
)

func AppList(c *gin.Context) {
	appManageService := service.NewAppManageService()
	appStatusList := appManageService.GetAppStatusList()
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "获取最新应用信息成功",
		"success": true,
		"result":  appStatusList,
	})
}

func StartOrStopApp(c *gin.Context) {
	var err error
	var message string
	var requestMap = make(map[string]any, 1)
	err = c.ShouldBindJSON(&requestMap)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "参数解析错误",
			"success": false,
			"result":  []string{},
		})
		return
	}
	id := requestMap["id"]
	direct := requestMap["direct"]
	if id == nil || direct == nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "参数缺失",
			"success": false,
			"result":  []string{},
		})
		return
	}
	switch direct {
	case "start":
		idUint, _ := strconv.ParseUint(id.(string), 0, 0)
		appManageService := service.NewAppManageService()
		err = appManageService.StartOrStopApp(uint(idUint), "start")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": err.Error(),
				"success": false,
				"result":  []string{},
			})
			return
		}
		message = "开启应用成功"
	case "stop":
		idUint, _ := strconv.ParseUint(id.(string), 0, 0)
		appManageService := service.NewAppManageService()
		err = appManageService.StartOrStopApp(uint(idUint), "stop")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": err.Error(),
				"success": false,
				"result":  []string{},
			})
			return
		}
		message = "停止应用成功"
	case "restart":
		idUint, _ := strconv.ParseUint(id.(string), 0, 0)
		appManageService := service.NewAppManageService()
		err = appManageService.RestartApp(uint(idUint))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": err.Error(),
				"success": false,
				"result":  []string{},
			})
			return
		}
		message = "重启应用成功"
	default:
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "命令输入有误",
			"success": false,
			"result":  []string{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": message,
		"success": true,
		"result":  []string{},
	})
}
