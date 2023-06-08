package service

import (
	"errors"
	"fmt"
	"os/exec"
)

type DatabaseService struct {
	DbName string
}

func NewDatabaseService(dbName string) *DatabaseService {
	return &DatabaseService{
		DbName: dbName,
	}
}

// CheckMysqlDump 检查mysqldump状态
func (o *DatabaseService) CheckMysqlDump() error {
	cmd := exec.Command("mysqldump --version && echo $?")
	out, err := cmd.CombinedOutput()
	if err != nil {
		G.Logger.Errorf("未检测到mysqldump,请先安装，具体原因:[%s]", err.Error())
		return errors.New("导出失败")
	}
	G.Logger.Errorf("执行检测mysqldump结果: [%s]", out)
	if string(out) == "0" {
		return nil
	}
	return errors.New("导出失败")
}

func (o *DatabaseService) CopyFile(originFilePath string, targetFilePath string) error {
	command := fmt.Sprintf("copy %s %s && echo $?", originFilePath, targetFilePath)
	cmd := exec.Command(command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		G.Logger.Errorf("脚本执行失败，失败原因:[%s]", err.Error())
		return err
	}
	if string(out) == "0" {
		return nil
	}
	return errors.New("导出失败")
}

// HandleMysqlDump 导出或备份数据库
func (o *DatabaseService) HandleMysqlDump(tempPath string, ignoreTables ...string) error {
	// 执行脚本导出数据库
	exportScript := "mysqldump -h" + G.C.DB.Mysql.Host + " -P" + G.C.DB.Mysql.Port + " " +
		"-u" + G.C.DB.Mysql.Username + " -p" + G.C.DB.Mysql.Password + " " + o.DbName
	if len(exportScript) > 0 {
		for _, ignoreTable := range ignoreTables {
			exportScript += " --ignore-table=" + o.DbName + "." + ignoreTable
		}
	}
	exportScript += " > " + tempPath
	G.Logger.Info("生成的数据库脚本为：", exportScript)
	// 开始执行脚本
	cmd := exec.Command(exportScript)
	err := cmd.Start()
	if err != nil {
		G.Logger.Errorf("脚本执行失败，失败原因:[%s]", err.Error())
		return err
	}
	return nil
}
