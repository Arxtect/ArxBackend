package routes

import (
	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/controllers"

	"github.com/gin-gonic/gin"
	"github.com/toheart/functrace"
)

type PromptController struct {
	promptController controllers.PromptController
}

func NewPromptRouteController(promptController controllers.PromptController) PromptController {
	defer functrace.Trace([]interface {
	}{promptController})()
	return PromptController{promptController}
}

func (dc *PromptController) PromptRoute(rg *gin.RouterGroup) {
	defer functrace.Trace([]interface {
	}{dc, rg})()

	prompt := rg.Group("prompt")
	{
		prompt.GET("list", dc.promptController.GetPromptList)
		prompt.GET("/:id", dc.promptController.GetPrompt)
		prompt.POST("", dc.promptController.CreatePrompt)
		prompt.PUT("", dc.promptController.UpdatePrompt)
		prompt.DELETE("/:id", dc.promptController.DeletePrompt)
	}

}
