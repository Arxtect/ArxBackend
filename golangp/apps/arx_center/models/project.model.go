package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	ProjectName string    `gorm:"type:text;" json:"projectName"`

	OwnerID uuid.UUID `gorm:"type:uuid;" json:"ownerId"`
	Owner   *User     `gorm:"foreignKey:OwnerID" json:"owner"`
	Members []*User   `gorm:"many2many:project_members;" json:"members"`
	S3Files []*S3File `gorm:"many2many:project_files;"`

	// Share
	CollaboratorShareToken *uuid.UUID `gorm:"type:uuid;" json:"collaboratorShareToken"`
	ViewerShareToken       *uuid.UUID `gorm:"type:uuid;" json:"viewerShareToken"`
	ShareTokenExpireAt     *time.Time `json:"shareTokenExpireAt"`

	// Legacy Project
	CreateBy  uuid.UUID `gorm:"type:uuid;" json:"create_by"`
	ShareLink string    `gorm:"type:text;" json:"share_link"`
	RoomName  string    `gorm:"type:text;" json:"room_name"`
	Email     string    `gorm:"type:text;" json:"email"`
	IsSync    bool      `gorm:"default:false" json:"is_sync"`
}

type CreateProjectPayload struct {
	ProjectName string `json:"projectName" binding:"required"`
}

type UpdateProjectPayload struct {
	ProjectName string `json:"project_name"`
	ShareLink   string `json:"share_link"`
	Email       string `json:"email" `
	RoomName    string `json:"room_name" `
	IsSync      bool   `json:"is_sync"`
}

func (p *Project) GetProjectViaID(db *gorm.DB, id uuid.UUID) (*Project, error) {
	var project Project
	if err := db.Where("id = ?", id).First(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (p *Project) CreateProject(db *gorm.DB) error {
	return db.Create(p).Error
}

func (p *Project) UpdateProject(db *gorm.DB) error {
	return db.Save(p).Error
}

func (p *Project) DeleteProject(db *gorm.DB) error {
	return db.Delete(p).Error
}
