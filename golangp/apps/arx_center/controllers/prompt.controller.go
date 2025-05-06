package controllers

import (
	"github.com/arxtect/ArxBackend/golangp/apps/arx_center/models"
	"github.com/arxtect/ArxBackend/golangp/apps/arx_center/service/dto"
	"github.com/arxtect/ArxBackend/golangp/common/constants"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/toheart/functrace"
	"gorm.io/gorm"
)

type PromptController struct {
	DB      *gorm.DB
	RedisDb *redis.Client
	Logger  *log.Logger
}

func NewPromptController(DB *gorm.DB, logger *log.Logger, redisDb *redis.Client) PromptController {
	defer functrace.Trace([]interface {
	}{DB, logger, redisDb})()
	return PromptController{
		DB:      DB,
		RedisDb: redisDb,
		Logger:  logger,
	}
}

func (pc *PromptController) GetPromptList(c *gin.Context) {
	defer functrace.Trace([]interface {
	}{pc, c})()
	req := dto.GetPromptRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid request"})
		return
	}

	var prompt []models.Prompt

	query := pc.DB.Model(&models.Prompt{})
	if req.Keyword != "" {

		query = query.Where("to_tsvector('simple', content) @@ plainto_tsquery('simple', ?)", req.Keyword)
	}

	var count int64
	err = query.Scopes(
		dto.Paginate(req.GetPageSize(), req.GetPageIndex()),
	).Order("created_at desc").
		Count(&count).
		Find(&prompt).Error
	if err != nil {
		pc.Logger.Printf("Error getting  prompt from DB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Error getting  prompt"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"prompt": prompt, "total": count, "pageIndex": req.GetPageIndex(), "pageSize": req.GetPageSize()}})
}

func (pc *PromptController) GetPrompt(c *gin.Context) {
	defer functrace.Trace([]interface {
	}{pc, c})()
	promptID := c.Param("id")
	if promptID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid request"})
		return
	}
	var prompt models.Prompt

	val, err := pc.RedisDb.Get(context.Background(), constants.RedisKeyAiPrompt+promptID).Result()
	if err == redis.Nil {

		err = pc.DB.First(&prompt, promptID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Prompt not found"})
			} else {
				pc.Logger.Printf("Error getting prompt from DB: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Error getting prompt"})
			}
			return
		}

		jsonPrompt, _ := json.Marshal(prompt)
		pc.RedisDb.Set(context.Background(), promptID, jsonPrompt, 0)
	} else if err != nil {
		pc.Logger.Printf("Error getting prompt from Redis: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Prompt not found"})
		return
	} else {

		json.Unmarshal([]byte(val), &prompt)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": prompt})
}

func (pc *PromptController) CreatePrompt(c *gin.Context) {
	defer functrace.Trace([]interface {
	}{pc, c})()
	var req dto.CreatePromptRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid request"})
		return
	}

	prompt := models.Prompt{
		Content:       req.Content,
		Keywords:      req.Keywords,
		Settings:      req.Settings,
		ReferenceFile: req.ReferenceFile,
	}

	err = pc.DB.Create(&prompt).Error
	if err != nil {
		pc.Logger.Printf("Error creating prompt in DB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Error creating prompt"})
		return
	}

	jsonPrompt, _ := json.Marshal(prompt)
	pc.RedisDb.Set(context.Background(), constants.RedisKeyAiPrompt+fmt.Sprintf("%v", prompt.ID), jsonPrompt, 0)

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": prompt})
}

func (pc *PromptController) UpdatePrompt(c *gin.Context) {
	defer functrace.Trace([]interface {
	}{pc, c})()
	var req dto.UpdatePromptRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid request"})
		return
	}
	if req.Base.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid request"})
		return
	}

	uuidID, err := uuid.Parse(req.Base.ID)
	if err != nil {
		pc.Logger.Printf("Error parsing UUID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid UUID"})
		return
	}

	var prompt models.Prompt
	result := pc.DB.First(&prompt, "id = ?", uuidID)
	if result.Error != nil {
		pc.Logger.Printf("Error finding prompt in DB: %v", result.Error)
		c.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Prompt not found"})
		return
	}

	prompt.Content = req.Content
	prompt.Keywords = req.Keywords
	prompt.Settings = req.Settings
	prompt.ReferenceFile = req.ReferenceFile

	result = pc.DB.Save(&prompt)
	if result.Error != nil {
		pc.Logger.Printf("Error updating prompt in DB: %v", result.Error)
		c.JSON(http.StatusOK, gin.H{"status": "fail", "message": "no prompt"})
		return
	}

	jsonPrompt, _ := json.Marshal(prompt)
	pc.RedisDb.Set(context.Background(), constants.RedisKeyAiPrompt+req.Base.ID, jsonPrompt, 0)

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Prompt updated"})
}

func (pc *PromptController) DeletePrompt(c *gin.Context) {
	defer functrace.Trace([]interface {
	}{pc, c})()
	promptID := c.Param("id")
	if promptID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid request"})
		return
	}

	uuidID, err := uuid.Parse(promptID)
	if err != nil {
		pc.Logger.Printf("Error parsing UUID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid UUID"})
		return
	}

	result := pc.DB.Delete(&models.Prompt{}, "id = ?", uuidID)
	if result.Error != nil {
		pc.Logger.Printf("Error deleting prompt from DB: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Error deleting prompt"})
		return
	}

	pc.RedisDb.Del(context.Background(), constants.RedisKeyAiPrompt+promptID)

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Prompt deleted"})
}
