package service

import (
	"hy.juck.com/go-publisher-server/dto/project"
	project2 "hy.juck.com/go-publisher-server/model/project"
	user2 "hy.juck.com/go-publisher-server/model/user"
)

type ProjectService struct {
}

func NewProjectService() *ProjectService {
	return &ProjectService{}
}

// GetProjectList 获取项目列表
func (o *ProjectService) GetProjectList(username string) (projectDtos []project.ResponseDto) {
	var projects []project2.Project
	// 根据当前登录用户获取当前用户所属的项目列表
	var user user2.User
	G.DB.Debug().Where("username = ?", username).Find(&user)
	if user.ID == 0 {
		return projectDtos
	}
	G.DB.Debug().Table("user_project").Select("project.*").
		Joins("left join project on user_project.project_id = project.id").
		Where("user_project.user_id = ?", user.ID).Scan(&projects)
	for _, p := range projects {
		var projectD = project.ResponseDto{
			Id:         p.ID,
			CreateDate: p.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdateDate: p.UpdatedAt.Format("2006-01-02 15:04:05"),
			Name:       p.Name,
		}
		projectDtos = append(projectDtos, projectD)
	}
	return projectDtos
}
