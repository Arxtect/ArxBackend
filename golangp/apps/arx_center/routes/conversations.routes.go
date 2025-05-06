package routes

import (
	"github.com/arxtect/ArxBackend/golangp/apps/arx_center/controllers"
	"github.com/arxtect/ArxBackend/golangp/common/middleware"

	"github.com/gin-gonic/gin"
	"github.com/toheart/functrace"
)

type ConversationsRouteController struct {
	ConversationsController controllers.ConversationsController
}

func NewConversationsRouteController(conversationsController controllers.ConversationsController) ConversationsRouteController {
	defer functrace.Trace([]interface {
	}{conversationsController})()
	return ConversationsRouteController{conversationsController}
}

func (dc *ConversationsRouteController) ConversationsRoute(rg *gin.RouterGroup) {
	defer functrace.Trace([]interface {
	}{dc, rg})()
	chatRouter := rg.Group("/chat")
	chatRouter.Use(middleware.DeserializeUser())
	chatRouter.GET("/access_token", dc.ConversationsController.GetConversationsAccessToken)
	chatRouter.GET("/app", dc.ConversationsController.GetAppList)
	chatRouter.POST("/chat-messages", dc.ConversationsController.ChatMessages)
	chatRouter.POST("/chat-messages/:task_id/stop", dc.ConversationsController.ChatMessagesStop)
	chatRouter.POST("/upload", dc.ConversationsController.HandleFileUploadForChat)
	chatRouter.POST("/auto-complete", dc.ConversationsController.AutoComplete)

	chatRouter.GET("/file/:file_id/preview", dc.ConversationsController.PreviewFile)
}
