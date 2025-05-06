package routes

import (
	"github.com/arxtect/ArxBackend/golangp/apps/arx_center/controllers"
	"github.com/arxtect/ArxBackend/golangp/common/middleware"

	"github.com/gin-gonic/gin"
	"github.com/toheart/functrace"
)

type PostRouteController struct {
	postController controllers.PostController
}

func NewRoutePostController(postController controllers.PostController) PostRouteController {
	defer functrace.Trace([]interface {
	}{postController})()
	return PostRouteController{postController}
}

func (pc *PostRouteController) PostRoute(rg *gin.RouterGroup) {
	defer functrace.Trace([]interface {
	}{pc, rg})()

	router := rg.Group("posts").Use(middleware.DeserializeUser())
	router.POST("", pc.postController.CreatePost)
	router.GET("", pc.postController.FindPosts)
	router.PUT("/:postId", pc.postController.UpdatePost)
	router.GET("/getLatestPost", pc.postController.FindLatestPost)
	router.DELETE("/:postId", pc.postController.DeletePost)
}
