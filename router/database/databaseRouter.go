package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"hy.juck.com/go-publisher-server/config"
	"hy.juck.com/go-publisher-server/dto/database"
	"hy.juck.com/go-publisher-server/middleware"
	"hy.juck.com/go-publisher-server/service"
	"hy.juck.com/go-publisher-server/utils"
)

var (
	G = config.G
)

// ExportTotal 导出或备份所有表
func ExportTotal(c *gin.Context) {
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
		defer os.RemoveAll(tempPathPrefix)
		// 为1则是备份数据库，先存到临时目录，压缩，再复制到备份目录下
		// 备份路径=备份路径+数据库名称+年月日+备份的数据库文件
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
		err = databaseService.CopyFile(zipTempFilePath, backPath)
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
			"message": "备份数据库" + totalDto.DbName + "成功",
		})
		return
	case 2:
		// 为2则是导出数据库，先存到临时目录，然后根据key格式为,total_export_${uuid}:{save_path:临时路径, expire: ${时间戳}, completed: false}存redis
		// 定义一个get接口从redis中读取文件路径，completed为true返回文件流，内部全部导出完成不清除临时目录，使用定时任务，每半个小时执行清理expire超过1小时的对应的临时目录和redis数据
		sqlFileName := totalDto.DbName + "-" + nowTime + ".sql"
		zipFileName := totalDto.DbName + "-" + nowTime + ".zip"
		zipTempFilePath := tempPathPrefix + "/" + zipFileName
		redisKey := fmt.Sprintf("total_export_%s", uuidStr)
		G.RedisClient.HSet(redisKey, "save_path", zipTempFilePath)
		G.RedisClient.HSet(redisKey, "create_time", time.Now().UnixMilli())
		G.RedisClient.HSet(redisKey, "completed", false)
		G.RedisClient.HSet(redisKey, "name", zipFileName)
		go func() {
			err = databaseService.HandleTotalMysqlDump(tempPath, totalDto.IgnoreTables...)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"success": false,
					"message": err.Error(),
				})
				return
			}

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
			G.RedisClient.HSet(redisKey, "completed", true)
		}()
		// 返回redis key给前端去调用
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": true,
			"message": "生成导出任务成功，稍后下载导出后的文件",
			"result":  redisKey,
		})
		return
		// c.FileAttachment(zipTempFilePath, zipFileName)
	default:
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": "导出或备份数据库失败，失败原因：类型错误",
		})
		return
	}
}

// 获取导出进度详情
func GetExportDetail(c *gin.Context) {
	// 获取post参数
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "参数解析错误",
		})
		return
	}
	var requestMap = make(map[string]any, 1)
	json.Unmarshal(body, &requestMap)
	redisKey := requestMap["key"]
	if redisKey == nil || redisKey == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "参数解析错误",
		})
		return
	}
	// 根据key从redis中获取对应状态数据和下载路径
	redisDataMap, err := G.RedisClient.HGetAll(redisKey.(string)).Result()
	if err != nil {
		G.Logger.Errorf("获取导出进度详情失败:[%s]\n", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "获取导出进度详情失败",
		})
		return
	}
	create_time := redisDataMap["create_time"]
	if create_time == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "未获取到相关导出进度详情",
		})
		return
	}
	var redisDataDto database.RedisDataDto
	if redisDataMap["completed"] == "1" {
		redisDataDto.Completed = true
	}
	redisDataDto.CreateTime = redisDataMap["create_time"]
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"success": true,
		"message": "获取进度详情成功",
		"result":  redisDataDto,
	})
}

// 下载导出文件流
func DownloadExportFile(c *gin.Context) {
	var errInfo error
	key := c.Query("key")
	token := c.Query("token")
	// 校验参数是否存在
	if key == "" || token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"success": false,
			"msg":     "参数缺失",
		})
		return
	}
	// 校验token正确性
	claims, errInfo := utils.ParseToken(token)
	if errInfo != nil {
		G.Logger.Errorf("token解析错误或已经失效:[%s]", errInfo)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "token解析错误或token已失效",
			"success": false,
			"result":  []string{},
		})
		c.Abort()
		return
	}
	errInfo = middleware.CheckLogoutRedis(claims.Username, token)
	if errInfo != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": errInfo,
			"success": false,
			"result":  []string{},
		})
		c.Abort()
		return
	}

	// 根据key从redis中获取下载路径
	resultMap, errInfo := G.RedisClient.HGetAll(key).Result()
	if errInfo != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "获取到相关下载信息出错",
			"success": false,
			"result":  []string{},
		})
		return
	}
	if resultMap["create_time"] == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "未获取到相关下载信息",
			"success": false,
			"result":  []string{},
		})
		return
	}
	downloadPath := resultMap["save_path"]
	fileName := resultMap["name"]
	// 如果已经导出完成则直接读取文件，否则提示正在处理中
	completed := resultMap["completed"]
	if completed == "1" {
		G.Logger.Infof("读取下载文件路径:[%s],文件名:[%s]", downloadPath, fileName)
		// 读取文件后直接返回流
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename="+fileName)
		c.Set("Content-Transfer-Encoding", "binary")
		c.File(downloadPath)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "系统正在处理中...",
		"success": false,
		"result":  []string{},
	})
}

// SingleExport 单独导出
func SingleExport(c *gin.Context) {
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
}

// DynamicSql 动态执行sql(需加权限控制)
func DynamicSql(c *gin.Context) {
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
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": "执行动态sql失败：当前用户不存在",
		})
		return
	}
	resultList, err := databaseService.DynamicExecSql(dyncSql, username.(string))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"result":  resultList,
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
}

// GetDbAndTableList 获取可操作的数据库和表
func GetDbAndTableList(c *gin.Context) {
	var databaseService = service.NewDatabaseService("")
	// 先检查系统是否安装了mysql
	err := databaseService.CheckMysql()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": "获取数据库列表失败，请联系管理员",
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
			"message": "获取数据库列表失败，请联系管理员",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"success": true,
		"result":  dataMap,
		"message": "获取数据库和表信息成功",
	})
}

// ExecSqlFile 执行sql文件(需加权限控制)
func ExecSqlFile(c *gin.Context) {
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
	defer os.RemoveAll("temp/" + tempFileName)
	databaseService := service.NewDatabaseService(dbName)
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": "执行sql失败：获取当前用户信息失败",
		})
		return
	}
	err = databaseService.ExecSqlFile(tempFilePath, username.(string))
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
		"message": "执行sql脚本成功",
	})
}
