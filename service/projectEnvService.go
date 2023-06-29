package service

import (
	"hy.juck.com/go-publisher-server/dto/projectEnv"
	"hy.juck.com/go-publisher-server/model"
)

type ProjectEnvService struct {
}

func NewProjectEnvService() *ProjectEnvService {
	return &ProjectEnvService{}
}

// GetProjectEnvList 获取项目列表
func (o *ProjectEnvService) GetProjectEnvList(projectId int64) (projectEnvDtos []projectEnv.ResponseDto) {
	var projectEnvs []model.ProjectEnv
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
