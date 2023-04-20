package main

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"hy.juck.com/go-publisher-server/common"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	// 代理静态文件
	router.StaticFile("/", path.Join("templates", "index.html"))
	router.Static("/static/js", "templates/static/js")
	router.Static("/static/css", "templates/static/css")
	router.StaticFile("manifest.json", path.Join("templates", "manifest.json"))
	router.StaticFile("logo192.png", path.Join("templates", "logo192.png"))
	router.POST("/publishProject", func(c *gin.Context) {
		var err error // 定义一个全局错误
		projectId, _ := c.GetPostForm("projectId")
		publishType, _ := c.GetPostForm("type")
		var projectIdList = []string{"1", "2"} // 项目id列表
		var typeList = []string{"app", "pc"}   // 类型列表
		// 判断是否是存在的项目
		isExistProject := common.EleIsExistSlice(projectId, projectIdList)
		if !isExistProject {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"success": false,
				"message": "发布失败: 项目不存在,请检查!",
			})
			return
		}
		// 判断存在类型
		isExistType := common.EleIsExistSlice(publishType, typeList)
		if !isExistType {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"success": false,
				"message": "发布失败: 发布类型不存在,请检查!",
			})
			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"success": false,
				"message": "发布失败: 请上传项目压缩文件!",
			})
			return
		}
		// 判断文件类型是否是zip,目前只支持zip格式
		fileIsRational := strings.HasSuffix(file.Filename, ".zip")
		if !fileIsRational {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"success": false,
				"message": "发布失败: 请上传zip格式文件!",
			})
			return
		}
		tempFileName := uuid.NewV4().String()
		// 把文件保存到临时目录
		err = c.SaveUploadedFile(file, "temp/"+tempFileName+"/"+file.Filename)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"success": false,
				"message": "项目发布失败:" + err.Error(),
			})
			return
		}

		msg, err := common.ExecCommand(true, fmt.Sprintf("scripts/%s/%s/build", projectId, publishType), tempFileName, file.Filename)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"success": false,
				"message": "项目发布失败,失败原因1:" + err.Error(),
			})
			return
		}
		if msg != "打包成功" {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"success": false,
				"message": "项目发布失败,具体原因3:" + msg,
			})
			return
		}
		defer func() {
			_, err = common.ExecCommand(true, "-c", "rm -rf temp/"+tempFileName+" && echo '移除成功'")
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"success": false,
					"message": "项目发布失败,失败原因2:" + err.Error(),
				})
				return
			}
		}()

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": true,
			"message": "项目发布成功",
		})
	})
	router.Run(":8002")
}
