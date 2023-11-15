package auth

type AuthRequest struct {
	Uid          string   `json:"uid" binding:"required"`          // 服务器唯一标识
	ServerIpList []string `json:"serverIpList" binding:"required"` // 服务器公网ip列表
	ProjectId    string   `json:"projectId" binding:"required"`    // 服务器公网ip
}
