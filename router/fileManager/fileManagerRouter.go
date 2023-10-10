package fileManager

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"hy.juck.com/go-publisher-server/config"
	"hy.juck.com/go-publisher-server/dto/fileManager"
	"hy.juck.com/go-publisher-server/middleware"
	"hy.juck.com/go-publisher-server/service"
	"hy.juck.com/go-publisher-server/utils"
)

var (
	G = config.G
)

// GetFileContent 获取文件内容
func GetFileContent(c *gin.Context) {
	var requestFileDto fileManager.RequestDto
	err := c.ShouldBindJSON(&requestFileDto)
	if err != nil {
		G.Logger.Errorf("参数解析错误，具体原因:[%s]", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "参数解析错误或参数缺失",
		})
		return
	}
	if requestFileDto.PathName == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "项目文件路径不能为空",
		})
		return
	}
	fileManagerService := service.NewFileManagerService()
	err = fileManagerService.SetProjectPath(requestFileDto)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("项目路径%s不存在，%s", requestFileDto.PathName, err.Error()),
		})
		return
	}
	_, err = fileManagerService.CheckProjectIsFile(requestFileDto.PathName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("路径[%s]解析错误,错误原因：%s", requestFileDto.PathName, err.Error()),
		})
		return
	}
	readPath := fileManagerService.CurrPath + "/" + requestFileDto.PathName
	file, err := os.ReadFile(readPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("路径[%s]解析错误,错误原因：%s", requestFileDto.PathName, err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"success": true,
		"result":  string(file),
		"message": "获取项目文件信息成功",
	})
}

// SaveFileContent 保存文件内容
func SaveFileContent(c *gin.Context) {
	var requestFileDto fileManager.RequestDto
	err := c.ShouldBindJSON(&requestFileDto)
	if err != nil {
		G.Logger.Errorf("参数解析错误，具体原因:[%s]", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "参数解析错误或参数缺失",
		})
		return
	}
	if requestFileDto.PathName == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "项目文件路径不能为空",
		})
		return
	}
	if requestFileDto.FileContent == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "项目文件内容不能为空",
		})
		return
	}
	fileManagerService := service.NewFileManagerService()
	err = fileManagerService.SetProjectPath(requestFileDto)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("项目路径%s不存在，%s", requestFileDto.PathName, err.Error()),
		})
		return
	}
	_, err = fileManagerService.CheckProjectIsFile(requestFileDto.PathName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("路径[%s]解析错误,错误原因：%s", requestFileDto.PathName, err.Error()),
		})
		return
	}
	editStatus, err := fileManagerService.SaveFileContent(requestFileDto.PathName, requestFileDto.FileContent)
	if !editStatus {
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
		"message": fmt.Sprintf("编辑文件%s成功", requestFileDto.PathName),
	})
}

// GetProjectList 获取项目信息列表
func GetProjectList(c *gin.Context) {
	// 1.查询项目列表
	projectService := service.NewProjectService()
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"msg":     "获取项目信息列表失败",
			"result":  map[string]any{},
		})
		return
	}
	projectList := projectService.GetProjectList(username.(string))
	// 2.根据项目id列表查询出所有的项目类型列表数据
	projectTypeService := service.NewProjectTypeService()
	var projectIdList []uint
	for _, project := range projectList {
		projectIdList = append(projectIdList, project.Id)
	}
	projectTypeList := projectTypeService.GetProjectTypeListByPIds(projectIdList)
	// 3.查询当前用户所有项目环境列表
	projectEnvService := service.NewProjectEnvService()
	projectEnvList := projectEnvService.GetProjectEnvListByPUsername(username.(string))
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"success": true,
		"msg":     "获取项目信息列表成功",
		"result": map[string]any{
			"projectEnvList":  projectEnvList,
			"projectList":     projectList,
			"projectTypeList": projectTypeList,
		},
	})
}

