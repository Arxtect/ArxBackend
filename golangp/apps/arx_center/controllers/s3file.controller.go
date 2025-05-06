package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/arxtect/ArxBackend/golangp/apps/arx_center/models"
	"github.com/arxtect/ArxBackend/golangp/common/xminio"
	"github.com/arxtect/ArxBackend/golangp/config"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"

	"gorm.io/gorm"
)

type S3FileController struct {
	db *gorm.DB
	s3 *xminio.S3Manager
}

func NewS3FileController(db *gorm.DB) *S3FileController {
	s3 := xminio.NewS3Manager(config.Env.MinioBucket, config.Env.MinioAccessKey, config.Env.MinioSecretKey, config.Env.MinioBucketUrl)
	return &S3FileController{
		db: db,
		s3: s3,
	}
}

func (fc *S3FileController) UploadFile(ctx *gin.Context) {
	// Parse form data
	projectId := ctx.PostForm("projectId")
	hash := ctx.PostForm("hash")
	if projectId == "" || hash == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "projectId and hash are required"})
		return
	}

	// Check if file with this hash already exists (early check)
	var existingFile models.S3File
	if err := fc.db.Where("hash = ?", hash).First(&existingFile).Error; err == nil {
		// File with this hash already exists, return immediately
		ctx.JSON(http.StatusOK, gin.H{
			"message": "File with this hash already exists",
			"path":    existingFile.Path,
		})
		return
	} else if err != gorm.ErrRecordNotFound {
		// Database error
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	// Get file from form
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "File is required"})
		return
	}
	defer file.Close()

	// Calculate actual SHA256 hash of the uploaded file
	calculatedHash, err := calculateFileHash(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to calculate file hash"})
		return
	}

	// Verify provided hash matches calculated hash
	if hash != calculatedHash {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message":         "Provided hash doesn't match file content",
			"calculated_hash": hash,
		})
		return
	}

	// Reset file pointer to beginning for upload
	_, err = file.Seek(0, 0)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to process file"})
		return
	}

	// Generate S3 path
	s3Path := fmt.Sprintf("project-assets/%s", hash)

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream" // Default content type if not provided
	}

	putOpts := minio.PutObjectOptions{
		ContentType: contentType,
		UserMetadata: map[string]string{
			"hash":     hash,
			"filename": header.Filename,
		},
		// Optional: Add server-side encryption if needed
		// ServerSideEncryption: minio.NewSSE(),
	}

	_, err = fc.s3.Client.PutObject(
		ctx,
		fc.s3.BucketName,
		s3Path,
		file,
		header.Size,
		putOpts,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to upload file to S3"})
		return
	}

	// Save file metadata to database
	s3File := models.S3File{
		Path: s3Path,
		Hash: hash,
	}

	if err := fc.db.Create(&s3File).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save file metadata"})
	}

	// TODO: Link file to project
	// err = fc.db.Transaction(func(tx *gorm.DB) error {
	// 	// Save file metadata to database
	// 	if err := tx.Create(&s3File).Error; err != nil {
	// 		return err // Transaction will rollback automatically
	// 	}

	// 	// Link file to project
	// 	if err := tx.Exec(`
	// 		INSERT INTO project_files (project_manager_id, s3_file_id)
	// 		VALUES (?, ?)`, projectId, s3File.ID).Error; err != nil {
	// 		return err
	// 	}

	// 	return nil
	// })

	// if err != nil {
	// 	// Cleanup: remove file from S3 if DB operations fail
	// 	if deleteErr := fc.s3.Client.RemoveObject(ctx, fc.s3.BucketName, s3Path, minio.RemoveObjectOptions{}); deleteErr != nil {
	// 		// Log the deletion error but don't fail the request because of it
	// 		// You might want to add logging here
	// 		fmt.Printf("Failed to cleanup S3 file %s: %v\n", s3Path, deleteErr)
	// 	}

	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"detail": err.Error(),
	// 	})
	// 	return
	// }

	ctx.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"path":    s3Path,
		"hash":    hash,
	})
}

func calculateFileHash(file multipart.File) (string, error) {
	hash := sha256.New()
	_, err := io.Copy(hash, file)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

type PresignDownloadRequest struct {
	Hashes []string `json:"hashes" binding:"required"`
}

// PresignedURLResponse defines the response structure for each file
type PresignedURLResponse struct {
	Hash  string `json:"hash"`
	URL   string `json:"url,omitempty"`
	Error string `json:"error,omitempty"`
}

// GeneratePresignedDownloadURLs generates presigned URLs for multiple files based on their hashes
func (fc *S3FileController) GeneratePresignedDownloadURLs(ctx *gin.Context) {
	var req PresignDownloadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	if len(req.Hashes) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "At least one hash is required"})
		return
	}

	// Query database for files matching the provided hashes
	var files []models.S3File
	if err := fc.db.Where("hash IN ?", req.Hashes).Find(&files).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	// Create a map of hash to file path for easier lookup
	hashToPath := make(map[string]string)
	for _, file := range files {
		hashToPath[file.Hash] = file.Path
	}

	// Generate presigned URLs for each requested hash
	results := make([]PresignedURLResponse, len(req.Hashes))
	urlExpiry := 24 * time.Hour // URL valid for 24 hours, adjust as needed

	for i, hash := range req.Hashes {
		result := PresignedURLResponse{Hash: hash}

		// Check if file exists in database
		s3Path, exists := hashToPath[hash]
		if !exists {
			result.Error = "File not found"
			results[i] = result
			continue
		}

		// Generate presigned URL
		url, err := fc.s3.Client.PresignedGetObject(
			ctx,
			fc.s3.BucketName,
			s3Path,
			urlExpiry,
			nil, // Query parameters (optional)
		)
		if err != nil {
			result.Error = "Failed to generate presigned URL"
			results[i] = result
			continue
		}

		result.URL = url.String()
		results[i] = result
	}

	ctx.JSON(http.StatusOK, gin.H{
		"urls": results,
	})
}
