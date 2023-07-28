package service

import (
	"hy.juck.com/go-publisher-server/dto/projectEnv"
	"hy.juck.com/go-publisher-server/model/project"
)

type ProjectEnvService struct {
}

func NewProjectEnvService() *ProjectEnvService {
	return &ProjectEnvService{}
}

// GetProjectEnvList 获取项目环境列表
func (o *ProjectEnvService) GetProjectEnvList() (projectEnvDtos []projectEnv.ResponseDto) {
	var projectEnvs []project.ProjectEnv
	G.DB.Debug().Find(&projectEnvs)
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