// GetProjectFileList 获取项目文件列表
func GetProjectFileList(c *gin.Context) {
	var projectFileDto fileManager.RequestDto
	err := c.ShouldBindJSON(&projectFileDto)
	if err != nil {
		G.Logger.Errorf("参数解析错误，具体原因:[%s]", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"msg":     "参数解析错误或参数缺失",
		})
		return
	}
	// 获取项目文件列表，包括文件夹和文件
	//path := "/Users/mac/Downloads/apps/stec-emerge-web/default"
	fileManagerService := service.NewFileManagerService()
	err = fileManagerService.SetProjectPath(projectFileDto)
	if err != nil {
		G.Logger.Errorf("获取项目文件失败，具体原因:[%s]", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"msg":     err.Error(),
			"result":  []string{},
		})
		return
	}
	fileManagerDtos, err := fileManagerService.GetFileList(projectFileDto.PathName)
	if err != nil {
		G.Logger.Errorf("获取项目文件失败，具体原因:[%s]", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"msg":     "项目目录不存在或不是目录",
			"result":  []string{},
		})
		return
	}
	if len(fileManagerDtos) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
			"msg":     "获取项目文件列表成功",
			"result":  []string{},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"success": true,
		"msg":     "获取项目文件列表成功",
		"result":  fileManagerDtos,
	})
}

// UploadProjectFile 上传项目文件
func UploadProjectFile(c *gin.Context) {
	projectId, _ := c.GetPostForm("projectId")
	projectEnvId, _ := c.GetPostForm("projectEnvId")
	projectTypeId, _ := c.GetPostForm("projectTypeId")
	pathName, _ := c.GetPostForm("pathName")
	if projectId == "" || projectEnvId == "" || projectTypeId == "" || pathName == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"msg":     "参数缺失",
		})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": "请上传项目文件!",
		})
		return
	}
	// 保存文件到项目目录下
	fileManagerService := service.NewFileManagerService()
	uintProjectId, err := strconv.ParseUint(projectId, 10, 64)
	uintProjectEnvId, err := strconv.ParseUint(projectEnvId, 10, 64)
	uintProjectTypeId, err := strconv.ParseUint(projectTypeId, 10, 64)
	var projectFileDto = fileManager.RequestDto{
		ProjectId:     uint(uintProjectId),
		ProjectEnvId:  uint(uintProjectEnvId),
		ProjectTypeId: uint(uintProjectTypeId),
	}
	err = fileManagerService.SetProjectPath(projectFileDto)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("上传项目文件%s失败，%s", file.Filename, err.Error()),
		})
		return
	}
	fileExist := fileManagerService.CheckProjectIsDir(pathName)
	if !fileExist {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("上传项目文件%s失败，文件或文件夹不存在", file.Filename),
		})
		return
	}
	err = c.SaveUploadedFile(file, fileManagerService.CurrPath+"/"+pathName+"/"+file.Filename)
	if err != nil {
		G.Logger.Errorf("上传项目文件失败：[%s]", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("上传项目文件%s失败", file.Filename),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"success": true,
		"message": fmt.Sprintf("上传项目文件%s成功", file.Filename),
	})
}

// UploadProjectFileChunk 上传项目chunk文件
func UploadProjectFileChunk(c *gin.Context) {
	fileName, isFileName := c.GetPostForm("fileName")
	chunkName, isChunkName := c.GetPostForm("chunkName")
	if !isFileName {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "文件名不能为空",
		})
		return
	}
	if !isChunkName {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "文件chunk名不能为空",
		})
		return
	}
	formFile, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "上传文件不能为空",
		})
		return
	}
	tempPath := fmt.Sprintf("temp/%s-chunk", fileName)
	_, err = os.Stat(tempPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(tempPath, os.ModePerm)
		if err != nil {
			fmt.Println("err", err)
			return
		}
	}
	err = c.SaveUploadedFile(formFile, tempPath+"/"+chunkName)
	if err != nil {
		fmt.Println("err", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "上传文件成功",
	})
}

