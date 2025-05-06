package routes

import (
	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/controllers"
	"github.com/Arxtect/ArxBackend/golangp/common/middleware"

	"github.com/gin-gonic/gin"
	"github.com/toheart/functrace"
)

type YRedisRouteController struct {
	yredisController controllers.YRedisController
}

func NewYRedisRouteController(yredisController controllers.YRedisController) YRedisRouteController {
	defer functrace.Trace([]interface {
	}{yredisController})()
	return YRedisRouteController{yredisController}
}

func (yrc *YRedisRouteController) YRediskRoute(sv *gin.RouterGroup) {
	defer functrace.Trace([]interface {
	}{yrc, sv})()
	rootrg := sv.Group("/")
	{
		rootrg.GET("/yredis/auth/token/:room", middleware.DeserializeUser(), yrc.yredisController.YRedisAuthToken)

		rootrg.GET("/yredis/auth/perm/:room/:userid", yrc.yredisController.YRedisRoomPermissionCallback)
		rootrg.PUT("/yredis/ydoc/:room", yrc.yredisController.YRedisYDocUpdateCallback)

		rootrg.GET("/yredis/room/share/user/:room", middleware.DeserializeUser(), yrc.yredisController.YRedisRoomShareUserGet)
		rootrg.PUT("/yredis/room/share/user", middleware.DeserializeUser(), yrc.yredisController.YRedisRoomShareUserUpdate)

		rootrg.POST("/yredis/room/create", middleware.DeserializeUser(), yrc.yredisController.YRedisRoomCreateRoom)

		rootrg.DELETE("/yredis/room/share/user", middleware.DeserializeUser(), yrc.yredisController.YRedisRoomShareUserRemove)

		rootrg.GET("/yredis/room/share/:room", middleware.DeserializeUser(), yrc.yredisController.YRedisRoomShareGet)
		rootrg.DELETE("/yredis/room/share", middleware.DeserializeUser(), yrc.yredisController.YRedisRoomShareClose)
		rootrg.POST("/yredis/room/share/reopen", middleware.DeserializeUser(), yrc.yredisController.YRedisRoomShareReopen)

	}
}
