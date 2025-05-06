package routes

import (
	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/controllers"
	"github.com/Arxtect/ArxBackend/golangp/common/middleware"

	"github.com/gin-gonic/gin"
)

type ProjectRouteController struct {
	projectController controllers.ProjectController
}

func NewProjectRouteController(projectController controllers.ProjectController) ProjectRouteController {
	return ProjectRouteController{projectController}
}

func (prc *ProjectRouteController) ProjectkRoute(rg *gin.RouterGroup) {
	router := rg.Group("projects").Use(middleware.DeserializeUser())
	router.GET("", prc.projectController.GetProjects)
	router.POST("", prc.projectController.CreateProject)
	router.GET("/:id", prc.projectController.GetProjectViaID)
	router.PUT("/:id", prc.projectController.UpdateProject)
	router.DELETE("/:id", prc.projectController.DeleteProject)

	router.GET("/:id/access", prc.projectController.GetProjectAccess)
	router.POST("/:id/share-tokens", prc.projectController.GenerateProjectShareToken)
	router.GET("/:id/share-tokens/:token_id", prc.projectController.GetTokenInfo)
	router.POST("/:id/members", prc.projectController.AddProjectMember)
}
