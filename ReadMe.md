# 应用发布系统

## 编译linux下运行文件命令
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
## 编译windows下运行文件命令
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build

## rsa加密
 * 生成私钥命令：openssl genrsa -out rsa_1024_priv.pem 1024
 * 生成公钥命令：openssl rsa -pubout -in rsa_1024_priv.pem -out rsa_1024_pub.pem

## 需求通知
### 0706需求(只可操作对应的项目目录）
 * 增加项目目录在线浏览
 * 增加项目文件在线编辑
 * 增加项目文件上传
 * 增加项目文件下载

### 后续需求
 * 数据库备份目录、版本说明及版本录入数据表
 * 项目备份目录、版本说明及版本录入数据表
 * 备份的数据库历史版本浏览、下载和回滚
 * 备份的项目历史版本浏览、下载和回滚
 * 增加查询后端jar包运行状态、停止jar包和启动jar包
 * 增加查看jar包运行日志
 * 白名单功能尚未真正实现（目前临时采用nginx指定ip访问）
 * token需要跟每一台服务器做绑定（不同服务器不通用）
