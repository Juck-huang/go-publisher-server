package dbBak

import (
	"github.com/gin-gonic/gin"
	"hy.juck.com/go-publisher-server/service"
	"net/http"
)

func BakReceive(c *gin.Context) {
	// 接收客户端发送的数据备份，form-data请求
	// 需要接收的form参数，第一个为接收类型type，1.数据库 2.配置文件 3.应用程序文件
	// 第二个参数为发送的文件，为了避免数据传输完整性，必须为zip文件
	// 第三个参数为生产项目id
	formFile, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "发送备份数据文件失败",
			"code":    http.StatusOK,
		})
		return
	}
	projectId, _ := c.GetPostForm("projectId")
	if projectId == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "参数缺失",
			"code":    http.StatusOK,
		})
		return
	}
	projectService := service.NewProjectService()
	username, _ := c.Get("username")
	projectList := projectService.GetProjectList(username.(string)) // 根据当前用户获取项目列表
	// 判断项目id是否在当前用户所拥有的项目列表中
	var existProjectId bool
	for _, project := range projectList {
		if string(project.Id) == projectId {
			existProjectId = true
			break
		}
	}
	if !existProjectId {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "项目不存在",
			"code":    http.StatusOK,
		})
		return
	}
	receiveType, _ := c.GetPostForm("type")
	switch receiveType {
	case "1":
		// 接收数据库zip文件(20230824暂时只接收数据库zip文件)
		c.SaveUploadedFile(formFile, "")
	case "2":
		// 接收配置zip文件
	case "3":
		// 接收附件zip文件
	default:
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "数据传输类型不识别，请重试",
			"code":    http.StatusOK,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
