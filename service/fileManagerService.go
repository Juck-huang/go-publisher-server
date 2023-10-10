package service

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"hy.juck.com/go-publisher-server/common"
	"hy.juck.com/go-publisher-server/dto/fileManager"
	"hy.juck.com/go-publisher-server/model/project"
)

type FileManagerService struct {
	CurrPath string // 项目路径
}

func NewFileManagerService() *FileManagerService {
	return &FileManagerService{}
}

// SetProjectPath 设置项目对应文件路径
func (o *FileManagerService) SetProjectPath(projectFileRequest fileManager.RequestDto) error {
	var projectDir project.ProjectDir
	G.DB.Where("project_id=? and project_env_id=? and project_type_id=?", projectFileRequest.ProjectId,
		projectFileRequest.ProjectEnvId, projectFileRequest.ProjectTypeId).First(&projectDir)
	if projectDir.ID == 0 {
		return errors.New("项目目录或文件夹不存在")
	}
	o.CurrPath = projectDir.ProjectPath
	return nil
}

// GetFileList 获取文件列表
func (o *FileManagerService) GetFileList(pathName string) (fileManagerDtos []fileManager.ResponseDto, err error) {
	readDir, err := os.ReadDir(o.CurrPath + "/" + pathName)
	if err != nil {
		return fileManagerDtos, err
	}
	// 如果读取到长度为0，则直接返回
	if len(readDir) == 0 {
		return fileManagerDtos, nil
	}
	for _, a := range readDir {
		fileInfo, err := a.Info()
		if err != nil {
			return nil, err
		}
		fileCompletePath := o.CurrPath + "/" + fileInfo.Name()
		var extName string // 文件才有扩展名
		fileType := "folder"
		if !fileInfo.IsDir() {
			extName = strings.Replace(path.Ext(fileCompletePath), ".", "", 1)
			fileType = "file"
		}
		updateDate := fileInfo.ModTime().Format("2006-01-02 15:04:05")
		fileSizeOrigin := fileInfo.Size()
		fileSize := ""
		if fileSizeOrigin >= 0 && fileSizeOrigin < 1024 {
			// 显示B
			fileSize = fmt.Sprintf("%.2fB", float64(fileSizeOrigin))
		} else if fileSizeOrigin >= 1024 && fileSizeOrigin < 1024*1024 {
			// 显示KB
			fileSize = fmt.Sprintf("%.2fKB", float64(fileSizeOrigin)/1024)
		} else if fileSizeOrigin >= 1024*1024 && fileSizeOrigin < 1024*1024*1024 {
			// 显示MB
			fileSize = fmt.Sprintf("%.2fMB", float64(fileSizeOrigin)/(1024*1024))
		} else if fileSizeOrigin >= 1024*1024*1024 && fileSizeOrigin < 1024*1024*1024*1024 {
			// 显示GB
			fileSize = fmt.Sprintf("%vGB", fileSizeOrigin/(1024*1024*1024))
		} else {
			fileSize = fmt.Sprintf("%.2fB", float64(fileSizeOrigin))
		}
		var fm = fileManager.ResponseDto{
			Name:       fileInfo.Name(),
			Size:       fileSize,
			Type:       fileType,
			ExtName:    extName,
			UpdateDate: updateDate,
		}
		fileManagerDtos = append(fileManagerDtos, fm)
		// 根据类型排序，文件夹放在最前面
		sort.SliceStable(fileManagerDtos, func(i, j int) bool {
			return fileManagerDtos[i].Type > fileManagerDtos[j].Type
		})
	}
	return fileManagerDtos, nil
}

// CheckProjectIsDir 校验项目完整路径名称
func (o *FileManagerService) CheckProjectIsDir(pathName string) bool {
	compPath := o.CurrPath + "/" + pathName
	G.Logger.Infof("当前项目文件完整文件夹名称:[%s]", compPath)
	s, err := os.Stat(compPath)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// CheckProjectIsFile 检查传递的项目是否是文件
func (o *FileManagerService) CheckProjectIsFile(pathName string) (bool, error) {
	compPath := o.CurrPath + "/" + pathName
	G.Logger.Infof("当前项目文件完整文件路径:[%s]", compPath)
	s, err := os.Stat(compPath)
	// 如果出错，直接返回
	if err != nil {
		return false, errors.New("文件路径错误或不是文件")
	}
	// 如果是目录，则直接返回
	if s.IsDir() {
		return false, errors.New("请确认上传路径是否是文件")
	}
	// 如果是文件，则判断文件是否存在，不存在则直接返回
	return true, nil
}

// SaveFileContent 保存项目文件内容
func (o *FileManagerService) SaveFileContent(pathName string, content string) (bool, error) {
	filePath := o.CurrPath + "/" + pathName
	openFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644) // O_WRONLY:只写入，O_TRUNC:覆写
	fmt.Println("content", content)
	if err != nil {
		G.Logger.Errorf("编辑文件失败:[%s]", err.Error())
		return false, errors.New("编辑文件失败")
	}
	defer openFile.Close()
	if _, err = io.WriteString(openFile, content); err != nil {
		G.Logger.Errorf("写入文件失败:[%s]", err.Error())
		return false, errors.New("编辑文件失败")
	}
	if err = openFile.Sync(); err != nil {
		G.Logger.Errorf("同步文件失败:[%s]", err.Error())
		return false, errors.New("编辑文件失败")
	}
	return true, nil
}

