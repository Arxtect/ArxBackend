package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/arxtect/ArxBackend/golangp/apps/arx_center/models"
	"github.com/arxtect/ArxBackend/golangp/common/constants"
	"github.com/arxtect/ArxBackend/golangp/common/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectController struct {
	DB *gorm.DB
}

func NewProjectController(DB *gorm.DB) ProjectController {
	return ProjectController{
		DB,
	}
}

func (controller *ProjectController) getOwnedProjects(userId uuid.UUID) ([]models.Project, error) {
	var projects []models.Project
	if err := controller.DB.Where("owner_id = ?", userId).Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (controller *ProjectController) getSharedProjects(userId uuid.UUID) ([]models.Project, error) {
	var projects []models.Project
	err := controller.DB.Joins("JOIN project_members ON projects.id = project_members.project_id").
		Where("project_members.user_id = ?", userId).
		Find(&projects).Error
	return projects, err
}

func (controller *ProjectController) GetProjects(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	relation := c.Query("relation") // Get relation query parameter

	// Validate relation parameter
	if relation != "" && relation != "owner" && relation != "member" {
		logger.Warning("Invalid relation parameter: %s", relation)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid relation parameter. Must be 'owner', or 'member'"})
		return
	}

	projects := make([]models.Project, 0)
	if relation == "" || relation == "owner" {
		ownedProjects, err := controller.getOwnedProjects(currentUser.ID)
		if err != nil {
			logger.Warning("Failed to retrieve owned projects: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve owned projects"})
			return
		}
		projects = append(projects, ownedProjects...)
	}

	if relation == "" || relation == "member" {
		sharedProjects, err := controller.getSharedProjects(currentUser.ID)
		if err != nil {
			logger.Warning("Failed to retrieve shared projects: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve shared projects"})
			return
		}
		projects = append(projects, sharedProjects...)
	}

	c.JSON(http.StatusOK, projects)
}

// GetProjectViaID retrieves a project by its ID
func (controller *ProjectController) GetProjectViaID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		logger.Warning("Invalid project ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid project ID"})
		return
	}

	var projectManager models.Project
	project, err := projectManager.GetProjectViaID(controller.DB, id)
	if err != nil {
		logger.Warning("Project not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "Project not found"})
		return
	}

	logger.Info("Retrieved project with ID: %s", id)
	c.JSON(http.StatusOK, project)
}

// CreateProject creates a new project
func (controller *ProjectController) CreateProject(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)

	fmt.Printf("Current user: %+v\n", currentUser)

	createPayload := models.CreateProjectPayload{}

	if err := c.ShouldBindJSON(&createPayload); err != nil {
		logger.Warning("Invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	projectManager := models.Project{
		ID:          uuid.New(),
		OwnerID:     currentUser.ID,
		ProjectName: createPayload.ProjectName,
	}

	if err := projectManager.CreateProject(controller.DB); err != nil {
		logger.Warning("Failed to create project: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create project"})
		return
	}

	logger.Info("Created new project with ID: %s", projectManager.ID)
	c.JSON(http.StatusCreated, projectManager)
}

// UpdateProject updates an existing project
func (controller *ProjectController) UpdateProject(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		logger.Warning("Invalid project ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid project ID"})
		return
	}

	var projectManager models.Project
	if err := controller.DB.First(&projectManager, "id = ?", id).Error; err != nil {
		logger.Warning("Project not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "Project not found"})
		return
	}

	var updatePayload models.UpdateProjectPayload

	if err := c.ShouldBindJSON(&updatePayload); err != nil {
		logger.Warning("Invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Update only if the field is not a zero value
	if updatePayload.ProjectName != "" {
		projectManager.ProjectName = updatePayload.ProjectName
	}
	if updatePayload.ShareLink != "" {
		projectManager.ShareLink = updatePayload.ShareLink
	}
	if updatePayload.Email != "" {
		projectManager.Email = updatePayload.Email
	}
	if updatePayload.RoomName != "" {
		projectManager.RoomName = updatePayload.RoomName
	}
	if updatePayload.IsSync != projectManager.IsSync {
		projectManager.IsSync = updatePayload.IsSync
	}

	if err := projectManager.UpdateProject(controller.DB); err != nil {
		logger.Warning("Failed to update project: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update project"})
		return
	}

	logger.Info("Updated project with ID: %s", id)
	c.JSON(http.StatusOK, projectManager)
}

// DeleteProject deletes a project by its ID
func (controller *ProjectController) DeleteProject(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		logger.Warning("Invalid project ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid project ID"})
		return
	}

	var projectManager models.Project
	if err := controller.DB.First(&projectManager, "id = ?", id).Error; err != nil {
		logger.Warning("Project not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "Project not found"})
		return
	}

	if err := projectManager.DeleteProject(controller.DB); err != nil {
		logger.Warning("Failed to delete project: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete project"})
		return
	}

	logger.Info("Deleted project with ID: %s", id)
	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

