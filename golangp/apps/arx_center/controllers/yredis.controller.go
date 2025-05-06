package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/arxtect/ArxBackend/golangp/apps/arx_center/models"
	"github.com/arxtect/ArxBackend/golangp/common/constants"
	"github.com/arxtect/ArxBackend/golangp/common/logger"
	"github.com/arxtect/ArxBackend/golangp/common/utils"
	"github.com/arxtect/ArxBackend/golangp/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type YRedisController struct {
	DB     *gorm.DB
	Logger *log.Logger
}

func NewYRedisController(DB *gorm.DB, logger *log.Logger) YRedisController {
	return YRedisController{
		DB,
		logger,
	}
}

func (yac *YRedisController) getUserSharePosition(currentUser models.User, yroom string) (position int, err error) {
	if yroom == "" {
		return position, fmt.Errorf("Invalid room")
	}

	var room models.Yroom
	if result := yac.DB.Unscoped().First(&room, "name = ?", yroom); result.Error != nil { // 忽略软删除
		return position, result.Error
	}

	if currentUser.ID == room.CreateBy { //是否房间创建者
		return 0, nil
	}
	var positionMap map[string]int
	if err := json.Unmarshal([]byte(room.PositionData), &positionMap); err != nil {
		return position, fmt.Errorf("Error parsing position data: %v", err)
	}

	// Get the position for the current user
	if pos, exists := positionMap[currentUser.ID.String()]; exists {
		return pos, nil
	}

	return position, nil
}

// YRedisAuthToken 前端 y-redis 获取授权 Token
func (yac *YRedisController) YRedisAuthToken(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	userid := currentUser.ID.String() // 使用后端用户ID
	yroom := c.Param("room")

	appName := config.Env.Domain

	ecKey, kerr := utils.ToYRedisECDSAPrivateKey(config.Env.YredisAuthPrivateKey)
	if kerr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Internal Server Error"})
		yac.Logger.Printf("Failed to convert y-redis private key: %v", kerr)
		return
	}
	token, err := utils.CreateYRedisToken(config.Env.YredisAccessTokenExpiresIn,
		appName, userid, ecKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Internal Server Error"})
		yac.Logger.Printf("Failed to create y-redis token: %v", err)
		return
	}
	position, err := yac.getUserSharePosition(currentUser, yroom)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "token": token, "position": position})
}

// YRedisRoomPermissionCallback y-redis server 获取文档房间访问权限回调
func (yac *YRedisController) YRedisRoomPermissionCallback(c *gin.Context) {
	yroom := c.Param("room")     // 项目名称+后端用户ID
	yuserid := c.Param("userid") // 后端用户ID
	if yroom == "" || yuserid == "" {
		c.String(http.StatusBadRequest, "Invalid room or userid")
		return
	}

	tuserid, err := uuid.Parse(yuserid)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid room or userid")
		return
	}

	var room models.Yroom
	if result := yac.DB.First(&room, "name = ?", yroom); result.Error != nil {
		c.String(http.StatusBadRequest, "Invalid room or userid")
		return
	}

	canWrite := false //可写
	canRead := false  //可读
	if room.CreateBy == tuserid || room.IsPublic || utils.Contains(room.RW, tuserid) {
		canWrite = true
	}
	if room.CreateBy == tuserid || room.IsPublic || utils.Contains(room.R, tuserid) {
		canRead = true
	}

	if !canRead && !canWrite {
		c.String(http.StatusForbidden, "Access Denied")
		return
	}

	var yaccess string //读写权限
	if canWrite {
		yaccess = constants.YAccessReadAndWrite
	} else {
		yaccess = constants.YAccessReadOnly
	}

	permission := map[string]string{
		"yroom":   yroom,
		"yaccess": yaccess,
		"yuserid": yuserid,
	}

	c.JSON(http.StatusOK, permission)
	yac.Logger.Printf("y-redis room permission for room %v and user %v, access %v", yroom, yuserid, yaccess)
}

// YRedisYDocUpdateCallback y-redis worker 文档更新回调
func (yac *YRedisController) YRedisYDocUpdateCallback(c *gin.Context) {
	// 调用频率高，这里主要为了满足调用过程，返回成功即可
	c.String(http.StatusOK, "OK")
}

// YRedisRoomShareUserGet y-redis 文档房间获取当前用户权限
func (yac *YRedisController) YRedisRoomShareUserGet(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	yroom := c.Param("room")  // 项目名称+后端用户ID
	tuserid := currentUser.ID // 后端用户ID
	if yroom == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid room"})
		return
	}

	var room models.Yroom
	if result := yac.DB.First(&room, "name = ?", yroom); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid room or userid"})
		return
	}

	canWrite := false //可写
	canRead := false  //可读
	if room.CreateBy == tuserid || room.IsPublic || utils.Contains(room.RW, tuserid) {
		canWrite = true
	}
	if room.CreateBy == tuserid || room.IsPublic || utils.Contains(room.R, tuserid) {
		canRead = true
	}

	var yaccess string //读写权限
	if canWrite {
		yaccess = constants.YAccessReadAndWrite
	} else if canRead {
		yaccess = constants.YAccessReadOnly
	} else {
		yaccess = "no"
	}

	permission := map[string]string{
		"status": "success",
		"access": yaccess,
	}

	c.JSON(http.StatusOK, permission)
}

