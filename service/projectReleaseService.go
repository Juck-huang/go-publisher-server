package service

import (
	"hy.juck.com/go-publisher-server/dto/projectRelease"
	"hy.juck.com/go-publisher-server/model"
)

type ProjectReleaseService struct {
	ProjectId string
}

func NewProjectReleaseService(projectId string) *ProjectReleaseService {
	return &ProjectReleaseService{
		ProjectId: projectId,
	}
}

func (o *ProjectReleaseService) GetProjectRelease(projectEnvId string, projectTypeId string) (projectReleaseDto projectRelease.ResponseDto) {
	var projectRelease = model.ProjectRelease{}
	G.DB.Debug().Where("project_id = ? and project_env_id = ? and project_type_id = ?", o.ProjectId, projectEnvId, projectTypeId).First(&projectRelease)
	projectReleaseDto.Id = projectRelease.ID
	projectReleaseDto.ProjectId = projectRelease.ProjectId
	projectReleaseDto.ProjectEnvId = projectRelease.ProjectEnvId
	projectReleaseDto.ProjectTypeId = projectRelease.ProjectTypeId
	projectReleaseDto.BuildScriptPath = projectRelease.BuildScriptPath
	projectReleaseDto.Params = projectRelease.Params
	return projectReleaseDto
}
