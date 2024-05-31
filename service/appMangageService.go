package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"hy.juck.com/go-publisher-server/common"
	"hy.juck.com/go-publisher-server/dto/application"
	model "hy.juck.com/go-publisher-server/model/application"
)

type AppManageService struct{}

func NewAppManageService() *AppManageService {
	return &AppManageService{}
}

// 获取最新app状态信息列表
func (o *AppManageService) GetAppStatusList() (applicationResponseDtos []application.ApplicationResponseDto) {
	var applications []model.Application
	err := G.DB.Debug().Find(&applications).Error
	if err != nil {
		return applicationResponseDtos
	}
	// 查找开发语言表
	var devLanguage []model.DevLanguage
	err = G.DB.Debug().Find(&devLanguage).Error
	if err != nil {
		return applicationResponseDtos
	}

	for _, app := range applications {
		var applicationResponseDto = application.ApplicationResponseDto{}
		applicationResponseDto.Id = app.ID
		applicationResponseDto.Name = app.Name
		applicationResponseDto.PackageTime = o.getPackageTime(app.AppPath)
		status, pid := o.getAppRunStatus(app.Code)
		applicationResponseDto.RunStatus = status
		for _, devLanguage := range devLanguage {
			if devLanguage.ID == app.DevLanguageId {
				applicationResponseDto.DevLanauge = devLanguage.Name
				break
			}
		}

		if status {
			startTime := o.getAppStartTime(pid)
			applicationResponseDto.StartTime = startTime
			runTime := o.getAppRunTime(pid)
			applicationResponseDto.RunTime = runTime
		}
		applicationResponseDtos = append(applicationResponseDtos, applicationResponseDto)
	}
	return applicationResponseDtos
}

// 获取运行状态
func (o *AppManageService) getAppRunStatus(code string) (bool, string) {
	command := fmt.Sprintf("pgrep -f %s", code)
	result, _ := o.execCommand(command)
	if result == "" {
		return false, ""
	}
	return true, result
}

// 获取运行时间
func (o *AppManageService) getAppRunTime(pid string) string {
	// 如果运行超过24小时，则命令返回格式：1-20:49:26，否则命令返回格式：20:49:26
	command := fmt.Sprintf("ps -p %s -o etime| awk 'NR==2'", pid)
	result, _ := o.execCommand(command)
	if result == "" {
		return ""
	}
	splites := strings.Split(strings.TrimSpace(result), ":")
	layout := "04:05"
	if len(splites) == 2 {
		// 说明是04:45格式,不超过一个小时
		parstTime, _ := time.Parse(layout, strings.TrimSpace(result))
		formattedTime := parstTime.Format("04分05秒")
		return formattedTime
	}
	splites = strings.Split(strings.TrimSpace(result), "-")
	layout = "15:04:05"
	if len(splites) == 1 {
		// 说明未超过24小时
		parstTime, _ := time.Parse(layout, splites[0])
		formattedTime := parstTime.Format("15时04分05秒")
		return formattedTime
	} else {
		// 说明超过24小时
		parstTime, _ := time.Parse(layout, splites[1])
		formattedTime := parstTime.Format("15时04分05秒")
		formattedTime = fmt.Sprintf("%s天%s", splites[0], formattedTime)
		return formattedTime
	}
}

// 获取启动时间
func (o *AppManageService) getAppStartTime(pid string) string {
	// 如果正在运行，则查询，否则不查询
	command := fmt.Sprintf("ps -p %s -o lstart | awk 'NR==2'", pid)
	result, _ := o.execCommand(command)
	command = "date -d " + "\"" + result + "\"" + " \"+%Y-%m-%d %H:%M:%S\""
	result, _ = o.execCommand(command)
	return result
}

// 获取发包时间
func (o *AppManageService) getPackageTime(packagePath string) string {
	// 如果正在运行，则查询，否则不查询
	fileInfo, err := os.Stat(packagePath)
	if err != nil {
		return ""
	}
	return fileInfo.ModTime().Format("2006-01-02 15:04:05")
}

// 开启或停止app
func (o *AppManageService) StartOrStopApp(id uint, direct string) error {
	var directName string
	if direct == "start" {
		directName = "开启"
	} else {
		directName = "停止"
	}
	var applications []model.Application
	err := G.DB.Debug().Find(&applications).Error
	if err != nil {
		return errors.New(directName + "应用失败:" + err.Error())
	}
	params := []string{}
	for _, application := range applications {
		if application.ID == id {
			params = append(params, application.ScriptPath)
			params = append(params, application.Code)
			params = append(params, strings.Split(application.Params, ",")...)
			break
		}
	}
	if len(params) == 0 {
		return errors.New(directName + "应用失败,未获取到当前应用信息")
	}
	// ./app-manage stec-emerge-service /data/apps/stec-emerge-service/default start
	params = append(params, direct)
	result, err := common.ExecCommand(true, params...)
	if err != nil {
		return errors.New(directName + "应用失败")
	}
	if result == "-1" {
		return errors.New(directName + "应用失败，或应用已" + directName)
	}
	return nil
}

// 重启app
func (o *AppManageService) RestartApp(id uint) error {
	err := o.StartOrStopApp(id, "stop")
	if err != nil {
		return errors.New("重启应用失败:" + err.Error())
	}
	err = o.StartOrStopApp(id, "start")
	if err != nil {
		return errors.New("重启应用失败")
	}
	return nil
}

// ExecCommand 执行命令,分管道和标准输出
func (o *AppManageService) execCommand(command string) (string, error) {
	G.Logger.Infof("执行的命令为：%s", command)
	cmd := exec.Command("bash", "-c", command)
	var msg string
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		G.Logger.Errorf("执行命令失败：失败原因[%s]", err.Error())
		msg = "执行命令失败"
		return msg, err
	}
	err = cmd.Start()
	if err != nil {
		G.Logger.Errorf("执行命令失败：失败原因[%s]", err.Error())
		msg = "执行命令失败"
		return msg, err
	}
	outBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		G.Logger.Errorf("执行命令失败：失败原因[%s]", err.Error())
		msg = "执行命令失败"
		return msg, err
	}
	defer stdout.Close()
	splitBytes := strings.Split(string(outBytes), "\n")
	G.Logger.Infof("执行命令结果:%s", splitBytes)
	if len(splitBytes) > 1 {
		return splitBytes[len(splitBytes)-2], nil
	} else {
		return "", errors.New("执行命令失败")
	}
}
