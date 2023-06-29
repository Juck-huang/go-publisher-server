package router

import (
	"github.com/gin-gonic/gin"
	"hy.juck.com/go-publisher-server/service"
	"net/http"
)

// ProjectTypeList 通过项目id获取项目类型列表
func ProjectTypeList(c *gin.Context) {
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
	projectTypeService := service.NewProjectTypeService()
	projects := projectTypeService.GetByProjectId(projectT.ProjectId)
	var result any
	if len(projects) > 0 {
		result = projects
	} else {
		result = []string{}
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "获取项目类型列表数据成功",
		"result":  result,
	})
	return
}
