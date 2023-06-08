package router

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"hy.juck.com/go-publisher-server/common"
	"hy.juck.com/go-publisher-server/config"
	"hy.juck.com/go-publisher-server/dto/database"
	"hy.juck.com/go-publisher-server/dto/login"
	"hy.juck.com/go-publisher-server/middleware"
	"hy.juck.com/go-publisher-server/service"
	"hy.juck.com/go-publisher-server/utils"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

var (
	G = config.G
)

func NewHttpRouter() {
	router := gin.Default()
	router.Use(cors.Default())
	if G.C.Zap.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	// 不需要认证的组
	noAuthGroup := router.Group("")
	{
		// 代理静态文件
		noAuthGroup.StaticFile("/", path.Join("templates", "index.html"))
		noAuthGroup.Static("/static/js", "templates/static/js")
		noAuthGroup.Static("/static/css", "templates/static/css")
		noAuthGroup.StaticFile("manifest.json", path.Join("templates", "manifest.json"))
		noAuthGroup.StaticFile("logo192.png", path.Join("templates", "logo192.png"))
		noAuthGroup.POST("/login", func(c *gin.Context) {
			var requestDto login.RequestDto
			err := c.ShouldBindJSON(&requestDto)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"success": false,
					"result":  map[string]any{},
					"message": err.Error(),
				})
				return
			}
			var userService = service.UserService{}
			err = userService.CheckUsernameAndPassword(requestDto.Username, requestDto.Password)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"success": false,
					"result":  map[string]any{},
					"message": err.Error(),
				})
				return
			}
			token, err := utils.GenToken(requestDto.Username)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"success": false,
					"result":  map[string]any{},
					"message": err.Error(),
				})
				return
			}
			// 登录接口
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"success": true,
				"result": map[string]any{
					"token": token,
				},
				"message": "登录成功",
			})
		})
	}

	// 需要认证的组
	authGroup := router.Group("/rest", middleware.AuthJwtToken())
	{
		// 项目发布
		authGroup.POST("/project/publish", func(c *gin.Context) {
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
		// 导出或备份数据库
		authGroup.POST("/database/export", func(c *gin.Context) {
			var requestDto database.RequestDto
			err := c.ShouldBindJSON(&requestDto)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"success": false,
					"message": "参数解析错误",
				})
				return
			}
			if requestDto.DbName == "" || requestDto.Status == 0 {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"success": false,
					"message": "参数缺失",
				})
				return
			}
			// 需要先判断系统是否有安装mysqldump，有则继续，否则退出程序
			var databaseService = service.NewDatabaseService(requestDto.DbName)
			err = databaseService.CheckMysqlDump()
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"success": false,
					"message": err.Error(),
				})
				return
			}
			uuidStr := uuid.NewV4().String()
			nowTime := time.Now().Format("20060102150405")
			tempPathPrefix := "temp/" + uuidStr
			tempPath := tempPathPrefix + "/" + requestDto.DbName + "-" + nowTime + ".sql"
			err = os.MkdirAll(tempPathPrefix, 0644)
			defer os.RemoveAll(tempPathPrefix)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"success": false,
					"message": err.Error(),
				})
				return
			}
			switch requestDto.Status {
			case 1:
				// 为1则是备份数据库，先存到临时目录，再复制到备份目录下
				err = databaseService.HandleMysqlDump(tempPath, requestDto.IgnoreTables...)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"code":    200,
						"success": false,
						"message": err.Error(),
					})
					return
				}
				err = databaseService.CopyFile(tempPath, G.C.DB.Mysql.BackUpPath)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"code":    200,
						"success": false,
						"message": err.Error(),
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"success": true,
					"message": "备份数据库成功",
				})
			case 2:
				// 为1则是导出数据库，先存到临时目录，再从临时目录读取返回文件流
				err = databaseService.HandleMysqlDump(tempPath, requestDto.IgnoreTables...)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"code":    200,
						"success": false,
						"message": err.Error(),
					})
					return
				}
				c.FileAttachment(tempPath, requestDto.DbName+"-"+nowTime+".sql")
			default:
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"success": false,
					"message": errors.New("导出或备份数据库失败，失败原因：类型错误"),
				})
			}
		})
	}

	router.Run(fmt.Sprintf(":%d", G.C.Server.Port))
}
