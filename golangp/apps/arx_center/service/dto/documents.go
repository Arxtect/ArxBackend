package dto

import (
	"github.com/google/uuid"
)

type FileResponse struct {
	FileStorageID     string `json:"file_storage_id"`
	FileName          string `json:"file_name"`
	FileSize          int64  `json:"file_size"`
	FileStorageBucket string `json:"file_storage_bucket"`
}

type AddTagToDocumentRequest struct {
	Title      string   `json:"title" binding:"required"`
	Tags       []string `json:"tags"`
	Content    string   `json:"content"`
	UploadType string   `json:"upload_type"`
}

type GetDocumentsByUserRequest struct {
	Pagination `search:"-"`
	Keyword    string   `form:"keyword"`
	Tags       []string `form:"tags"`
}

type GetDocumentDownloadUrlRequest struct {
	FileStorageID string `json:"file_storage_id" binding:"required"`
}

type SearchListDocumentsResp struct {
	ID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name  string    `gorm:"type:varchar(255);not null" json:"name"`
	Email string    `gorm:"uniqueIndex;not null" json:"email"`
}

type GenCommitDocumentRequest struct {
	SummaryDiff string `json:"summary_diff" binding:"required"`
	Language    string `json:"language"`
}

type GenerateS3UploadURLRequest struct {
	Filename string `json:"filename" binding:"required"`
}
