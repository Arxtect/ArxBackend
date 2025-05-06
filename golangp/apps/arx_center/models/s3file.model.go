package models

import "gorm.io/gorm"

type S3File struct {
	gorm.Model
	Path     string     `gorm:"type:text;not null" json:"path"`
	Hash     string     `gorm:"type:text;not null" json:"hash"`
	Projects []*Project `gorm:"many2many:project_files;"`
}
