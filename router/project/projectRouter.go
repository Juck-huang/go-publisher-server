package project

import (
	"github.com/gin-gonic/gin"
	"hy.juck.com/go-publisher-server/service"
	"net/http"
)

// ProjectList 获取项目列表
func ProjectList(c *gin.Context) {
	projectService := service.NewProjectService()
	projects := projectService.GetProjectList()
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "获取项目列表数据成功",
		"result":  projects,
	})
}
