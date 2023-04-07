# 应用发布系统 

### 流程图

![img.png](img.png)

## 编译linux下运行文件命令
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
## 编译windows下运行文件命令
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build
