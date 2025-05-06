package controllers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/dify"
	"github.com/Arxtect/ArxBackend/golangp/common/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ConversationsController struct {
	DB *gorm.DB
}

type AppList struct {
	ID            string                 `json:"id,omitempty"`
	Name          string                 `json:"name,omitempty"`
	AccessToken   string                 `json:"access_token,omitempty"`
	Default       bool                   `json:"default,omitempty"`
	FileUpload    dify.FileUploadPayload `json:"file_upload,omitempty"`
	Selection     bool                   `json:"selection,omitempty"`
	UserInputForm []any                  `json:"user_input_form,omitempty"`
	UploadPrompt  string                 `json:"upload_prompt,omitempty"`
}

func InitDifyClient() *dify.DifyClient {
	client, err := dify.GetDifyClient()
	if err != nil {
		logger.Warning(err.Error())
		return nil // 返回 nil 以符合返回类型
	}

	_, err = client.GetUserToken()
	if err != nil {
		logger.Warning(err.Error())
		return nil // 返回 nil 以符合返回类型
	}
	return client
}

func NewConversationsController(DB *gorm.DB) ConversationsController {
	return ConversationsController{DB}
}

func GetAppAuthorization(difyClient *dify.DifyClient, ctx *gin.Context) (string, error) {
	accessToken := ctx.GetHeader("APP-Authorization")
	if accessToken != "" {
		return accessToken, nil
	}

	result, err := difyClient.GetAppsAccessToken("7gINmdCL3hhqyRsq")

	if err != nil {
		return "", err
	}

	return "Bearer " + result.AccessToken, nil
}

func (c *ConversationsController) GetAppList(ctx *gin.Context) {
	list := []AppList{
		{
			ID:      "9a1686f4-e641-474a-96c5-c6812908e046",
			Name:    "Ask ai anything",
			Default: true,
		},
		{
			ID:           "22930441-0b94-4252-9db9-c05256550002",
			Name:         "Continue writing",
			Default:      false,
			Selection:    true,
			UploadPrompt: "Help me continue writing",
		},
		{
			ID:           "23f33933-67f5-4a8c-9f94-be3265c7ffd9",
			Name:         "Add an equation",
			Default:      false,
			UploadPrompt: "Generate laTex equation",
		},
		{
			ID:           "3fe91c2b-0924-4080-82fb-47c4f6ee8929",
			Name:         "Add a table",
			Default:      false,
			UploadPrompt: "Generate laTex table",
		},
		//{
		//	ID:        "cc3378c3-dd99-4c9a-8a12-8f81e2e71acd",
		//	Name:      "Section Polisher",
		//	Default:   false,
		//	Selection: true,
		//},
	}

	difyClient := InitDifyClient()

	for i := 0; i < len(list); i++ {
		app, _ := difyClient.GetAppsHandler(list[i].ID)
		list[i].FileUpload = app.ModelConfig.FileUpload
		list[i].UserInputForm = app.ModelConfig.UserInputForm
		list[i].AccessToken = app.Site.AccessToken
	}

	ctx.JSON(http.StatusOK, list)
}

func (c *ConversationsController) GetConversationsAccessToken(ctx *gin.Context) {
	accessToken := ctx.Query("access_token")

	difyClient := InitDifyClient()

	result, err := difyClient.GetAppsAccessToken(accessToken)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": result.AccessToken})
}

func (c *ConversationsController) AutoComplete(ctx *gin.Context) {
	// 解析请求体
	payload := dify.AutoCompletePayload{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// 初始化 Dify 客户端
	difyClient := InitDifyClient()
	if difyClient == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to initialize Dify client"})
		return
	}

	// 获取 APP-Authorization 令牌
	payload.AcessToken = "Bearer " + dify.API_AUTOCOMPLETE_ACCESSTOKEN

	// 调用 Dify 的 AutoComplete API
	result, err := difyClient.AutoComplete(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// 返回AI建议结果
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": result.Data.Outputs.Suggestion})
}

func (c *ConversationsController) ChatMessages(ctx *gin.Context) {
	payload := dify.ChatMessagesPayload{}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	difyClient := InitDifyClient()

	appAuthorization, err := GetAppAuthorization(difyClient, ctx)

	if appAuthorization == "" || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "missing APP-Authorization header"})
		return
	}

	if payload.ResponseMode == "blocking" {
		result, err := difyClient.ChatMessages(payload.Query, payload.Inputs, payload.ConversationID, payload.Files, appAuthorization)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": result})
	} else {
		// 设置流式响应头
		ctx.Writer.Header().Set("Content-Type", "text/event-stream")
		ctx.Writer.Header().Set("Connection", "keep-alive")
		ctx.Writer.Header().Set("Cache-Control", "no-cache")
		ctx.Writer.Header().Set("X-Accel-Buffering", "no") // 关闭缓冲
		ctx.Writer.Header().Set("Content-Type", "application/octet-stream")

		dataChan, err := difyClient.ChatMessagesStreaming(payload.Query, payload.Inputs, payload.ConversationID, payload.Files, appAuthorization)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
			return
		}

		// 使用 ctx.Stream 进行流式传输
		ctx.Stream(func(w io.Writer) bool {
			if msg, ok := <-dataChan; ok {
				_, err := w.Write([]byte(msg + "\n\n"))
				if err != nil {
					fmt.Printf("Error writing message: %s\n", err)
					return false
				}
				//fmt.Printf("Message sent: %s\n", msg)
				ctx.Writer.Flush() // 立即刷新输出缓冲区
				return true
			}
			return false
		})
	}
}

func (c *ConversationsController) ChatMessagesStop(ctx *gin.Context) {
	task_id := ctx.Param("task_id") // 从URL参数获取 task_id
	if task_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "missing task_id parameter"})
		return
	}
	difyClient := InitDifyClient()
	AppAuthorization, err := GetAppAuthorization(difyClient, ctx)
	if AppAuthorization == "" || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "missing APP-Authorization header"})
		return
	}

	_, err = difyClient.ChatMessagesStop(task_id, AppAuthorization)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (c *ConversationsController) HandleFileUploadForChat(ctx *gin.Context) {
	// 获取文件
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to get file from request"})
		return
	}
	defer file.Close()

	difyClient := InitDifyClient()

	appAuthorization, err := GetAppAuthorization(difyClient, ctx)

	if appAuthorization == "" || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "missing APP-Authorization header"})
		return
	}

	// 上传文件
	result, err := difyClient.DatasetsFileUploadChat(file, header.Filename, appAuthorization)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Failed to upload file: %v", err)})
		return
	}

	// Debugging: Print the response
	fmt.Println("Response:")
	fmt.Println(result)

	ctx.JSON(http.StatusOK, result)
}

func (c *ConversationsController) PreviewFile(ctx *gin.Context) {
	fileID := ctx.Param("file_id")

	difyClient := InitDifyClient()

	body, contentType, err := difyClient.PreviewFile(fileID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.Header("Content-Type", contentType)
	ctx.Status(http.StatusOK)
	ctx.Writer.Write(body)
}