// YRedisRoomShareUserUpdate y-redis 文档房间分享添加、更新用户
func (yac *YRedisController) YRedisRoomShareUserUpdate(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	var payload *models.ShareProjectUserAddInput

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var targetUser models.User
	if result := yac.DB.First(&targetUser, "email = ?", strings.ToLower(payload.ShareEmail)); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or project"})
		return
	}
	if !targetUser.Verified {
		c.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "Invalid email or project"})
		return
	}
	if targetUser.ID == currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "You can't share a project with yourself"})
		return
	}

	var room models.Yroom
	resultR := yac.DB.Unscoped().First(&room, "name = ?", payload.ProjectName) // 忽略软删除
	now := time.Now()
	needCreate := resultR.Error == gorm.ErrRecordNotFound
	if needCreate {
		// 记录不存在，需要新建
		room = models.Yroom{
			Name:     payload.ProjectName, // 项目名称+后端用户ID
			CreateBy: currentUser.ID,
			RW:       make(models.UUIDArray, 0, 1),
			R:        make(models.UUIDArray, 0, 1),
			IsPublic: false,
			CreateAt: now,
		}
	} else if resultR.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or project"})
		logger.Warning("Failed to get room: %v", resultR.Error)
		return
	}

	if currentUser.ID != room.CreateBy && currentUser.Role != constants.AppRoleAdmin { //是否房间创建者
		c.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "You are not the creator of this room"})
		return
	}
	if targetUser.ID == room.CreateBy { // 目标用户是否是创建者
		c.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "You can't share a project with owner"})
		return
	}

	var isNewUser bool = false                  // 是否是新协作者
	if !utils.Contains(room.R, targetUser.ID) { // 读权限，默认读
		room.R = append(room.R, targetUser.ID)
		_ = room.AddKeyIfNotExists(targetUser.ID.String())
		isNewUser = true
	}

	if payload.Access == constants.YAccessReadAndWrite { // 读写权限
		if !utils.Contains(room.RW, targetUser.ID) {
			room.RW = append(room.RW, targetUser.ID)
		}
	} else {
		room.RW = utils.RemoveAll(room.RW, targetUser.ID)
	}

	room.UpdateAt = now

	if needCreate {
		if result := yac.DB.Create(&room); result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Duplicate project name"})
			return
		} else if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Internal Server Error"})
			logger.Warning("Failed to create room: %v", result.Error)
			return
		}
	} else {
		if result := yac.DB.Save(&room); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Internal Server Error"})
			logger.Warning("Failed to update room: %v", result.Error)
			return
		}
	}

	ProjectName := strings.Split(payload.ProjectName, currentUser.ID.String())[0]
	Subject := fmt.Sprintf(`"%s" - shared by %s`, ProjectName, currentUser.Email)
	// Send Email
	emailData := utils.ShareProjectEmailData{
		AuthorizedUser: targetUser.Name,
		SharerUser:     currentUser.Name,
		SharerEmail:    currentUser.Email,
		ProjectName:    ProjectName,
		ProjectLink:    payload.ShareLink,
		Subject:        Subject,
	}

	if isNewUser {
		go utils.SendShareProjectEmail(&targetUser, &emailData, "shareProject.html")
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// YRedisRoomCreateRoom y-redis 创建文档房间
func (yac *YRedisController) YRedisRoomCreateRoom(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	var payload *models.ShareProjectAddInput

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var room models.Yroom
	resultR := yac.DB.Unscoped().First(&room, "name = ?", payload.ProjectName) // 忽略软删除
	now := time.Now()
	needCreate := resultR.Error == gorm.ErrRecordNotFound
	if needCreate {
		// 记录不存在，需要新建
		room = models.Yroom{
			Name:         payload.ProjectName, // 项目名称+后端用户ID
			CreateBy:     currentUser.ID,
			RW:           make(models.UUIDArray, 0, 1),
			R:            make(models.UUIDArray, 0, 1),
			IsPublic:     false,
			CreateAt:     now,
			PositionData: "{}",
		}
	} else if resultR.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or project"})
		logger.Warning("Failed to get room: %v", resultR.Error)
		return
	}

	if currentUser.ID != room.CreateBy && currentUser.Role != constants.AppRoleAdmin { //是否房间创建者
		c.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "You are not the creator of this room"})
		return
	}

	if needCreate {
		if result := yac.DB.Create(&room); result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Duplicate project name"})
			return
		} else if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Internal Server Error"})
			logger.Warning("Failed to create room: %v", result.Error)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// YRedisRoomShareGet y-redis 获取文档房间用户权限列表
