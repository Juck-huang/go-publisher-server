package service

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	user2 "hy.juck.com/go-publisher-server/model/user"
)

type DatabaseService struct {
	DbName      string
	dbPrivilege string
}

func NewDatabaseService(dbName string) *DatabaseService {
	return &DatabaseService{
		DbName: dbName,
	}
}

// CheckMysqlDump 检查mysqldump状态
func (o *DatabaseService) CheckMysqlDump() error {
	command := "source /etc/profile && mysqldump --version"
	_, err := o.execCommand("检查mysqldump", command)
	if err != nil {
		G.Logger.Errorf("未检测到mysqldump,请先安装，具体原因:[%s]", err.Error())
		return errors.New("导出或备份数据库失败")
	}
	return nil
}

// CheckMysql 检查mysql状态
func (o *DatabaseService) CheckMysql() error {
	command := "source /etc/profile && mysql --version"
	_, err := o.execCommand("检查mysql", command)
	if err != nil {
		G.Logger.Errorf("未检测到mysql,请先安装，具体原因:[%s]", err.Error())
		return errors.New("导出或备份数据库失败")
	}
	return nil
}

// CopyFile 复制文件
func (o *DatabaseService) CopyFile(originFilePath string, targetFilePath string) error {
	_, err := os.Stat(originFilePath)
	if os.IsNotExist(err) {
		// 说明文件不存在存在
		G.Logger.Errorf("源文件不存在无法复制，失败原因:[%s]", err.Error())
		return errors.New("导出或备份数据失败")
	}
	command := fmt.Sprintf("cp %s %s", originFilePath, targetFilePath)
	_, err = o.execCommand("复制文件", command)
	if err != nil {
		// G.Logger.Errorf("复制文件脚本执行失败，失败原因:[%s]", err.Error())
		return errors.New("导出或备份数据库失败")
	}
	return nil
}

// HandleTotalMysqlDump HandleMysqlDump 导出或备份整个数据库
func (o *DatabaseService) HandleTotalMysqlDump(tempPath string, ignoreTables ...string) error {
	// 执行脚本导出数据库
	//exportScript := "mysqldump -h" + G.C.Ops.Mysql.Host + " -P" + G.C.Ops.Mysql.Port + " " +
	//	"-u" + G.C.Ops.Mysql.Username + " -p'" + G.C.Ops.Mysql.Password + "' " + o.DbName
	mysqlData := G.C.Ops.Mysql
	exportScript := fmt.Sprintf("mysqldump -h %s -u %s -p'%s' -P %s %s", mysqlData.Host, mysqlData.Username, mysqlData.Password, mysqlData.Port, o.DbName)
	if len(ignoreTables) > 0 {
		for _, ignoreTable := range ignoreTables {
			exportScript += " --ignore-table=" + o.DbName + "." + ignoreTable
		}
	}
	exportScript += " > " + tempPath
	_, err := o.execCommand("导出或备份", exportScript)
	if err != nil {
		G.Logger.Errorf("导出或备份数据脚本执行失败，失败原因:[%s]", err.Error())
		return err
	}
	return nil
}

// 执行命令
func (o *DatabaseService) execCommand(commandName string, command string) ([]string, error) {
	command += " && echo $?"
	G.Logger.Infof("[%s]正在执行命令: [%s]", commandName, command)
	// 开始执行脚本
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.CombinedOutput()
	outStr := strings.TrimSpace(string(out))
	strs := strings.Split(outStr, "\n")
	var newArr []string
	if err != nil {
		G.Logger.Errorf("[%s]执行脚本失败，具体状态:[%s], 失败原因: [%s]", commandName, err.Error(), strs)
		return newArr, errors.New(fmt.Sprintf("执行sql脚本失败,具体原因：%s", strs))
	}
	G.Logger.Infof("[%s]脚本执行结果,状态:[ %s]", commandName, outStr)
	// 如果只有一位则直接判断
	if len(strs) == 1 && strs[0] == "0" {
		return strs, nil
	}
	newArr = append(newArr, strs[0:len(strs)-1]...)
	// 判断数组最后一位是否是0，为0则代表脚本执行成功
	if strs[len(strs)-1] == "0" {
		return newArr, nil
	}

	return newArr, errors.New("执行脚本失败")
}

