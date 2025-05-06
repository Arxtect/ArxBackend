package service

import (
	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/models"
	"github.com/Arxtect/ArxBackend/golangp/common/initializers"

	"github.com/toheart/functrace"
)

func GetDocumentByFileHash(hash string) (*models.Document, error) {
	defer functrace.Trace([]interface {
	}{hash})()
	var doc models.Document
	err := initializers.DB.Model(&doc).Where("file_hash = ?", hash).First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func GetTagsByName(names []string) ([]models.Tag, error) {
	defer functrace.Trace([]interface {
	}{names})()
	var tags []models.Tag
	err := initializers.DB.Model(&models.Tag{}).Where("name IN (?)", names).Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func UpdateDocumentTitleAndTags(doc *models.Document, title string, tags []models.Tag, content string) error {
	defer functrace.Trace([]interface {
	}{doc, title, tags, content})()
	update := models.Document{}
	if title != "" {
		update.Title = title
	}
	if tags != nil && len(tags) > 0 {
		update.Tags = tags
	}
	if content != "" {
		update.Content = content
	}

	return initializers.DB.Model(doc).Omit("User").Updates(update).Error
}

func GetDocumentByKey(key string) (*models.Document, error) {
	defer functrace.Trace([]interface {
	}{key})()
	var doc models.Document
	err := initializers.DB.Model(&doc).
		Preload("User").Preload("Tags").
		Where("storage_key = ?", key).First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}
