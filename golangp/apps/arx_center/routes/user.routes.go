package routes

import (
	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/controllers"
	"github.com/Arxtect/ArxBackend/golangp/common/middleware"

	"github.com/gin-gonic/gin"
	"github.com/toheart/functrace"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	defer functrace.Trace([]interface {
	}{userController})()
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {
	defer functrace.Trace([]interface {
	}{uc, rg})()

	router := rg.Group("users").Use(middleware.DeserializeUser())
	router.GET("/me", uc.userController.GetMe)
	router.POST("/admin_update_balance", uc.userController.AdminUpdateBalance)
	router.GET("/gitea/token", uc.userController.GetUserAccessTokens)
	router.POST("/gitea/repo", uc.userController.CreateGiteaRepo)
	router.GET("/gitea/repoList", uc.userController.GetUserRepoList)
	router.DELETE("/gitea/repo/:name", uc.userController.DeleteGiteaRepo)
	router.GET("/gitea/:token/validate", uc.userController.ValidateAccessToken)

	routerWs := rg.Group("ws").Use(middleware.DeserializeUser())
	{
		routerWs.POST("/establishWs", uc.userController.WsEditingRoom)
		routerWs.POST("/room/:fileId", uc.userController.CreateRoom)

	}

}
