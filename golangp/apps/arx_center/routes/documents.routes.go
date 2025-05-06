package routes

import (
	"github.com/arxtect/ArxBackend/golangp/apps/arx_center/controllers"
	"github.com/arxtect/ArxBackend/golangp/common/middleware"

	"github.com/gin-gonic/gin"
	"github.com/toheart/functrace"
)

type DocumentsController struct {
	documentController controllers.DocumentsController
}

func NewDocumentRouteController(documentController controllers.DocumentsController) DocumentsController {
	defer functrace.Trace([]interface {
	}{documentController})()
	return DocumentsController{documentController}
}

func (dc *DocumentsController) DocumentRoute(rg *gin.RouterGroup) {
	defer functrace.Trace([]interface {
	}{dc, rg})()
	documents := rg.Group("documents").Use(middleware.DeserializeUser())
	{
		documents.POST("/upload", dc.documentController.UploadDocumentsByUser, middleware.Sentinel())
		documents.GET("/drafts", dc.documentController.GetDraftsByUser)
		documents.GET("/drafts/:key", dc.documentController.GetDocumentByKey)
		documents.POST("/gen/commitInfo", dc.documentController.GenCommitDocument)

		documents.POST("/generateS3UploadURL", dc.documentController.GenerateS3UploadURL)

	}

	documentsNoauth := rg.Group("documents")
	{
		documentsNoauth.GET("/tags/list", dc.documentController.GetDocumentsAllTags)
		documentsNoauth.GET("/list/search", dc.documentController.GetDocumentsListSearch)
		documentsNoauth.GET("/list/search-v2", dc.documentController.GetDocumentsListSearchV2)

		documentsNoauth.GET("/pre/download/:key", dc.documentController.GetDocumentDownloadUrl)
		documentsNoauth.GET("/pre/preview/:key", dc.documentController.PreViewFile)
	}
}
