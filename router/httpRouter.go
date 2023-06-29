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
	} else if G.C.Zap.Mode == "pro" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		panic("启动格式不正确，应为dev(开发模式)或pro(生产模式)")
	}
	// 不需要认证的组
	noAuthGroup := router.Group("")
	{
		// 代理静态文件
		noAuthGroup.StaticFile("/", path.Join("templates", "index.html"))
		noAuthGroup.Static("/static/js", "templates/static/js")
		noAuthGroup.Static("/static/css", "templates/static/css")
		noAuthGroup.Static("/static/media", "templates/static/media")
		noAuthGroup.StaticFile("manifest.json", path.Join("templates", "manifest.json"))
		noAuthGroup.StaticFile("logo192.png", path.Join("templates", "logo192.png"))
		noAuthGroup.POST(G.C.Server.Path+"/login", func(c *gin.Context) {
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
			var userService = service.NewUserService()
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
	authGroup := router.Group(G.C.Server.Path+"/rest", middleware.WhiteAuth(), middleware.AuthJwtToken())
	{
		authGroup.GET("test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"msg":  "请求成功",
			})
			return
		})
		// 用户相关
		userGroup := authGroup.Group("/user")
		{
			// 获取用户信息
			userGroup.POST("/getUserInfo", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"code":    http.StatusOK,
					"message": "获取用户数据成功",
					"success": true,
					"result":  map[string]any{},
				})
			})
			// 登出
			userGroup.POST("/logout", func(c *gin.Context) {
				// 登出逻辑
				// 1.登出后把该用户的token加入redis，并设置过期时间和token的有效时间一致
				// 2.请求接口时，校验token通过后，先去redis中读取该token，如果存在，则证明该用户已经登出，直接返回没有权限
				//（存入redis的token有效时间永远比生成的token过期时间长），需要重新登陆获取token
				token := c.GetHeader("x-token")
				username, exist := c.Get("username")
				if !exist {
					c.JSON(http.StatusOK, gin.H{
						"code":    http.StatusOK,
						"message": "登出失败",
						"success": false,
					})
					return
				}
				userService := service.NewUserService()
				err := userService.Logout(token, username.(string))
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"code":    http.StatusOK,
						"message": "登出失败:" + err.Error(),
						"success": false,
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"code":    http.StatusOK,
					"message": "登出成功",
					"success": true,
				})
			})
		}
		// 发布管理
		publishGroup := authGroup.Group("/publish") // 发布组路由
		{
			// 发布管理
			publishGroup.POST("/release", Release)
		}
		// 项目管理
		projectGroup := authGroup.Group("/project") // 项目组路由
		{
			// 项目管理
			projectGroup.POST("/list", ProjectList)
			// 获取环境列表
			projectGroup.POST("/getProjectEnvList", ProjectEnvList)
			// 获取项目类型
			projectGroup.POST("/getProjectTypeList", ProjectTypeList)
		}
		// 数据库管理
		databaseGroup := authGroup.Group("/database") // 数据库组路由
		{
			// 导出或备份整个数据库
			databaseGroup.POST("/total/export", func(c *gin.Context) {
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
				fmt.Println(totalDto)
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
				err = databaseService.HandleTotalMysqlDump(tempPath, totalDto.IgnoreTables...)
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
					// 备份路径=备份路径+数据库名称+年月日+备份的数据库文件
					backPath := G.C.Ops.Mysql.BackUpPath + "/" + totalDto.DbName + "/" + time.Now().Format("20060102")
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
					return
				case 2:
					// 为1则是导出数据库，先存到临时目录，再从临时目录读取返回文件流
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
			databaseGroup.POST("/single/export", func(c *gin.Context) {
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
			databaseGroup.POST("/dynamic/execSql", func(c *gin.Context) {
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
			databaseGroup.POST("/dbAndTable/list", func(c *gin.Context) {
				var databaseService = service.NewDatabaseService("")
				// 先检查系统是否安装了mysql
				err := databaseService.CheckMysql()
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"code":    200,
						"success": false,
						"message": "获取数据库列表失败，请联系系统管理员",
					})
					return
				}
				// 默认可操作除information_schema、mysql、performance_schema和sys这几个数据库之外的表数据
				ignoreDbs := G.C.Ops.Mysql.IgnoreDbs
				dataMap, err := databaseService.GetDbAndTableList(ignoreDbs)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"code":    200,
						"success": false,
						"message": "获取数据库列表失败",
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
			// 执行sql文件
			databaseGroup.POST("execSqlFile", func(c *gin.Context) {
				dbName, flag := c.GetPostForm("dbName")
				if !flag {
					c.JSON(http.StatusOK, gin.H{
						"code":    200,
						"success": false,
						"message": "参数错误",
					})
					return
				}
				file, err := c.FormFile("file")
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"code":    200,
						"success": false,
						"message": "请先上传sql文件",
					})
					return
				}
				// 判断是否是.sql文件结尾的文件
				hasSqlSuffix := strings.HasSuffix(file.Filename, ".sql")
				if !hasSqlSuffix {
					c.JSON(http.StatusOK, gin.H{
						"code":    200,
						"success": false,
						"message": "请上传sql类型文件",
					})
					return
				}
				tempFileName := uuid.NewV4().String()
				tempFilePath := "temp/" + tempFileName + "/" + file.Filename
				// 把文件保存到临时目录
				err = c.SaveUploadedFile(file, tempFilePath)
				if err != nil {
					G.Logger.Errorf("执行sql文件失败,失败原因[%s]", err.Error())
					c.JSON(http.StatusOK, gin.H{
						"code":    200,
						"success": false,
						"message": "执行sql文件失败，请联系管理员",
					})
					return
				}
				databaseService := service.NewDatabaseService(dbName)
				err = databaseService.ExecSqlFile(tempFilePath)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"code":    200,
						"success": false,
						"message": err,
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"success": true,
					"message": "执行sql文件成功",
				})
			})
		}
	}

	router.Run(fmt.Sprintf(":%d", G.C.Server.Port))
}