// MergeFileChunk 合并上传的chunk文件为一个
func MergeFileChunk(c *gin.Context) {
	var mapData = make(map[string]any)
	err := c.ShouldBindJSON(&mapData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "参数解析错误:" + err.Error(),
		})
		return
	}
	fileName := mapData["fileName"].(string)
	if fileName == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "文件名不能为空",
		})
		return
	}
	fileHash := mapData["fileHash"].(string)
	if fileHash == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "文件hash不能为空",
		})
		return
	}
	chunkPath := fmt.Sprintf("temp/%s-chunk", fileName) // 文件夹chunk名
	targetPath := "temp/" + fileName
	fileManagerService := service.NewFileManagerService()
	err = fileManagerService.HandleMergeFile(chunkPath, targetPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": fmt.Sprintf("合并文件%s失败", fileName),
		})
		return
	}
	isHash, err := fileManagerService.CheckFileHash(targetPath, fileHash)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": fmt.Sprintf("合并文件%s失败", fileName),
		})
		return
	}
	if !isHash {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": fmt.Sprintf("合并文件%s失败", fileName),
		})
		return
	}
	G.Logger.Infof("合并文件%s完成", fileName)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("合并文件%s完成", fileName),
	})
}

// DownloadProjectFile 下载项目文件或文件夹
func DownloadProjectFile(c *gin.Context) {
	var requestFileDto fileManager.RequestDto
	err := c.ShouldBindJSON(&requestFileDto)
	if err != nil {
		G.Logger.Errorf("参数解析错误，具体原因:[%s]", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"msg":     "参数解析错误或参数缺失",
		})
		return
	}
	if requestFileDto.PathName == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"msg":     "项目文件路径不能为空",
		})
		return
	}
	fileManagerService := service.NewFileManagerService()
	err = fileManagerService.SetProjectPath(requestFileDto)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("项目路径%s不存在，%s", requestFileDto.PathName, err.Error()),
		})
		return
	}
	_, err = fileManagerService.CheckProjectIsFile(requestFileDto.PathName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("路径[%s]解析错误,错误原因：%s", requestFileDto.PathName, err.Error()),
		})
		return
	}
	downloadPath := fileManagerService.CurrPath + "/" + requestFileDto.PathName
	filePaths := strings.Split(requestFileDto.PathName, "/")
	fileName := filePaths[len(filePaths)-1]
	G.Logger.Infof("读取下载文件路径:[%s],文件名:[%s]", downloadPath, fileName)
	// 读取文件后直接返回流
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Set("Content-Transfer-Encoding", "binary")
	c.File(downloadPath)
}

// GetProjectFile 下载项目文件或文件夹
func GetProjectFile(c *gin.Context) {
	var requestFileDto fileManager.RequestDto
	var errInfo error
	requestFileDto.PathName = c.Query("pathName")
	token := c.Query("token")
	projectId, err := strconv.ParseUint(c.Query("projectId"), 0, 0)
	if err != nil {
		errInfo = err
	}
	requestFileDto.ProjectId = uint(projectId)
	projectEnvId, err := strconv.ParseUint(c.Query("projectEnvId"), 0, 0)
	if err != nil {
		errInfo = err
	}
	requestFileDto.ProjectEnvId = uint(projectEnvId)
	projectTypeId, err := strconv.ParseUint(c.Query("projectTypeId"), 0, 0)
	if err != nil {
		errInfo = err
	}
	requestFileDto.ProjectTypeId = uint(projectTypeId)
	if errInfo != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"msg":     "参数解析错误",
		})
		return
	}
	if requestFileDto.PathName == "" || projectId == 0 || projectEnvId == 0 || projectTypeId == 0 || token == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"msg":     "参数缺失",
		})
		return
	}
	// token解析错误不通过
	claims, errInfo := utils.ParseToken(token)
	if errInfo != nil {
		G.Logger.Errorf("token解析错误或已经失效:[%s]", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "token解析错误或token已失效",
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
			"result":  []string{},
		})
		c.Abort()
		return
	}
	fileManagerService := service.NewFileManagerService()
	errInfo = fileManagerService.SetProjectPath(requestFileDto)
	if errInfo != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("项目路径%s不存在，%s", requestFileDto.PathName, err.Error()),
		})
		return
	}
	_, errInfo = fileManagerService.CheckProjectIsFile(requestFileDto.PathName)
	if errInfo != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("路径[%s]解析错误,错误原因：%s", requestFileDto.PathName, err.Error()),
		})
		return
	}
	downloadPath := fileManagerService.CurrPath + "/" + requestFileDto.PathName
	filePaths := strings.Split(requestFileDto.PathName, "/")
	fileName := filePaths[len(filePaths)-1]
	G.Logger.Infof("读取下载文件路径:[%s],文件名:[%s]", downloadPath, fileName)
	// 读取文件后直接返回流
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Set("Content-Transfer-Encoding", "binary")
	c.File(downloadPath)
}