// RemoveFile 删除文件或文件夹
func (o *FileManagerService) RemoveFile(pathName string) error {
	rmPath := o.CurrPath + "/" + pathName
	fmt.Println("rmPath", rmPath)
	err := os.RemoveAll(rmPath)
	if err != nil {
		G.Logger.Errorf("删除文件或文件夹[%s]失败，失败原因:[%s]", pathName, err.Error())
		return errors.New("删除文件或文件夹失败")
	}
	return nil
}

// AddFolder 新建文件夹
func (o *FileManagerService) AddFolder(pathName string, addFolderName string) error {
	addPath := o.CurrPath + "/" + pathName + "/" + addFolderName
	err := os.Mkdir(addPath, os.ModePerm)
	if err != nil {
		G.Logger.Errorf("创建文件夹失败：" + err.Error())
		return errors.New("创建文件夹失败")
	}
	return nil
}

// AddFile 新建文件
func (o *FileManagerService) AddFile(pathName string, addFileName string) error {
	addPath := o.CurrPath + "/" + pathName + "/" + addFileName
	file, err := os.Create(addPath)
	if err != nil {
		G.Logger.Errorf("创建文件失败：" + err.Error())
		return errors.New("创建文件失败")
	}
	defer file.Close()
	return nil
}

// ReNameFile 重命名文件或文件夹
func (o *FileManagerService) ReNameFile(pathName string, newFileName string) error {
	oldPath := o.CurrPath + "/" + pathName
	newPath := o.CurrPath + "/" + newFileName
	err := os.Rename(oldPath, newPath)
	if err != nil {
		G.Logger.Errorf("重命名文件夹或文件失败：" + err.Error())
		return err
	}
	return nil
}

// MoveOrCopyFile 移动或复制文件或文件夹
func (o *FileManagerService) MoveOrCopyFile(operationType int, extName string, srcPath string, targetPath string) error {
	srcPath = o.CurrPath + "/" + srcPath
	targetPath = o.CurrPath + "/" + targetPath
	//fmt.Println("srcPath", srcPath)
	//fmt.Println("targetPath", targetPath)
	switch operationType {
	case 1:
		// 移动
		err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)
		if err != nil {
			G.Logger.Errorf("创建文件夹或文件失败：" + err.Error())
			return err
		}
		err = os.Rename(srcPath, targetPath)
		if err != nil {
			G.Logger.Errorf("移动文件夹或文件失败：" + err.Error())
			return err
		}
	case 2:
		// 操作
		if extName == "" {
			// 是文件夹
			sss := strings.Split(targetPath, "/")
			targetPath = strings.Join(sss[0:len(sss)-1], "/")
			command := fmt.Sprintf("cp -r %s %s && echo '复制文件夹成功'", srcPath, targetPath)
			_, err := common.ExecCommand(true, "-c", command)
			if err != nil {
				G.Logger.Errorf("复制文件夹失败：" + err.Error())
				return err
			}
		} else {
			// 是文件
			srcOpen, err := os.Open(srcPath)
			if err != nil {
				G.Logger.Errorf("复制文件失败：" + err.Error())
				return err
			}
			defer srcOpen.Close()
			targetOpen, err := os.Create(targetPath)
			if err != nil {
				G.Logger.Errorf("复制文件失败：" + err.Error())
				return err
			}
			defer targetOpen.Close()
			_, err = io.Copy(targetOpen, srcOpen)
			if err != nil {
				G.Logger.Errorf("复制文件失败：" + err.Error())
				return err
			}
		}
	default:
		return errors.New("操作类型不支持")
	}
	return nil
}

// CompressFileOrFolder 压缩文件夹或目录
func (o *FileManagerService) CompressFileOrFolder(pathName string, extName string) error {
	splitNames := strings.Split(pathName, "/")
	compressPath := o.CurrPath + "/" + strings.Join(splitNames[0:len(splitNames)-1], "/")
	var compressSrcName = splitNames[len(splitNames)-1]
	var compressTargetName string
	if extName == "" {
		// 说明是文件夹
		compressTargetName = compressSrcName + ".zip"
	} else {
		// 说明是文件
		compressTargetName = strings.Split(compressSrcName, extName)[0] + ".zip"
	}
	command := fmt.Sprintf("cd %s && zip -qr %s %s && echo $?", compressPath, compressTargetName, compressSrcName)
	out, err := common.ExecCommand(true, "-c", command)
	if err != nil {
		G.Logger.Errorf("压缩文件或文件夹失败:[%s]", err.Error())
		return errors.New("压缩文件或文件夹失败")
	}
	if out != "0" {
		G.Logger.Errorf("压缩文件或文件夹失败: 执行命令有误")
		return errors.New("压缩文件或文件夹失败")
	}
	return nil
}