// ZipFile 压缩文件
func (o *DatabaseService) ZipFile(currentPath string, originPath string, targetPath string) error {
	command := fmt.Sprintf("cd %s && zip -qj %s %s", currentPath, targetPath, originPath)
	_, err := o.execCommand("压缩文件", command)
	if err != nil {
		// G.Logger.Errorf("压缩文件脚本执行失败，失败原因:[%s]", err.Error())
		return err
	}
	return nil
}

// SingleExportTables 单独导出多个表
func (o *DatabaseService) SingleExportTables(tempPath string, tableNames ...string) error {
	// mysqldump -h127.0.0.1 -P3306 -uroot -p123456 stec-cdsa --tables sys_user > cdsa
	mysqlData := G.C.Ops.Mysql
	command := fmt.Sprintf("mysqldump -h %s -u%s -p'%s' -P %s %s", mysqlData.Host, mysqlData.Username, mysqlData.Password, mysqlData.Port, o.DbName)
	if len(tableNames) > 0 {
		command += " --tables"
		for _, tableName := range tableNames {
			command += " " + tableName
		}
	}
	command += " > " + tempPath
	_, err := o.execCommand("导出或备份", command)
	if err != nil {
		// G.Logger.Errorf("导出表脚本执行失败，失败原因:[%s]", err.Error())
		return err
	}
	return nil
}

// DynamicExecSql 动态执行sql
func (o *DatabaseService) DynamicExecSql(sql string, username string) (map[string]any, error) {

	var dataMap = make(map[string]any, 1)
	err := o.SetDbPrivilege(username)
	if err != nil {
		return dataMap, err
	}
	// 从数据库获取当前登录用户数据库的权限,若该用户同时包含读和写权限，则返回写权限，否则就是单独的权限
	sql = strings.ReplaceAll(sql, "\"", "\\\"")
	// command := fmt.Sprintf("mysql --login-path=%s %s -e \"%s\"", o.dbPrivilege, o.DbName, sql)
	mysqlData := G.C.Ops.Mysql
	command := fmt.Sprintf("mysql -u%s -p'%s' -h%s -P%s %s -e \"%s\"", mysqlData.Username,
		mysqlData.Password, mysqlData.Host, mysqlData.Port, o.DbName, sql)
	resultList, err := o.execCommand("动态执行sql", command)
	if err != nil {
		// 如果有错误，则返回格式还是之前的格式
		dataMap["title"] = []string{"err_msg"}
		var errInfo = err.Error()
		if strings.Contains(err.Error(), "1142") {
			errInfo = "您没有操作权限"
		}
		dataMap["content"] = []any{[]string{errInfo}}
		return dataMap, errors.New(errInfo)
	}
	if len(resultList) > 0 {
		resultList = resultList[1:]
	}
	var dataContentList []any
	for i, result := range resultList {
		var dataList []string
		result = strings.ReplaceAll(result, "NULL", "")
		strList := strings.Split(result, "\t")
		// 说明是标题
		if i == 0 {
			dataMap["title"] = strList
			continue
		}
		dataList = append(dataList, strList...)
		dataContentList = append(dataContentList, dataList)
	}
	if len(dataContentList) > 0 {
		dataMap["content"] = dataContentList
	}
	return dataMap, nil
}

