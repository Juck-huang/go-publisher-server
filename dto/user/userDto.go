package user

type ResponseDto struct {
	Id             uint               `json:"id"`         // id
	Name           string             `json:"name"`       // 名称
	Avatar         string             `json:"avatar"`     // 个人头像地址
	State          int                `json:"state"`      // 用户状态，0.正常， 1.禁用
	Username       string             `json:"username"`   // 用户名
	CreateDate     string             `json:"createDate"` // 创建日期
	UpdateDate     string             `json:"updateDate"` // 更新日期
	UserPrivileges []UserPrivilegeDto `json:"userPrivileges"`
}

type UserPrivilegeDto struct {
	PrivilegeType string `json:"privilegeType"` // 权限类型 0.用户，1.数据库
	PrivilegeCode string `json:"privilegeCode"` // 权限编码
}
