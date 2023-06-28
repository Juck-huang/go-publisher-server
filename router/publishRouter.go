package router

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"hy.juck.com/go-publisher-server/common"
	"hy.juck.com/go-publisher-server/service"
)

// Release 应用发布
func Release(c *gin.Context) {
	var err error // 定义一个全局错误
	envId, _ := c.GetPostForm("envId")
	projectId, _ := c.GetPostForm("projectId")
	typeId, _ := c.GetPostForm("typeId")
	if projectId == "" || typeId == "" || envId == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": "发布失败: 参数缺失",
		})
		return
	}
	// 根据项目类型id和项目id查询对应的数据
	projectReleaseService := service.NewProjectReleaseService(projectId)
	projectReleaseDto := projectReleaseService.GetProjectRelease(envId, typeId)
	if projectReleaseDto.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": "发布失败: 项目不存在",
		})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": "发布失败: 请上传项目项目文件!",
		})
		return
	}
	// // 判断文件类型是否是zip,目前只支持zip格式
	// fileIsRational := strings.HasSuffix(file.Filename, ".zip")
	// if !fileIsRational {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"code":    200,
	// 		"success": false,
	// 		"message": "发布失败: 请上传zip格式文件!",
	// 	})
	// 	return
	// }
	tempFileName := uuid.NewV4().String()
	executable, err := os.Getwd()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": "发布失败!",
		})
		return
	}
	tempFilePath := executable + "/temp/" + tempFileName // 临时文件绝对目录
	tempFile := tempFilePath + "/" + file.Filename       //临时文件全路径
	// 把文件保存到临时目录
	err = c.SaveUploadedFile(file, tempFile)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": "项目发布失败:" + err.Error(),
		})
		return
	}
	defer func() {
		_, err = common.ExecCommand(true, "-c", "rm -rf "+tempFilePath+" && echo '移除成功'")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"success": false,
				"message": "项目发布失败,失败原因2:" + err.Error(),
			})
			return
		}
	}()

	// 执行的构建脚本路径需要在数据库配置好
	split := strings.Split(projectReleaseDto.Params, ",")
	command := []string{
		projectReleaseDto.BuildScriptPath,
		tempFilePath,
	}
	command = append(command, split...)
	msg, err := common.ExecCommand(true, command...)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": "项目发布失败,失败原因1:" + err.Error(),
		})
		return
	}
	if msg != "项目发布成功" {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": "项目发布失败,具体原因3:" + msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"success": true,
		"message": "项目发布成功",
	})
}
