package service

import (
	"hy.juck.com/go-publisher-server/dto/project"
	"hy.juck.com/go-publisher-server/model"
)

type ProjectService struct {
}

func NewProjectService() *ProjectService {
	return &ProjectService{}
}

// GetProjectList 获取项目列表
func (o *ProjectService) GetProjectList() (projectDtos []project.ResponseDto) {
	var projects []model.Project
	var projectEnvs []model.ProjectEnv
	G.DB.Debug().Find(&projects)
	G.DB.Debug().Find(&projectEnvs)
	for _, p := range projects {
		for _, projectEnv := range projectEnvs {
			if p.ID == projectEnv.ProjectId {
				var projectD = project.ResponseDto{
					Id:           p.ID,
					CreateDate:   p.CreatedAt.Format("2006-01-02 15:04:05"),
					UpdateDate:   p.UpdatedAt.Format("2006-01-02 15:04:05"),
					Name:         p.Name + "(" + projectEnv.Name + ")",
					ProjectEnvId: projectEnv.ID,
				}
				projectDtos = append(projectDtos, projectD)
			}
		}
	}
	return projectDtos
}
