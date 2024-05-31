package application

type ApplicationResponseDto struct {
	Id          uint   `json:"id"`          // id
	Name        string `json:"name"`        // 应用名称
	RunStatus   bool   `json:"runStatus"`   // 运行状态, true为运行，false为停止
	PackageTime string `json:"packageTime"` // 包发布时间
	RunTime     string `json:"runTime"`     // 运行时长
	StartTime   string `json:"startTime"`   // 开启时间
	DevLanauge  string `json:"devLanauge"`  // 开发语言
}
