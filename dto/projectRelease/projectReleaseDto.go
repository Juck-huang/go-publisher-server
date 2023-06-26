package projectRelease

type ResponseDto struct {
	Id              uint   `json:"id"`
	ProjectId       int64  `json:"projectId"`
	ProjectEnvId    int64  `json:"projectEnvId"`
	ProjectTypeId   int64  `json:"projectTypeId"`
	BuildScriptPath string `json:"buildScriptPath"`
	Params          string `json:"params"`
}
