package fileManager

type ResponseDto struct {
	Name       string `json:"name"`       // 文件名称
	Size       string `json:"size"`       // 文件大小,低于1024KB显示为KB，超过1025KB小于1024MB显示为MB，超过1024MB，显示为GB
	Type       string `json:"type"`       // 文件类型，文件夹为folder，文件为file
	ExtName    string `json:"extName"`    // 文件扩展名，文件夹返回为空，文件返回其扩展名
	UpdateDate string `json:"updateDate"` // 文件更新日期
}

type BaseRequest struct {
	ProjectId     uint   `json:"projectId" binding:"required"`     // 项目id
	ProjectEnvId  uint   `json:"projectEnvId" binding:"required"`  // 项目环境id
	ProjectTypeId uint   `json:"projectTypeId" binding:"required"` // 项目类型id
	PathName      string `json:"pathName"`                         // 项目路径名称，传入会动态查询对应的目录，初次查询默认为空
}

type RequestDto struct {
	ProjectId     uint   `json:"projectId" binding:"required"`     // 项目id
	ProjectEnvId  uint   `json:"projectEnvId" binding:"required"`  // 项目环境id
	ProjectTypeId uint   `json:"projectTypeId" binding:"required"` // 项目类型id
	PathName      string `json:"pathName"`                         // 项目路径名称，传入会动态查询对应的目录，初次查询默认为空
	FileContent   string `json:"fileContent"`                      // 编辑保存文件内容
	AddFolderName string `json:"addFolderName"`                    // 新增文件夹名称
	AddFileName   string `json:"addFileName"`                      // 新增文件名称
}

// ReNameRequestDto 重命名
type ReNameRequestDto struct {
	BaseRequest
	NewFileName string `json:"newFileName" binding:"required"` // 重命名后文件或文件夹的名称
}

type MoveOrCopyRequestDto struct {
	BaseRequest
	Type   int    `json:"type" binding:"required"` // 操作类型, 1.是移动，2是复制
	ToPath string `json:"toPath"`                  // 目标路径
}
