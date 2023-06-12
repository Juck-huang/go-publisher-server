package router

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

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
		// 导出或备份整个数据库
		authGroup.POST("/database/total/export", func(c *gin.Context) {
			var totalDto database.TotalDto
			err := c.ShouldBindJSON(&totalDto)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"success": false,
					"message": "参数解析错误",
				})
				return
			}
			if totalDto.DbName == "" || totalDto.Type == 0 {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"success": false,
					"message": "参数缺失",
				})
				return
			}
			// 需要先判断系统是否有安装mysqldump，有则继续，否则退出程序
			var databaseService = service.NewDatabaseService(totalDto.DbName)
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
			tempPath := tempPathPrefix + "/" + totalDto.DbName + "-" + nowTime + ".sql"
			err = os.MkdirAll(tempPathPrefix, os.ModePerm)
			defer os.RemoveAll(tempPathPrefix)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"success": false,
					"message": err.Error(),
				})
				return
			}
			switch totalDto.Type {
			case 1:
				// 为1则是备份数据库，先存到临时目录，再复制到备份目录下
				err = databaseService.HandleTotalMysqlDump(tempPath, totalDto.IgnoreTables...)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"code":    200,
						"success": false,
						"message": err.Error(),
					})
					return
				}
				// 备份路径=备份路径+年月日+备份的数据库文件
				backPath := G.C.DB.Mysql.BackUpPath + "/" + time.Now().Format("20060102")
				_, err = os.Stat(backPath)
				if os.IsNotExist(err) {
					if err = os.MkdirAll(backPath, os.ModePerm); err != nil {
						G.Logger.Error("创建备份文件夹失败:", err.Error())
						c.JSON(http.StatusOK, gin.H{
							"code":    200,
							"success": false,
							"message": "备份或导出失败",
						})
						return
					}
				}
				err = databaseService.CopyFile(tempPath, backPath)
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
				err = databaseService.HandleTotalMysqlDump(tempPath, totalDto.IgnoreTables...)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"code":    200,
						"success": false,
						"message": err.Error(),
					})
					return
				}
				sqlFileName := totalDto.DbName + "-" + nowTime + ".sql"
				zipFileName := totalDto.DbName + "-" + nowTime + ".zip"
				zipTempFilePath := tempPathPrefix + "/" + zipFileName
				// 压缩文件
				err = databaseService.ZipFile(tempPathPrefix, sqlFileName, zipFileName)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"code":    200,
						"success": false,
						"message": err.Error(),
					})
					return
				}
				c.FileAttachment(zipTempFilePath, zipFileName)
			default:
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"success": false,
					"message": "导出或备份数据库失败，失败原因：类型错误",
				})
				return
			}
		})
		// 单独导出某几个表
		authGroup.POST("/database/single/export", func(c *gin.Context) {
			var singleDto database.SingleDto
			err := c.BindJSON(&singleDto)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"success": false,
					"message": "参数解析错误",
				})
				return
			}
			if singleDto.DbName == "" || len(singleDto.ExportTables) == 0 {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"success": false,
					"message": "参数缺失",
				})
				return
			}
			uuidStr := uuid.NewV4().String()
			nowTime := time.Now().Format("20060102150405")
			tempPathPrefix := "temp/" + uuidStr
			tempPath := tempPathPrefix + "/" + singleDto.DbName + "-" + nowTime + ".sql"
			err = os.MkdirAll(tempPathPrefix, os.ModePerm)
			defer os.RemoveAll(tempPathPrefix)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"success": false,
					"message": err.Error(),
				})
				return
			}
			// 需要先判断系统是否有安装mysqldump，有则继续，否则退出程序
			var databaseService = service.NewDatabaseService(singleDto.DbName)
			err = databaseService.CheckMysqlDump()
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"success": false,
					"message": err.Error(),
				})
				return
			}
			// 导出表
			err = databaseService.SingleExportTables(tempPath, singleDto.ExportTables...)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"success": false,
					"message": err.Error(),
				})
				return
			}
			sqlFileName := singleDto.DbName + "-" + nowTime + ".sql"
			zipFileName := singleDto.DbName + "-" + nowTime + ".zip"
			zipTempFilePath := tempPathPrefix + "/" + zipFileName
			// 压缩文件
			err = databaseService.ZipFile(tempPathPrefix, sqlFileName, zipFileName)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"success": false,
					"message": err.Error(),
				})
				return
			}
			c.FileAttachment(zipTempFilePath, zipFileName)
		})
		// 动态执行sql
		authGroup.POST("/database/dynamic/execSql", func(c *gin.Context) {
			var dynamicExecDto database.DynamicExecDto
			err := c.BindJSON(&dynamicExecDto)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"success": false,
					"message": "参数解析错误",
				})
				return
			}
			if dynamicExecDto.DbName == "" || dynamicExecDto.Sql == "" {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"success": false,
					"message": "参数缺失",
				})
				return
			}
			// 需要先判断系统是否有安装mysqldump，有则继续，否则退出程序
			var databaseService = service.NewDatabaseService(dynamicExecDto.DbName)
			err = databaseService.CheckMysqlDump()
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"success": false,
					"message": err.Error(),
				})
				return
			}
			// 动态执行sql
			// 首先需要先移除sql多余的空格
			dyncSql := strings.TrimSpace(dynamicExecDto.Sql)
			resultList, err := databaseService.DynamicExecSql(dyncSql)
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
				"result":  resultList,
				"message": "执行动态sql成功",
			})
		})
		// 获取可操作的数据库和对应表列表
		authGroup.POST("/database/dbAndTable/list", func(c *gin.Context) {
			var databaseService = service.NewDatabaseService("")
			// 默认可操作除information_schema、mysql、performance_schema和sys这几个数据库之外的表数据
			ignoreDbs := G.C.DB.Mysql.IgnoreDbs
			dataMap, err := databaseService.GetDbAndTableList(ignoreDbs)
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
				"result":  dataMap,
				"message": "获取数据库和表信息成功",
			})
		})
	}

	router.Run(fmt.Sprintf(":%d", G.C.Server.Port))
}
