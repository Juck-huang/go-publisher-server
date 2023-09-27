package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hy.juck.com/go-publisher-server/config"
	"hy.juck.com/go-publisher-server/middleware"
	database2 "hy.juck.com/go-publisher-server/router/database"
	"hy.juck.com/go-publisher-server/router/fileManager"
	"hy.juck.com/go-publisher-server/router/project"
	"hy.juck.com/go-publisher-server/router/publish"
	"hy.juck.com/go-publisher-server/router/user"
)

var (
	G = config.G
)

func NewHttpRouter() {
	router := gin.Default()
	router.Use(middleware.Cors())
	if G.C.Zap.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
	} else if G.C.Zap.Mode == "pro" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		panic("启动格式不正确，应为dev(开发模式)或pro(生产模式)")
	}

	//// 代理静态文件组
	//staticGroup := router.Group("aps")
	//{
	//	staticGroup.StaticFile("/", path.Join("templates", "index.html"))
	//	staticGroup.Static("/static/js/", "./templates/static/js")
	//	staticGroup.Static("/static/css/", "./templates/static/css")
	//	staticGroup.Static("/static/media/", "./templates/static/media")
	//	staticGroup.StaticFile("manifest.json", path.Join("templates", "manifest.json"))
	//	staticGroup.StaticFile("logo192.png", path.Join("templates", "logo192.png"))
	//}

	// 父组，带请求前缀，先校验白名单
	parentGroup := router.Group(G.C.Server.Path, middleware.WhiteAuth())
	{
		// 不需要认证的组，例如登录
		noAuthGroup := parentGroup.Group("")
		{
			// 登录
			noAuthGroup.POST("/login", user.Login)
		}
		// 需要认证的组，需要，然后校验token
		authGroup := parentGroup.Group("/rest", middleware.AuthJwtToken())
		{
			// 用户相关
			userGroup := authGroup.Group("/user")
			{
				// 获取用户信息
				userGroup.POST("/getUserInfo", user.GetUserInfo)
				// 登出
				userGroup.POST("/logout", user.Logout)
				// 自动上报ip服务接收方
				userGroup.POST("/updateLoginWhiteIp", user.UpdateLoginWhiteIp)
			}
			// 发布管理
			publishGroup := authGroup.Group("/publish") // 发布组路由
			{
				// 发布管理
				publishGroup.POST("/release", publish.Release)
			}
			// 项目管理
			projectGroup := authGroup.Group("/project") // 项目组路由
			{
				// 项目管理
				projectGroup.POST("/list", project.ProjectList)
				// 获取环境列表
				projectGroup.POST("/getProjectEnvList", project.ProjectEnvList)
				// 获取项目类型
				projectGroup.POST("/getProjectTypeList", project.ProjectTypeList)
			}
			// 数据库管理
			databaseGroup := authGroup.Group("/database") // 数据库组路由
			{
				// 导出或备份整个数据库
				databaseGroup.POST("/total/export", database2.ExportTotal)
				// 单独导出某几个表
				databaseGroup.POST("/single/export", database2.SingleExport)
				// 动态执行sql
				databaseGroup.POST("/dynamic/execSql", database2.DynamicSql)
				// 获取可操作的数据库和对应表列表
				databaseGroup.POST("/dbAndTable/list", database2.GetDbAndTableList)
				// 执行sql文件
				databaseGroup.POST("/execSqlFile", database2.ExecSqlFile)
			}
			// 文件管理
			fileManageGroup := authGroup.Group("/fileManager") // 文件管理组路由
			{
				// 获取项目信息列表
				fileManageGroup.POST("/getProjectList", fileManager.GetProjectList)
				// 获取项目文件列表
				fileManageGroup.POST("/getProjectFileList", fileManager.GetProjectFileList)
				// 上传项目文件，单独上传
				fileManageGroup.POST("/uploadProjectFile", fileManager.UploadProjectFile)
				// 下载项目文件,包括项目文件夹，单独下载
				fileManageGroup.POST("/downloadProjectFile", fileManager.DownloadProjectFile)
				// 上传项目切片文件
				fileManageGroup.POST("/uploadProjectFileChunk", fileManager.UploadProjectFileChunk)
				// 合并切片文件成一个
				fileManageGroup.POST("/mergeFileChunk", fileManager.MergeFileChunk)
				// 获取需要下载文件大小
				fileManageGroup.POST("/getFileSize", fileManager.GetFileSize)
				// 读取文件信息
				fileManageGroup.POST("/getFileContent", fileManager.GetFileContent)
				// 保存文件内容
				fileManageGroup.POST("/saveFileContent", fileManager.SaveFileContent)
				// 删除文件或文件夹
				fileManageGroup.POST("/removeFile", fileManager.RemoveFile)
				// 新建文件夹
				fileManageGroup.POST("/addFolder", fileManager.AddFolder)
				// 新建文件
				fileManageGroup.POST("/addFile", fileManager.AddFile)
				// 重命名文件夹或文件
				fileManageGroup.POST("/reNameFile", fileManager.ReNameFile)
				// 移动或复制文件或文件夹
				fileManageGroup.POST("/moveOrCopyFile", fileManager.MoveOrCopyFile)
				// 压缩文件夹或文件
				fileManageGroup.POST("/compressFileOrFolder", fileManager.CompressFileOrFolder)
				// 解压文件
				fileManageGroup.POST("/decompressionFile", fileManager.DecompressionFile)
			}
		}

		//wsAuthGroup := parentGroup.Group("/ws")
		//{
		//	// 查看实时日志
		//	wsAuthGroup.GET("/getRealTimeLog", fileManager.GetRealTimeLog)
		//}
	}

	router.Run(fmt.Sprintf(":%d", G.C.Server.Port))
}
