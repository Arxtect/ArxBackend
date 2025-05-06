package routes

import (
	"github.com/arxtect/ArxBackend/golangp/apps/arx_center/controllers"

	"github.com/gin-gonic/gin"
)

type S3FileRouteController struct {
	S3FileController *controllers.S3FileController
}

func NewS3FileRouteController(S3FileController controllers.S3FileController) S3FileRouteController {
	return S3FileRouteController{&S3FileController}
}

func (dc *S3FileRouteController) S3FileRoute(rg *gin.RouterGroup) {
	// .Use(middleware.DeserializeUser())
	S3File := rg.Group("s3files")
	{
		S3File.POST("", dc.S3FileController.UploadFile)
		S3File.POST("presigned-urls", dc.S3FileController.GeneratePresignedDownloadURLs)
	}

}
