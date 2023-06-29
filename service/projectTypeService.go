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

func (o *ProjectTypeService) GetByProjectId(projectId int64) []projectType.ResponseDto {
	var projectTypes []model.ProjectType
	var projectTypeDtos []projectType.ResponseDto
	G.DB.Debug().Where("project_id = ?", projectId).Find(&projectTypes)
	if len(projectTypes) > 0 {
		for _, pt := range projectTypes {
			p := projectType.ResponseDto{
				Id:        pt.ID,
				Name:      pt.Name,
				ProjectId: pt.ProjectId,
				TreeId:    pt.TreeId,
				ParentId:  pt.ParentId,
				TreeLevel: pt.TreeLevel,
				IsLeaf:    pt.IsLeaf,
			}
			projectTypeDtos = append(projectTypeDtos, p)
		}
	}
	return projectTypeDtos
}