// DecompressionFile 解压文件
func (o *FileManagerService) DecompressionFile(pathName string, extName string) error {
	splitNames := strings.Split(pathName, "/")
	deCompressPath := o.CurrPath + "/" + strings.Join(splitNames[0:len(splitNames)-1], "/")
	var deCompressSrcName = splitNames[len(splitNames)-1]
	command := fmt.Sprintf("cd %s && unzip -q %s && echo $?", deCompressPath, deCompressSrcName)
	out, err := common.ExecCommand(true, "-c", command)
	if err != nil {
		G.Logger.Errorf("解压文件失败:[%s]", err.Error())
		return errors.New("解压文件失败")
	}
	if out != "0" {
		G.Logger.Errorf("解压文件失败: 执行命令有误")
		return errors.New("解压文件失败")
	}
	return nil
}

// HandleMergeFile 合并文件
func (o *FileManagerService) HandleMergeFile(chunkPath string, targetFilePath string) error {
	dirEntries, err := os.ReadDir(chunkPath)
	if err != nil {
		G.Logger.Errorf("合并文件失败[%s]", err.Error())
		return errors.New(err.Error())
	}
	targetFile, err := os.OpenFile(targetFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm|os.ModeAppend)
	defer targetFile.Close()
	if err != nil {
		G.Logger.Errorf("合并文件失败[%s]", err.Error())
		return errors.New(err.Error())
	}
	G.Logger.Infof("共%d个chunk文件需要合并\n", len(dirEntries))
	// 需要先把文件根据切片名称进行升序
	sort.SliceStable(dirEntries, func(i, j int) bool {
		fileName := dirEntries[i].Name()
		fileNameSplit := strings.Split(fileName, "-")
		nextFileName := dirEntries[j].Name()
		nextFileNameSplit := strings.Split(nextFileName, "-")
		atoi, err := strconv.Atoi(fileNameSplit[len(fileNameSplit)-1])
		if err != nil {
			G.Logger.Errorf("合并文件失败[%s]", err.Error())
			return false
		}
		atoiNext, err := strconv.Atoi(nextFileNameSplit[len(nextFileNameSplit)-1])
		if err != nil {
			G.Logger.Errorf("合并文件失败[%s]", err.Error())
			return false
		}
		return atoi < atoiNext
	})
	for index, entry := range dirEntries {
		openFile, err := os.OpenFile(chunkPath+"/"+entry.Name(), os.O_RDONLY, os.ModePerm)
		if err != nil {
			G.Logger.Errorf("合并文件失败[%s]", err.Error())
			return errors.New(err.Error())
		}
		data, err := ioutil.ReadAll(openFile)
		if err != nil {
			G.Logger.Errorf("合并文件失败[%s]", err.Error())
			return errors.New(err.Error())
		}
		targetFile.Write(data)
		openFile.Close()
		G.Logger.Infof("合并第%d个文件成功,文件名称%s\n", index+1, entry.Name())
	}
	return nil
}

// CheckFileHash 校验合并文件完成生成的hash与原文件hash是否一致
func (o *FileManagerService) CheckFileHash(filePath string, originHash string) (bool, error) {
	openFile, err := os.Open(filePath)
	defer openFile.Close()
	if err != nil {
		return false, err
	}
	//buf := make([]byte, 1024)
	md5h := md5.New()
	_, err = io.Copy(md5h, openFile)
	if err != nil {
		return false, err
	}
	//for {
	//	n, err := openFile.Read(buf)
	//	if err != nil {
	//		if err == io.EOF {
	//			break
	//		}
	//		fmt.Println("读取文件失败")
	//	}
	//	md5h.Write(buf[:n])
	//}
	md5Str := hex.EncodeToString(md5h.Sum(nil))
	fmt.Println("文件hash:", originHash, md5Str)
	return originHash == md5Str, nil
}

func (o *FileManagerService) GetRealLog() {
	filePath := "/Users/mac/Downloads/apps/emerge/stec-emerge-service/default/logs/123.log"
	lineNum := 5
	f, err := os.Open(filePath)
	if err != nil {
		os.Exit(1)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		os.Exit(1)
	}
	var offset int64
	if fi.Size() > int64(1024*lineNum) {
		offset = fi.Size() - int64(1024*lineNum)
	}
	_, err = f.Seek(offset, 0)
	if err != nil {
		os.Exit(1)
	}
	reader := bufio.NewReader(f)
	for {
		line, a, err := reader.ReadLine()
		fmt.Println("line", string(line), a)
		if err != nil {
			fmt.Println("err", err)
			if err.Error() == "EOF" {
				time.Sleep(time.Second * 5)
				continue
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		fmt.Println("line ", string(line))
	}
}
