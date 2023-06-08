package login

// RequestDto 登录dto
type RequestDto struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}
