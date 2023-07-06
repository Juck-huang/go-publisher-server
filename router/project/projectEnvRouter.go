package project

import (
	"github.com/gin-gonic/gin"
	"hy.juck.com/go-publisher-server/service"
	"net/http"
)

// ProjectEnvList 通过项目id获取项目类型列表
func ProjectEnvList(c *gin.Context) {
	var projectT struct {
		ProjectId int64 `json:"projectId"`
	}
	err := c.ShouldBindJSON(&projectT)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "参数解析错误",
			"result":  []string{},
		})
		return
	}
	if projectT.ProjectId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "项目id不能为空",
			"result":  []string{},
		})
		return
	}
	projectEnvService := service.NewProjectEnvService()
	projectEnvs := projectEnvService.GetProjectEnvList(projectT.ProjectId)
	var result any
	if len(projectEnvs) > 0 {
		result = projectEnvs
	} else {
		result = []string{}
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "获取项目环境列表数据成功",
		"result":  result,
	})
	return
}
