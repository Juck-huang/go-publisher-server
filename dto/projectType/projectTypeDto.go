package projectType

type ResponseDto struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	ProjectId int64  `json:"projectId"`
	ParentId  int64  `json:"parentId,omitempty"`
	IsLeaf    int64  `json:"isLeaf"`
	TreeId    string `json:"treeId"`
	TreeLevel int64  `json:"treeLevel"`
}