func (controller *ProjectController) validateProjectId(c *gin.Context, idString string, writable bool) (*models.Project, error) {
	id, err := uuid.Parse(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid project ID"})
		return nil, err
	}

	var project models.Project
	if err := controller.DB.First(&project, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Project not found"})
		return nil, err
	}

	if writable && project.OwnerID != c.MustGet("currentUser").(models.User).ID {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have permission to access this project"})
		return nil, fmt.Errorf("permission denied")
	}

	return &project, nil
}

type GenerateProjectShareTokenPayload struct {
	Validity string `json:"validity" binding:"required"`
}

func (controller *ProjectController) GenerateProjectShareToken(c *gin.Context) {
	payload := GenerateProjectShareTokenPayload{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	project, err := controller.validateProjectId(c, c.Param("id"), true)
	if err != nil {
		return
	}

	// Parse validity period as duration (e.g., "24h", "7d")
	duration, err := time.ParseDuration(payload.Validity)
	if err != nil {
		fmt.Println("Error parsing duration:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid validity period format. Use duration like '24h' or '7d'"})
		return
	}

	// Generate new tokens and expiration times
	collaboratorToken := uuid.New()
	viewerToken := uuid.New()
	expireAt := time.Now().Add(duration)

	// Update project with new tokens and expiration times
	project.CollaboratorShareToken = &collaboratorToken
	project.ViewerShareToken = &viewerToken
	project.ShareTokenExpireAt = &expireAt

	// Save to database
	if err := controller.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate share tokens"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{})
}

func (controller *ProjectController) GetTokenInfo(c *gin.Context) {
	project, err := controller.validateProjectId(c, c.Param("id"), false)
	if err != nil {
		return
	}

	var tokenType string
	if c.Param("token_id") == project.CollaboratorShareToken.String() {
		tokenType = "collaborator"
	} else if c.Param("token_id") == project.ViewerShareToken.String() {
		tokenType = "viewer"
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid token ID"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"tokenType":   tokenType,
		"expireAt":    project.ShareTokenExpireAt,
		"projectId":   project.ID,
		"projectName": project.ProjectName,
	})
}

type AddProjectMembersPayload struct {
	Token string `json:"token" binding:"required"`
}

func (controller *ProjectController) AddProjectMember(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)

	// Validate project ID
	project, err := controller.validateProjectId(c, c.Param("id"), false)
	if err != nil {
		return
	}

	// Bind JSON payload
	var payload AddProjectMembersPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Validate role and token
	TokenMap := map[string]string{
		constants.ProjectRoleCollaborator: project.CollaboratorShareToken.String(),
		constants.ProjectRoleViewer:       project.ViewerShareToken.String(),
	}

	var validRole string
	for role, token := range TokenMap {
		if payload.Token == token {
			validRole = role
			break
		}
	}

	if validRole == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid token"})
		return
	}

	// Check if user is owner
	fmt.Println("Current user ID:", currentUser.ID)
	fmt.Println("Project owner ID:", project.OwnerID)
	if currentUser.ID == project.OwnerID {

		c.JSON(http.StatusBadRequest, gin.H{"message": "Owner cannot be added as a member"})
		return
	}

	// Add or update project member
	member := models.ProjectMember{
		UserID:    currentUser.ID,
		ProjectID: project.ID,
		Role:      validRole,
	}
	if err := controller.DB.Save(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save project member"})
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{"message": "Member saved", "role": member.Role})
}

func (controller *ProjectController) GetProjectAccess(c *gin.Context) {
	user := c.MustGet("currentUser").(models.User)
	project, err := controller.validateProjectId(c, c.Param("id"), false)
	if err != nil {
		return
	}

	body := gin.H{
		"user": gin.H{
			"id":   user.ID,
			"name": user.Name,
		},
		"readOnly": true,
	}

	if user.Role == constants.AppRoleAdmin {
		body["readOnly"] = false
		c.JSON(http.StatusOK, body)
		return
	}

	// Check if the user is the owner
	if project.OwnerID == user.ID {
		body["readOnly"] = false
		c.JSON(http.StatusOK, body)
		return
	}

	// Check if the user is a member
	var member models.ProjectMember
	if err := controller.DB.Where("user_id = ? AND project_id = ?", c.MustGet("currentUser").(models.User).ID, project.ID).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "You do not have access to this project",
		})
		return
	}

	body["readOnly"] = member.Role != constants.ProjectRoleCollaborator
	c.JSON(http.StatusOK, body)
}
