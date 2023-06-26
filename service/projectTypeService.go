package service

import (
	"hy.juck.com/go-publisher-server/dto/projectType"
	"hy.juck.com/go-publisher-server/model"
)

type ProjectTypeService struct {
}

func NewProjectTypeService() *ProjectTypeService {
	return &ProjectTypeService{}
}

func (o *ProjectTypeService) GetProjectType(projectId string, projectTypeId string) (projectTypeDto projectType.ResponseDto) {
	var projectType = model.ProjectType{}
	G.DB.Debug().Where("id = ? and project_id = ?", projectTypeId, projectId).First(&projectType)
	projectTypeDto.Id = projectType.ID
	projectTypeDto.Name = projectType.Name
	projectTypeDto.ProjectId = projectType.ProjectId
	//projectTypeDto.BuildScriptPath = projectType.BuildScriptPath
	return projectTypeDto
}
