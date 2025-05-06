package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/gitea"
	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/models"
	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/service/ws"
	"github.com/Arxtect/ArxBackend/golangp/common/constants"
	"github.com/Arxtect/ArxBackend/golangp/common/logger"
	"github.com/Arxtect/ArxBackend/golangp/common/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
	cs *CreditSystem
}

func NewUserController(DB *gorm.DB, creditSystem *CreditSystem) UserController {
	return UserController{DB, creditSystem}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	userResponse := &models.UserResponse{
		ID:        currentUser.ID,
		Name:      currentUser.Name,
		Email:     currentUser.Email,
		Photo:     currentUser.Photo,
		Role:      currentUser.Role,
		Verified:  currentUser.Verified,
		Balance:   currentUser.Balance,
		Provider:  currentUser.Provider,
		CreatedAt: currentUser.CreatedAt,
		UpdatedAt: currentUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}

func (uc *UserController) AdminUpdateBalance(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	if currentUser.Role != constants.AppRoleAdmin {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Only Admin can update credits"})
		return
	}
	var payload *models.UpdateBalanceInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	newBalance, err := uc.cs.UpdateBalanceByUserEmail(payload.Email, payload.Amount)
	if err != nil {
		logger.Warning(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to update credits"})
		return
	}
	now := time.Now()
	bill := models.BillingHistory{
		Operator:          currentUser.Email,
		AccountEmail:      payload.Email,
		Amount:            payload.Amount,
		TransactionType:   constants.TransactionTypeAdmin,
		TransactionDetail: "Admin updates balance",
		TransactionTime:   now,
	}
	err = uc.DB.Create(&bill).Error
	if err != nil {
		logger.Warning("Failed to create billing history", err)
	}
	logger.Info("Create billing history", bill)

	var user models.User
	err = uc.DB.First(&user, "email = ?", payload.Email).Error
	if err != nil {
		logger.Warning(err.Error())
	}
	var firstName = user.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	emailData := utils.AccountEmailData{
		FirstName: firstName,
		Subject:   "Pointer.ai 充值提醒",
		Amount:    payload.Amount,
		Balance:   newBalance,
	}
	go func() {
		utils.SendAccountEmail(&user, &emailData, "chargeSucceed.html")

	}()
	logger.Info("Charge balance of user", payload.Email, "for amount", payload.Amount, "newBalance: ", newBalance)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"updated_email": payload.Email, "new_balance": newBalance}})
}

// WsEditingRoom 用户socket用来控制接收和写入,达到协同编辑
func (uc *UserController) WsEditingRoom(c *gin.Context) {
	ws.HandlerWs(c)

}

// CreateRoom 创建房间
func (uc *UserController) CreateRoom(c *gin.Context) {
	fileId := c.Param("fileId")
	room, Invitation, err := ws.HandleCreateRoom(fileId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": room, "invitation": Invitation})
}

// WsSubscribeEdit 用户socket用来控制接收和写入,达到协同编辑
//func (uc *UserController) WsSubscribeEdit(c *gin.Context) {
//	currentUser := c.MustGet("currentUser").(models.User)
//
//	var req *models.SubscribeRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
//		return
//	}
//	if err := ws.Subscriber.Subscribe(req.RoomId, subscription{Connection: c.Writer, UserInfo: &currentUser}); err != nil {
//		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success"})
//}

func (uc *UserController) GetUserAccessTokens(c *gin.Context) {
	inputToken := c.Query("token")
	user := c.MustGet("currentUser").(models.User)

	fmt.Println(user.Name, user.Password)

	client, err := gitea.CreateGiteaUserClient(user.Name, user.Password)
	if err != nil || client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if inputToken != "" {
		isValid, err := client.ValidateAccessToken(inputToken)
		if err != nil {
			logger.Warning(err.Error())
		}
		if isValid {
			c.JSON(http.StatusOK, gin.H{"status": "success", "data": inputToken})
			return
		}
	}

	_, token, err := client.CreateAccessToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": token})
}

func (uc *UserController) CreateGiteaRepo(c *gin.Context) {
	var payload gitea.CreateRepoOption
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user := c.MustGet("currentUser").(models.User)
	client, err := gitea.CreateGiteaUserClient(user.Name, user.Password)
	if err != nil || client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	repo, err := client.CreateRepo(payload.Name, payload.Description, payload.Private)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": repo})
}

func (uc *UserController) GetUserRepoList(c *gin.Context) {
	user := c.MustGet("currentUser").(models.User)

	// 使用 DefaultQuery 方法设置默认值
	page := c.DefaultQuery("page", "1")
	perPage := c.DefaultQuery("limit", "10")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid page number"})
		return
	}

	perPageInt, err := strconv.Atoi(perPage)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid limit"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	client, err := gitea.CreateGiteaUserClient(user.Name, user.Password)

	if err != nil || client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	repos, total, err := client.ListUserRepos(pageInt, perPageInt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": repos, "total": total})
}

func (uc *UserController) ValidateAccessToken(c *gin.Context) {
	user := c.MustGet("currentUser").(models.User)

	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "token is required"})
		return
	}

	client, err := gitea.CreateGiteaUserClient(user.Name, user.Password)

	if err != nil || client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	isValid, err := client.ValidateAccessToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": isValid})
}

func (uc *UserController) DeleteGiteaRepo(c *gin.Context) {
	repoName := c.Param("name")
	if repoName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "name is required"})
		return
	}

	user := c.MustGet("currentUser").(models.User)
	client, err := gitea.CreateGiteaUserClient(user.Name, user.Password)
	if err != nil || client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	err = client.DeleteRepo(user.Name, repoName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
