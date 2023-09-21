package service

import (
	"hy.juck.com/go-publisher-server/dto/projectEnv"
	"hy.juck.com/go-publisher-server/model/project"
	user2 "hy.juck.com/go-publisher-server/model/user"
)

type ProjectEnvService struct {
}

func NewProjectEnvService() *ProjectEnvService {
	return &ProjectEnvService{}
}

// GetProjectEnvListByPUsername 通过项目id列表获取项目环境列表
func (o *ProjectEnvService) GetProjectEnvListByPUsername(username string) (projectEnvDtos []projectEnv.ResponseDto) {
	var projectEnvs []project.ProjectEnv
	var user user2.User
	G.DB.Debug().Where("username = ?", username).First(&user)
	G.DB.Debug().Table("project_env").Select("project_env.*").
		Joins("left join user_project_env on project_env.id = user_project_env.project_env_id").
		Where("user_project_env.user_id = ?", user.ID).Find(&projectEnvs)
	for _, p := range projectEnvs {
		var projectD = projectEnv.ResponseDto{
			Id:        p.ID,
			Name:      p.Name,
			ProjectId: p.ProjectId,
		}
		projectEnvDtos = append(projectEnvDtos, projectD)
	}
	return projectEnvDtos
}

// GetProjectEnvListByPId 通过项目id获取项目环境列表
func (o *ProjectEnvService) GetProjectEnvListByPId(projectId int64) (projectEnvDtos []projectEnv.ResponseDto) {
	var projectEnvs []project.ProjectEnv
	G.DB.Debug().Where("project_id = ?", projectId).Find(&projectEnvs)
	for _, p := range projectEnvs {
		var projectD = projectEnv.ResponseDto{
			Id:        p.ID,
			Name:      p.Name,
			ProjectId: p.ProjectId,
		}
		projectEnvDtos = append(projectEnvDtos, projectD)
	}
	return projectEnvDtos
}
