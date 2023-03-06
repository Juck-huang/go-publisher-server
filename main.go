package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	// 代理静态文件
	router.StaticFile("/", path.Join("templates", "index.html"))
	router.Static("/js", "templates/js")
	router.Static("/css", "templates/css")
	router.POST("/uploadFile", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "请求失败:" + err.Error(),
			})
			return
		}
		c.SaveUploadedFile(file, "temp/"+file.Filename)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "发布成功",
		})
	})
	router.Run(":8002")
}
