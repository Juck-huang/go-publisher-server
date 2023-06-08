package model

type User struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	CreateDate    string `json:"createDate"`
	LastLoginDate string `json:"lastLoginDate"`
}
