package models

import (
	"github.com/google/uuid"
	"github.com/toheart/functrace"
)

type Document struct {
	Base
	Title         string    `gorm:"type:varchar(255);not null" json:"title"`
	Content       string    `gorm:"type:text" json:"content"`
	AuthorID      uuid.UUID `gorm:"type:uuid;not null" json:"author_id"`
	StorageBucket string    `gorm:"type:varchar(100);not null" json:"storage_bucket"`
	StorageKey    string    `gorm:"type:varchar(100);not null" json:"storage_key"`
	StorageSize   int64     `gorm:"type:bigint" json:"storage_size"`
	FileHash      string    `gorm:"type:varchar(100);unique;not null" json:"file_hash"`
	StorageZip    string    `gorm:"type:varchar(100)" json:"storage_zip"`
	Tags          []Tag     `gorm:"many2many:document_tags;" json:"tags"`
	Cover         string    `gorm:"type:varchar(100)" json:"cover"`
	User          SafeUser  `gorm:"foreignKey:AuthorID" json:"user"`
}

type Tag struct {
	Base
	Name      string     `gorm:"type:varchar(100);unique;not null" json:"name"`
	Documents []Document `gorm:"many2many:document_tags;" json:"documents"`
}

func (s *Document) TableName() string {
	defer functrace.Trace([]interface {
	}{s})()
	return "documents"
}

type DocumentMigrate struct {
	Base
	Title         string    `gorm:"type:varchar(255);not null" json:"title"`
	Content       string    `gorm:"type:text" json:"content"`
	AuthorID      uuid.UUID `gorm:"type:uuid;not null" json:"author_id"`
	StorageBucket string    `gorm:"type:varchar(100);not null" json:"storage_bucket"`
	StorageKey    string    `gorm:"type:varchar(100);not null" json:"storage_key"`
	StorageSize   int64     `gorm:"type:bigint" json:"storage_size"`
	FileHash      string    `gorm:"type:varchar(100);unique;not null" json:"file_hash"`
	StorageZip    string    `gorm:"type:varchar(100)" json:"storage_zip"`
	Tags          []Tag     `gorm:"many2many:document_tags;" json:"tags"`
	Cover         string    `gorm:"type:varchar(100)" json:"cover"`
}

func (s *DocumentMigrate) TableName() string {
	defer functrace.Trace([]interface {
	}{s})()
	return "documents"
}