func (yac *YRedisController) YRedisRoomShareGet(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	yroom := c.Param("room") // 项目名称+后端用户ID

	if yroom == "" {
		c.JSON(http.StatusOK, gin.H{"status": "fail", "message": "Invalid room"})
		return
	}

	var room models.Yroom
	if result := yac.DB.Unscoped().First(&room, "name = ?", yroom); result.Error != nil { // 忽略软删除
		c.JSON(http.StatusOK, gin.H{"status": "fail", "message": "Invalid room"})
		return
	}

	if currentUser.ID != room.CreateBy && currentUser.Role != constants.AppRoleAdmin { //是否房间创建者
		c.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "You are not the creator of this room"})
		return
	}

	maxCapacity := len(room.RW) + len(room.R) // 预先计算 map 和切片的容量
	userPermMap := make(map[uuid.UUID]string, maxCapacity)
	for _, userID := range room.RW { // 读写权限用户
		userPermMap[userID] = constants.YAccessReadAndWrite
	}
	for _, userID := range room.R { // 读权限的用户
		if _, exists := userPermMap[userID]; !exists { // 如果用户不在读写权限列表中，则添加
			userPermMap[userID] = constants.YAccessReadOnly
		}
	}

	userPerms := make([]models.YRedisRoomShareGetResponseUserPermission, 0, maxCapacity) // 用户权限列表
	for userID, access := range userPermMap {
		var userPerm models.YRedisRoomShareGetResponseUserPermission
		userPerm.ID = userID
		userPerm.Access = access

		var user models.User
		if err := yac.DB.First(&user, "id = ?", userID).Error; err == nil { // 若找到用户
			userPerm.Email = user.Email
			userPerm.Name = user.Name
		}

		userPerms = append(userPerms, userPerm)
	}

	response := models.YRedisRoomShareGetResponse{ // 响应数据
		ID:        room.ID,
		Name:      room.Name,
		CreatorID: room.CreateBy,
		IsPublic:  room.IsPublic,
		CreateAt:  room.CreateAt,
		UpdateAt:  room.UpdateAt,
		IsClosed:  room.DeletedAt.Valid,
		Users:     userPerms,
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"room": response}})
}

// YRedisRoomShareUserRemove y-redis 文档房间分享移除用户
func (yac *YRedisController) YRedisRoomShareUserRemove(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	var payload *models.ShareProjectUserRemoveInput

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var targetUser models.User
	if result := yac.DB.First(&targetUser, "email = ?", strings.ToLower(payload.RemoveEmail)); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or project"})
		return
	}
	if targetUser.ID == currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "You cannot remove yourself from a project"})
		return
	}

	var room models.Yroom
	resultR := yac.DB.Unscoped().First(&room, "name = ?", payload.ProjectName) // 忽略软删除
	now := time.Now()
	if resultR.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or project"})
		logger.Warning("Failed to get room: %v", resultR.Error)
		return
	}

	if currentUser.ID != room.CreateBy && currentUser.Role != constants.AppRoleAdmin { //是否房间创建者
		c.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "You are not the creator of this room"})
		return
	}

	room.R = utils.RemoveAll(room.R, targetUser.ID)   //移除读
	room.RW = utils.RemoveAll(room.RW, targetUser.ID) //移除读写

	room.UpdateAt = now

	if result := yac.DB.Save(&room); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Internal Server Error"})
		logger.Warning("Failed to update room: %v", result.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// YRedisRoomShareClose y-redis 文档房间分享关闭
func (yac *YRedisController) YRedisRoomShareClose(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	var payload *models.ShareProjectCloseInput

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var room models.Yroom
	resultR := yac.DB.First(&room, "name = ?", payload.ProjectName)
	now := time.Now()
	if resultR.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid project"})
		logger.Warning("Failed to get room: %v", resultR.Error)
		return
	}

	if currentUser.ID != room.CreateBy && currentUser.Role != constants.AppRoleAdmin { //是否房间创建者
		c.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "You are not the creator of this room"})
		return
	}

	room.UpdateAt = now

	if result := yac.DB.Delete(&room); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Internal Server Error"})
		logger.Warning("Failed to delete room: %v", result.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// YRedisRoomShareReopen y-redis 文档房间分享重新打开
func (yac *YRedisController) YRedisRoomShareReopen(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	var payload *models.ShareProjectCloseInput

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var room models.Yroom
	resultR := yac.DB.Unscoped().First(&room, "name = ?", payload.ProjectName) // 忽略软删除
	now := time.Now()
	if resultR.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid project"})
		logger.Warning("Failed to get room: %v", resultR.Error)
		return
	}

	if currentUser.ID != room.CreateBy && currentUser.Role != constants.AppRoleAdmin { //是否房间创建者
		c.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "You are not the creator of this room"})
		return
	}

	room.UpdateAt = now

	if room.DeletedAt.Valid { //如果房间被删除过，则重新打开
		room.DeletedAt.Time = time.Time{}
		room.DeletedAt.Valid = false
	}

	if result := yac.DB.Save(&room); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Internal Server Error"})
		logger.Warning("Failed to open room: %v", result.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// YRedisRoomShareGet y-redis 文档房间分享获取
