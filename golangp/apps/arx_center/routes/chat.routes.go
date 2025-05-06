package routes

import (
	"github.com/arxtect/ArxBackend/golangp/apps/arx_center/controllers"
	"github.com/arxtect/ArxBackend/golangp/common/middleware"

	"github.com/gin-gonic/gin"
	"github.com/toheart/functrace"
)

type ChatRouteController struct {
	chatController controllers.ChatController
}

func NewChatRouteController(chatController controllers.ChatController) ChatRouteController {
	defer functrace.Trace([]interface {
	}{chatController})()
	return ChatRouteController{chatController}
}

func (crc *ChatRouteController) ChatRoute(rg *gin.RouterGroup) {
	defer functrace.Trace([]interface {
	}{crc, rg})()

	chat := rg.Group("chat").Use(middleware.DeserializeUser())
	{
		chat.POST("/completion_with_model_info", crc.chatController.CompletionWithModelInfo)
	}
}