// CheckDownloadFile 校验需要下载文件信息
func CheckDownloadFile(c *gin.Context) {
	var requestFileDto fileManager.RequestDto
	err := c.ShouldBindJSON(&requestFileDto)
	if err != nil {
		G.Logger.Errorf("参数解析错误，具体原因:[%s]", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"msg":     "参数解析错误或参数缺失",
		})
		return
	}
	if requestFileDto.PathName == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"msg":     "项目文件路径不能为空",
		})
		return
	}
	fileManagerService := service.NewFileManagerService()
	err = fileManagerService.SetProjectPath(requestFileDto)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("项目路径%s不存在，%s", requestFileDto.PathName, err.Error()),
		})
		return
	}
	_, err = fileManagerService.CheckProjectIsFile(requestFileDto.PathName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("路径[%s]解析错误,错误原因：%s", requestFileDto.PathName, err.Error()),
		})
		return
	}
	filePath := fileManagerService.CurrPath + "/" + requestFileDto.PathName
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": "解析项目路径失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "校验项目文件成功",
		"result": map[string]any{
			"pathName": requestFileDto.PathName,
			"size":     fileInfo.Size(),
		},
	})
}

// RemoveFile 删除文件或文件夹
func RemoveFile(c *gin.Context) {
	var requestFileDto fileManager.RequestDto
	err := c.ShouldBindJSON(&requestFileDto)
	if err != nil {
		G.Logger.Errorf("参数解析错误，具体原因:[%s]", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "参数解析错误或参数缺失",
		})
		return
	}
	if requestFileDto.PathName == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "项目文件路径不能为空",
		})
		return
	}
	if requestFileDto.PathName == "/" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "项目根路径无法删除",
		})
		return
	}
	fileManagerService := service.NewFileManagerService()
	err = fileManagerService.SetProjectPath(requestFileDto)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("项目路径%s不存在，%s", requestFileDto.PathName, err.Error()),
		})
		return
	}
	fileCompletePath := fileManagerService.CurrPath + "/" + requestFileDto.PathName
	fileInfo, err := os.Stat(fileCompletePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("项目路径或文件夹%s不存在，%s", requestFileDto.PathName, err.Error()),
		})
		return
	}
	if fileInfo.IsDir() {
		// 说明是目录
		isDir := fileManagerService.CheckProjectIsDir(requestFileDto.PathName)
		if !isDir {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"success": false,
				"message": fmt.Sprintf("项目路径%s不存在", requestFileDto.PathName),
			})
			return
		}
	} else {
		_, err = fileManagerService.CheckProjectIsFile(requestFileDto.PathName)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"success": false,
				"message": fmt.Sprintf("项目文件路径%s不存在", requestFileDto.PathName),
			})
			return
		}
	}

	err = fileManagerService.RemoveFile(requestFileDto.PathName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": "删除文件或文件夹失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"success": true,
		"message": "删除项目文件或文件夹成功",
	})
}

// AddFolder 新建文件夹
func AddFolder(c *gin.Context) {
	var requestFileDto fileManager.RequestDto
	err := c.ShouldBindJSON(&requestFileDto)
	if err != nil {
		G.Logger.Errorf("参数解析错误，具体原因:[%s]", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "参数解析错误或参数缺失",
		})
		return
	}
	// 正则校验新建文件夹名称，必须是中文，大小写字母、数字和汉字组成,[\p{L}\p{N}]表示匹配一个Unicode字母或数字
	pattern := "^[a-zA-Z0-9\u4e00-\u9fa5.\\-_]+$"
	matchFolder, err := regexp.MatchString(pattern, requestFileDto.AddFolderName)
	if !matchFolder {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "新建文件夹名称格式不正确",
		})
		return
	}
	fileManagerService := service.NewFileManagerService()
	err = fileManagerService.SetProjectPath(requestFileDto)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("项目路径%s不存在，%s", requestFileDto.PathName, err.Error()),
		})
		return
	}
	// 校验项目文件夹路径是否存在
	isDir := fileManagerService.CheckProjectIsDir(requestFileDto.PathName)
	if !isDir {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("项目路径%s不存在", requestFileDto.PathName),
		})
		return
	}
	// 校验新增的文件夹路径是否存在，如果存在则提示文件夹已存在，否则运行通过
	dirIs := fileManagerService.CheckProjectIsDir(requestFileDto.PathName + "/" + requestFileDto.AddFolderName)
	if dirIs {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("文件夹%s已存在", requestFileDto.AddFolderName),
		})
		return
	}
	err = fileManagerService.AddFolder(requestFileDto.PathName, requestFileDto.AddFolderName)
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
		"message": "文件夹创建成功",
	})
}