// GetDbAndTableList 获取可操作的数据库列表
func (o *DatabaseService) GetDbAndTableList(ignoreDbs []string) (map[string]any, error) {
	// sql:mysql -uroot -pcjxx2022 -h127.0.0.1 -P3306 -e "SHOW DATABASES WHERE \`Database\` NOT IN ('information_schema', 'sys', 'performance_schema', 'mysql')"
	mysqlData := G.C.Ops.Mysql
	command := fmt.Sprintf("mysql -u%s -p'%s' -h%s -P%s %s -e \"SHOW DATABASES WHERE \\`Database\\` NOT IN (", mysqlData.Username, mysqlData.Password, mysqlData.Host, mysqlData.Port, o.DbName)
	for i, db := range ignoreDbs {
		if i == len(ignoreDbs)-1 {
			command += "'" + db + "')\""
		} else {
			command += "'" + db + "',"
		}
	}

	var dbTableMap = make(map[string]any, 1)
	dataList, err := o.execCommand("获取数据库列表", command)
	if err != nil {
		return dbTableMap, err
	}
	if len(dataList) > 0 {
		dataList = dataList[2:]
	}

	for _, db := range dataList {
		tables, err := o.__getTableList(db)
		if err != nil {
			G.Logger.Error("获取数据库列表失败，", err)
			return dbTableMap, errors.New("获取数据库列表失败")
		}
		dbTableMap[db] = tables
	}

	return dbTableMap, nil
}

func (o *DatabaseService) __getTableList(dbName string) ([]string, error) {
	// mysql -uroot -pcjxx2022 -h127.0.0.1 -P3306 stec_bytd -e "SHOW TABLES;"
	mysqlData := G.C.Ops.Mysql
	command := fmt.Sprintf("mysql -u%s -p'%s' -h%s -P%s %s -e \"SHOW TABLES;\"", mysqlData.Username, mysqlData.Password, mysqlData.Host, mysqlData.Port, dbName)
	dataList, err := o.execCommand("获取数据表列表", command)
	if err != nil {
		return dataList, err
	}
	if len(dataList) > 0 {
		dataList = dataList[2:]
	}
	return dataList, nil
}

// ExecSqlFile 执行sql文件
func (o *DatabaseService) ExecSqlFile(tempPath string, username string) error {
	// mysql -uroot  -P3306 -p123456 -h127.0.0.1 数据库名称 < sql脚本全路径
	err := o.SetDbPrivilege(username)
	if err != nil {
		return err
	}
	mysqlData := G.C.Ops.Mysql
	command := fmt.Sprintf("mysql -u%s -p'%s' -h%s -P%s %s < %s", mysqlData.Username, mysqlData.Password, mysqlData.Host, mysqlData.Port, o.DbName, tempPath)
	_, err = o.execCommand("执行sql文件", command)
	if err != nil {
		//G.Logger.Errorf("执行sql文件失败,失败原因:[%s]", err.Error())
		var errInfo = err.Error()
		if strings.Contains(err.Error(), "1142") {
			errInfo = "您没有操作权限"
		}
		return errors.New("执行sql脚本失败：" + errInfo)
	}
	return nil
}

// SetDbPrivilege 设置数据库权限
func (o *DatabaseService) SetDbPrivilege(username string) error {
	// 从数据库获取当前登录用户数据库的权限,若该用户同时包含读和写权限，则返回写权限，否则就是单独的权限
	var user user2.User
	G.DB.Where("username = ?", username).First(&user)
	if user.ID == 0 {
		return errors.New("当前用户不存在")
	}
	var privilegeCodes []string
	G.DB.Debug().Table("user_privilege").Select("privilege.code as privilegeCode").
		Joins("left join privilege on privilege.id = user_privilege.privilege_id").
		Where("user_privilege.user_id = ? and privilege.type = 1", user.ID).Scan(&privilegeCodes)
	var loginPath string
	if len(privilegeCodes) == 0 {
		return errors.New("您没有操作权限")
	} else if len(privilegeCodes) == 1 {
		switch privilegeCodes[0] {
		case "read":
			loginPath = privilegeCodes[0]
		case "write":
			loginPath = privilegeCodes[0]
		default:
			return errors.New("您没有操作权限")
		}
	} else if len(privilegeCodes) == 2 {
		for _, privilegeCode := range privilegeCodes {
			if "read" == privilegeCode {
				loginPath = "read"
				continue
			}
			if "write" == privilegeCode {
				// 说明是写权限
				loginPath = "write"
				break
			}
		}
	}
	o.dbPrivilege = loginPath
	return nil
}
