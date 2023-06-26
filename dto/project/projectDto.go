package project

type ResponseDto struct {
	Id           uint   `json:"id"`
	CreateDate   string `json:"createDate"`
	UpdateDate   string `json:"updateDate"`
	Name         string `json:"name"`
	ProjectEnvId uint   `json:"projectEnvId"`
}