// AddFile 新建文件
func AddFile(c *gin.Context) {
	var requestFileDto fileManager.RequestDto
	err := c.ShouldBindJSON(&requestFileDto)
	if err != nil {
		G.Logger.Errorf("参数解析错误，具体原因:[%s]", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "参数解析错误或参数缺失",
		})
		return
	}
	// 正则校验新建文件名称，必须是中文，大小写字母、数字和点汉字组成,[\p{L}\p{N}]表示匹配一个Unicode字母或数字
	pattern := "^[a-zA-Z0-9\u4e00-\u9fa5.\\-._]+$"
	matchFolder, err := regexp.MatchString(pattern, requestFileDto.AddFileName)
	if !matchFolder {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "新建文件名称格式不正确",
		})
		return
	}
	fileManagerService := service.NewFileManagerService()
	err = fileManagerService.SetProjectPath(requestFileDto)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("项目路径%s不存在，%s", requestFileDto.PathName, err.Error()),
		})
		return
	}
	// 校验项目文件路径是否存在
	isDir := fileManagerService.CheckProjectIsDir(requestFileDto.PathName)
	if !isDir {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("项目路径%s不存在", requestFileDto.PathName),
		})
		return
	}
	// 校验新增的文件路径是否存在，如果存在则提示文件已存在，否则运行通过
	dirIs, err := fileManagerService.CheckProjectIsFile(requestFileDto.PathName + "/" + requestFileDto.AddFileName)
	if dirIs {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("文件%s已存在", requestFileDto.AddFileName),
		})
		return
	}
	err = fileManagerService.AddFile(requestFileDto.PathName, requestFileDto.AddFileName)
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
		"message": "文件创建成功",
	})
}

// ReNameFile 重命名文件或文件夹
func ReNameFile(c *gin.Context) {
	var requestFileDto fileManager.ReNameRequestDto
	err := c.ShouldBindJSON(&requestFileDto)
	if err != nil {
		G.Logger.Errorf("参数解析错误，具体原因:[%s]", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "参数解析错误或参数缺失",
		})
		return
	}
	if requestFileDto.PathName == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "项目文件路径不能为空",
		})
		return
	}
	// 正则校验重命名文件夹名称，必须是中文，大小写字母、数字和汉字组成,[\p{L}\p{N}]表示匹配一个Unicode字母或数字
	pattern := "^[a-zA-Z0-9\u4e00-\u9fa5.\\-_]+$"
	matchFolder, err := regexp.MatchString(pattern, requestFileDto.NewFileName)
	if !matchFolder {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "重命名文件夹或文件名称格式不正确",
		})
		return
	}
	fileManagerService := service.NewFileManagerService()
	requestFile := fileManager.RequestDto{
		ProjectId:     requestFileDto.ProjectId,
		ProjectEnvId:  requestFileDto.ProjectEnvId,
		ProjectTypeId: requestFileDto.ProjectTypeId,
		PathName:      requestFileDto.PathName,
	}
	err = fileManagerService.SetProjectPath(requestFile)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("项目路径%s不存在，%s", requestFileDto.PathName, err.Error()),
		})
		return
	}
	fileCompletePath := fileManagerService.CurrPath + "/" + requestFileDto.PathName
	fmt.Println("fileCompletePath", fileCompletePath)
	fileInfo, err := os.Stat(fileCompletePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("项目路径%s不存在，%s", requestFileDto.PathName, err.Error()),
		})
		return
	}
	splitPathNames := strings.Split(requestFileDto.PathName, "/") // 取/分割后前面的所有
	newFileName := strings.Join(splitPathNames[0:len(splitPathNames)-1], "/") + "/" + requestFileDto.NewFileName
	// 说明是目录
	if fileInfo.IsDir() {
		// 校验新文件夹路径
		isDir := fileManagerService.CheckProjectIsDir(newFileName)
		if isDir {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"success": false,
				"message": fmt.Sprintf("重命名路径%s已存在", requestFileDto.NewFileName),
			})
			return
		}
	} else {
		// 说明是文件,则有扩展名
		exitFile, _ := fileManagerService.CheckProjectIsFile(newFileName)
		if exitFile {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"success": false,
				"message": fmt.Sprintf("项目文件路径%s已存在", requestFileDto.NewFileName),
			})
			return
		}
	}
	err = fileManagerService.ReNameFile(requestFileDto.PathName, newFileName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": "重命名文件或文件夹失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"success": true,
		"message": "重命名文件或文件夹成功",
	})
}

