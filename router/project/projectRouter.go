package project

import (
	"github.com/gin-gonic/gin"
	"hy.juck.com/go-publisher-server/service"
	"net/http"
)

// ProjectList 获取项目列表
func ProjectList(c *gin.Context) {
	projectService := service.NewProjectService()
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "获取项目列表数据失败",
			"result":  []string{},
			"success": false,
		})
	}
	projects := projectService.GetProjectList(username.(string))

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "获取项目列表数据成功",
		"result":  projects,
	})
}
