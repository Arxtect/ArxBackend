package models

import (
	"errors"

	"github.com/Arxtect/ArxBackend/golangp/common/constants"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectMember struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	ProjectID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Role      string
}

var ProjectMemberRoles = []string{constants.ProjectRoleCollaborator, constants.ProjectRoleViewer}

func (sc *ProjectMember) BeforeSave(tx *gorm.DB) error {
	for _, validRole := range ProjectMemberRoles {
		if sc.Role == validRole {
			return nil
		}
	}

	roleOptions := ""
	for _, role := range ProjectMemberRoles {
		roleOptions += role + ", "
	}

	return errors.New("invalid role: " + sc.Role + ", valid roles are: " + roleOptions[:len(roleOptions)-2])
}
