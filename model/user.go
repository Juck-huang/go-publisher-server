package model

// User 用户结构体
type User struct {
	Id            int64  `json:"id"`            // id
	Name          string `json:"name"`          // 名称
	Username      string `json:"username"`      // 用户名
	Password      string `json:"password"`      // 密码，保存argon2加密后的密码
	State         int    `json:"state"`         // 用户状态 0为正常， 1为禁用， 默认为0
	CreateDate    string `json:"createDate"`    // 记录创建日期
	LastLoginDate string `json:"lastLoginDate"` // 上次登录日期
	LoginIp       string `json:"loginIp"`       // 用户登录的ip地址
	LoginIpPlace  string `json:"loginIpPlace"`  // 登录ip归属地
	Avatar        string `json:"avatar"`        // 个人头像地址
}
