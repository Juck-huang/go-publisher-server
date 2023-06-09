package service

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
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
	command := "source /etc/profile && mysqldump --version && echo $?"
	err := o.execCommand("检查mysqldump", command)
	if err != nil {
		G.Logger.Errorf("未检测到mysqldump,请先安装，具体原因:[%s]", err.Error())
		return errors.New("导出或备份数据库失败")
	}
	return nil
}

// 复制文件
func (o *DatabaseService) CopyFile(originFilePath string, targetFilePath string) error {
	_, err := os.Stat(originFilePath)
	if os.IsNotExist(err) {
		// 说明文件不存在存在
		G.Logger.Errorf("源文件不存在无法复制，失败原因:[%s]", err.Error())
		return errors.New("导出或备份数据失败")
	}
	command := fmt.Sprintf("cp %s %s && echo $?", originFilePath, targetFilePath)
	err = o.execCommand("复制文件", command)
	if err != nil {
		G.Logger.Errorf("复制文件脚本执行失败，失败原因:[%s]", err.Error())
		return errors.New("导出或备份数据库失败")
	}
	return nil
}

// HandleMysqlDump 导出或备份整个数据库
func (o *DatabaseService) HandleTotalMysqlDump(tempPath string, ignoreTables ...string) error {
	// 执行脚本导出数据库
	exportScript := "source /etc/profile && mysqldump -h" + G.C.DB.Mysql.Host + " -P" + G.C.DB.Mysql.Port + " " +
		"-u" + G.C.DB.Mysql.Username + " -p" + G.C.DB.Mysql.Password + " " + o.DbName
	if len(exportScript) > 0 {
		for _, ignoreTable := range ignoreTables {
			exportScript += " --ignore-table=" + o.DbName + "." + ignoreTable
		}
	}
	exportScript += " > " + tempPath + " && echo $?"
	err := o.execCommand("导出或备份", exportScript)
	if err != nil {
		G.Logger.Errorf("导出或备份数据脚本执行失败，失败原因:[%s]", err.Error())
		return errors.New("导出或备份数据库失败")
	}
	return nil
}

// 执行命令
func (o *DatabaseService) execCommand(commondName string, commond string) error {
	G.Logger.Infof("[%s]正在执行命令: [%s]", commondName, commond)
	// 开始执行脚本
	cmd := exec.Command("bash", "-c", commond)
	out, err := cmd.CombinedOutput()
	if err != nil {
		G.Logger.Errorf("[%s]执行脚本失败，失败原因:[%s]", commondName, err.Error())
		return errors.New(commondName + "执行失脚本失败")
	}
	outStr := strings.TrimSpace(string(out))
	strs := strings.Split(outStr, "\n")
	G.Logger.Infof("[%s]脚本执行结果,状态: [%s]", commondName, outStr)
	// 判断数组最后一位是否是0，为0则代表脚本执行成功
	if strs[len(strs)-1] == "0" {
		return nil
	}
	return errors.New(commondName + "执行脚本失败")
}

// 压缩文件
func (o *DatabaseService) ZipFile(currentPath string, originPath string, targetPath string) error {
	command := fmt.Sprintf("cd %s && zip -qj %s %s && echo $?", currentPath, targetPath, originPath)
	err := o.execCommand("压缩文件", command)
	if err != nil {
		G.Logger.Errorf("压缩文件脚本执行失败，失败原因:[%s]", err.Error())
		return errors.New("导出或备份数据库失败")
	}
	return nil
}

// 单独导出多个表
func (o *DatabaseService) SimpleExportTables(tempPath string, tableNames ...string) error {
	command := fmt.Sprintf("mysqldump -h%s -P%s -u%s -p%s %s", G.C.DB.Mysql.Host, G.C.DB.Mysql.Port,
		G.C.DB.Mysql.Username, G.C.DB.Mysql.Password, o.DbName)
	if len(tableNames) > 0 {
		command += " --tables"
		for _, tableName := range tableNames {
			command += " " + tableName
		}
	}
	command += " > " + tempPath + " && echo $?"
	err := o.execCommand("导出或备份", command)
	if err != nil {
		G.Logger.Errorf("导出表脚本执行失败，失败原因:[%s]", err.Error())
		return errors.New("导出表失败")
	}
	return nil
}

// DynamicExecSql 动态执行sql
func (o *DatabaseService) DynamicExecSql(sql string) error {
	command := fmt.Sprintf("mysql -u%s -p%s %s -e \"%s\";", G.C.DB.Mysql.Username, G.C.DB.Mysql.Password, o.DbName, sql)
	err := o.execCommand("动态执行sql", command)
	if err != nil {
		G.Logger.Errorf("动态执行sql脚本执行失败，失败原因:[%s]", err.Error())
		return errors.New("语句执行失败")
	}
	return nil
}