// MoveOrCopyFile 移动或复制文件或文件夹
func MoveOrCopyFile(c *gin.Context) {
	var requestFileDto fileManager.MoveOrCopyRequestDto
	err := c.ShouldBindJSON(&requestFileDto)
	if err != nil {
		G.Logger.Errorf("参数解析错误，具体原因:[%s]", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "参数解析错误或参数缺失",
		})
		return
	}
	if requestFileDto.PathName == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "项目文件路径不能为空",
		})
		return
	}

	fileManagerService := service.NewFileManagerService()
	requestFile := fileManager.RequestDto{
		ProjectId:     requestFileDto.ProjectId,
		ProjectEnvId:  requestFileDto.ProjectEnvId,
		ProjectTypeId: requestFileDto.ProjectTypeId,
		PathName:      requestFileDto.PathName,
	}
	err = fileManagerService.SetProjectPath(requestFile)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("项目路径%s不存在，%s", requestFileDto.PathName, err.Error()),
		})
		return
	}
	fileCompletePath := fileManagerService.CurrPath + "/" + requestFileDto.PathName
	extName := path.Ext(fileCompletePath)
	var toFileName string
	splitPathNames := strings.Split(requestFileDto.PathName, "/") // 取/分割后前面的所有
	toFileName = requestFileDto.ToPath + "/" + splitPathNames[len(splitPathNames)-1]
	if extName == "" {
		// 说明是目录
		isDir := fileManagerService.CheckProjectIsDir(requestFileDto.PathName)
		if !isDir {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"success": false,
				"message": fmt.Sprintf("项目路径%s不存在", requestFileDto.PathName),
			})
			return
		}
		isDir = fileManagerService.CheckProjectIsDir(requestFileDto.ToPath)
		if !isDir {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"success": false,
				"message": fmt.Sprintf("项目目标路径%s不存在", requestFileDto.ToPath),
			})
			return
		}
	} else {
		// 说明是文件
		_, err = fileManagerService.CheckProjectIsFile(requestFileDto.PathName)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"success": false,
				"message": fmt.Sprintf("项目文件路径%s不存在", requestFileDto.PathName),
			})
			return
		}
	}
	err = fileManagerService.MoveOrCopyFile(requestFileDto.Type, extName, requestFileDto.PathName, toFileName)
	var msg string
	if requestFileDto.Type == 1 {
		msg = "移动"
	} else if requestFileDto.Type == 2 {
		msg = "复制"
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("%s文件或文件夹失败", msg),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"success": true,
		"message": fmt.Sprintf("%s文件或文件夹成功", msg),
	})
}

