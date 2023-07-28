# 应用发布系统

## 编译linux下运行文件命令
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
## 编译windows下运行文件命令
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build

## rsa加密
 * 生成私钥命令：openssl genrsa -out rsa_1024_priv.pem 1024
 * 生成公钥命令：openssl rsa -pubout -in rsa_1024_priv.pem -out rsa_1024_pub.pem

## mysql config editor
 * 参考：https://blog.51cto.com/u_10950710/4973128
 * 生成：mysql_config_editor set --login-path=read -h localhost --user=root --port=3306 --password
 * 重置所有：mysql_config_editor reset
 * 移除单个：mysql_config_editor remove --login-path=auth
 * 需要给读用户加上的权限：select
 * 创建只读权限教程：https://blog.csdn.net/weixin_43573186/article/details/121607548

## 需求通知
### 0706需求(只可操作对应的项目目录）
 * 增加项目目录在线浏览，文件和文件夹可重命名、移动到项目指定目录
 * 增加项目文件在线编辑
 * 增加项目文件上传
 * 增加新建文件夹和文件功能，默认指定读写权限
 * 增加项目文件删除，项目文件夹删除(删除文件夹时需校验文件夹是否为空，为空则可以删除)
 * 增加项目文件夹多选、删除多选的内容(删除时需要校验所有子文件夹是否为空，均为空则可以删除)
 * 增加项目文件和文件夹下载,下载文件夹时则默认下载压缩后的文件

### 后续需求
 * 数据库备份目录、版本说明及版本录入数据表
 * 项目备份目录、版本说明及版本录入数据表
 * 备份的数据库历史版本浏览、下载和回滚
 * 备份的项目历史版本浏览、下载和回滚
 * 增加查询后端jar包运行状态、停止jar包和启动jar包
 * 增加查看jar包运行日志
 * 白名单功能尚未真正实现（目前临时采用nginx指定ip访问）
 * token需要跟每一台服务器做绑定（不同服务器不通用）
