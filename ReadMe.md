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

## curl
 * 发送post请求：curl -H "Content-Type: application/json" -X POST -d '{"username":"example","password":"test"}' http://www.baidu.com

## 需求通知
### 需求(只可操作对应的项目目录）
 * 白名单功能尚未真正实现（目前临时采用nginx指定ip访问）(20231114已实现，从数据库读取，定时把公网ip写库)
 * 增加项目文件在线编辑(20231017已实现)
 * 增加项目文件上传(20231017已实现)
 * 增加新建文件夹和文件功能，默认指定读写权限(20231017已实现)
 * 增加项目文件删除，项目文件夹删除(20231017已实现)
 * token需要跟每一台服务器做绑定（不同服务器不通用）(20231017已实现)
 * 增加项目目录在线浏览，文件和文件夹可重命名(20231017已实现)

### 后续需求
 * 增加项目目录移动到项目指定目录
 * 数据库备份目录、版本说明及版本录入数据表
 * 项目备份目录、版本说明及版本录入数据表
 * 备份的数据库历史版本浏览、下载和回滚
 * 增加查询后端jar包运行状态、停止jar包和启动jar包
 * 增加查看jar包运行日志
 * 增加项目文件夹多选、删除多选的内容(删除时需要校验所有子文件夹是否为空，均为空则可以删除)
 * 增加项目文件和文件夹下载,下载文件夹时则默认下载压缩后的文件
 * 增加监控服务器状态，通过websocket实时获取CPU、硬盘、内存、负载(1分钟、5分钟、30分钟)、运行时间、ip等信息,(第一次进来获取所有信息，后续获取实时信息)
 * 增加上述指标超过规定值系统自动报警的功能，如通过发送邮件告知
 * 增加项目管理、项目环境管理和项目类型管理页面

### 20231114新增需求
 * 应用发布：绑定用户项目名称和环境(如该用户的该项目只对应一个环境)
 * 应用发布：增加上传列表显示，支持分片上传，增加暂停、取消、删除当前已上传按钮
 * 应用管理：备份的项目历史版本浏览、下载和回滚
 * 应用发布：发布时可选择是否备份原项目，默认备份
 * 应用发布：增加发布按钮，发布版本号、发布日期，点击发布完后自动重启应用，自动弹出日志打印页（需手动关闭）
 * 应用发布：项目文件支持拖拽上传，一次只能上传一个项目文件。
 * 新增定时任务功能，支持通过系统添加定时任务，指定脚本路径和参数信息，包括开启和关闭状态，任务说明