// CompressFileOrFolder 压缩文件夹或目录
func CompressFileOrFolder(c *gin.Context) {
	var requestFileDto fileManager.RequestDto
	err := c.ShouldBindJSON(&requestFileDto)
	if err != nil {
		G.Logger.Errorf("参数解析错误，具体原因:[%s]", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "参数解析错误或参数缺失",
		})
		return
	}
	if requestFileDto.PathName == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "项目文件路径不能为空",
		})
		return
	}

	fileManagerService := service.NewFileManagerService()
	err = fileManagerService.SetProjectPath(requestFileDto)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("项目路径%s不存在，%s", requestFileDto.PathName, err.Error()),
		})
		return
	}
	fileCompletePath := fileManagerService.CurrPath + "/" + requestFileDto.PathName
	extName := path.Ext(fileCompletePath)
	if extName == "" {
		// 说明是目录
		isDir := fileManagerService.CheckProjectIsDir(requestFileDto.PathName)
		if !isDir {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"success": false,
				"message": fmt.Sprintf("项目路径%s不存在", requestFileDto.PathName),
			})
			return
		}
	} else {
		// 说明是文件
		_, err = fileManagerService.CheckProjectIsFile(requestFileDto.PathName)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"success": false,
				"message": fmt.Sprintf("项目文件路径%s不存在", requestFileDto.PathName),
			})
			return
		}
	}
	err = fileManagerService.CompressFileOrFolder(requestFileDto.PathName, extName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("压缩文件或文件夹失败"),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"success": true,
		"message": fmt.Sprintf("压缩文件或文件夹成功"),
	})
}

// DecompressionFile 解压文件
func DecompressionFile(c *gin.Context) {
	var requestFileDto fileManager.RequestDto
	err := c.ShouldBindJSON(&requestFileDto)
	if err != nil {
		G.Logger.Errorf("参数解析错误，具体原因:[%s]", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "参数解析错误或参数缺失",
		})
		return
	}
	if requestFileDto.PathName == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": false,
			"message": "项目文件路径不能为空",
		})
		return
	}

	fileManagerService := service.NewFileManagerService()
	err = fileManagerService.SetProjectPath(requestFileDto)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("项目路径%s不存在，%s", requestFileDto.PathName, err.Error()),
		})
		return
	}
	fileCompletePath := fileManagerService.CurrPath + "/" + requestFileDto.PathName
	extName := path.Ext(fileCompletePath)
	if extName == "" || extName != ".zip" {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": "只支持解压zip文件",
		})
		return
	}
	// 校验文件是否存在
	_, err = fileManagerService.CheckProjectIsFile(requestFileDto.PathName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("项目文件路径%s不存在", requestFileDto.PathName),
		})
		return
	}
	err = fileManagerService.DecompressionFile(requestFileDto.PathName, extName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": false,
			"message": fmt.Sprintf("解压文件失败"),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"success": true,
		"message": fmt.Sprintf("解压文件成功"),
	})
}

// GetRealTimeLog 查看实时日志
func GetRealTimeLog(c *gin.Context) {
	// 升级http为websocket连接
	//upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	//	return true
	//}}
	//w := c.Writer
	//r := c.Request
	//conn, err := upgrader.Upgrade(w, r, nil)
	//if err != nil {
	//	// 连接升级失败，返回错误
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//authSuccess := false
	//// 处理websocket请求
	//for {
	//	// 读取客户端发送的消息
	//	_, msg, err := conn.ReadMessage()
	//	if err != nil {
	//		// 读取消息失败，断开连接
	//		break
	//	}
	//	fmt.Println("读取到客户端发送的消息:", string(msg))
	//	if strings.Contains(string(msg), "x-ws-token") {
	//		var wsMap = make(map[string]any, 1)
	//		fmt.Println("是认证消息")
	//		err = json.Unmarshal(msg, &wsMap)
	//		if err != nil {
	//			authSuccess = false
	//		}
	//		token := wsMap["x-ws-token"].(string)
	//		_, err = utils.ParseToken(token)
	//		if err != nil {
	//			authSuccess = false
	//		} else {
	//			authSuccess = true
	//		}
	//	}
	//	if authSuccess {
	//		fmt.Println("身份认证成功, 开始读取实时日志")
	//		err = conn.WriteMessage(websocket.TextMessage, []byte("身份认证成功，开始读取实时日志"))
	//		if err != nil {
	//			// 发送消息失败，断开连接
	//			break
	//		}
	//
	//	} else {
	//		fmt.Println("身份认证失败")
	//		err = conn.WriteMessage(websocket.TextMessage, []byte("身份认证失败"))
	//		if err != nil {
	//			// 发送消息失败，断开连接
	//			break
	//		}
	//		// 认证失败则直接端开连接
	//		conn.Close()
	//	}
	//}
	//conn.Close()
	//

	managerService := service.NewFileManagerService()
	managerService.GetRealLog()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
