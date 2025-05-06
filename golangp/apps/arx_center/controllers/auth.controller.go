package controllers

import (
	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/gitea"
	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/models"
	"github.com/Arxtect/ArxBackend/golangp/common/constants"
	"github.com/Arxtect/ArxBackend/golangp/common/utils"
	"github.com/Arxtect/ArxBackend/golangp/config"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

// SignUpUser SignUp User
func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var payload *models.SignUpInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	now := time.Now()
	newUser := models.User{
		Name:      payload.Name,
		Email:     strings.ToLower(payload.Email),
		Password:  hashedPassword,
		Role:      constants.AppRoleUser,
		Verified:  false,
		Photo:     "test",
		Provider:  "local",
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := ac.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists, " +
			"try use forget password to reset it."})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened"})
		return
	}
	// Generate Verification Code
	code := randstr.String(6)

	verificationCode := utils.Encode(code)

	// Update User in Database
	newUser.VerificationCode = verificationCode
	ac.DB.Save(newUser)

	configCopy := config.Env

	var firstName = newUser.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// ? Send Email
	emailData := utils.AccountEmailData{
		URL:              configCopy.ClientOrigin + "/verify-email?code=" + code,
		VerificationCode: code,
		FirstName:        firstName,
		Subject:          "Verify your arXtect email address",
	}

	go utils.SendAccountEmail(&newUser, &emailData, "verificationCode.html")

	message := "We sent an email with a verification code to " + newUser.Email
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": message})

}

func (ac *AuthController) GetDomains(ctx *gin.Context, config config.Config) string {
	requestOrigin := ctx.Request.Header.Get("Referer")
	var requestHost string
	if requestOrigin != "" {
		u, err := url.Parse(requestOrigin)
		if err == nil {
			requestHost = u.Host
		}
	}
	if requestHost == "" {
		requestHost = ctx.Request.Host
	}

	domains := strings.Split(config.Domain, ",")
	for _, domain := range domains {
		domain = strings.TrimSpace(domain)
		if strings.Contains(requestHost, domain) {
			return domain
		}
	}

	if len(domains) > 0 {
		return strings.TrimSpace(domains[0])
	}
	return ""
}

func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var payload *models.SignInInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	if !user.Verified {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "Please verify your email"})
		return
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	configCopy := config.Env

	// Generate Tokens
	accessToken, err := utils.CreateToken(configCopy.AccessTokenExpiresIn, user.ID, configCopy.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refreshToken, err := utils.CreateToken(configCopy.RefreshTokenExpiresIn, user.ID, configCopy.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	domain := ac.GetDomains(ctx, configCopy)

	ctx.SetCookie("access_token", accessToken, configCopy.AccessTokenMaxAge*60, "/", domain, false, true)
	ctx.SetCookie("refresh_token", refreshToken, configCopy.RefreshTokenMaxAge*60, "/", domain, false, true)
	ctx.SetCookie("logged_in", "true", configCopy.AccessTokenMaxAge*60, "/", domain, false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

// RefreshAccessToken Refresh Access Token
func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	cookie, err := ctx.Cookie("refresh_token")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	configCopy := config.Env

	sub, err := utils.ValidateToken(cookie, configCopy.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "id = ?", fmt.Sprint(sub))
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}

	accessToken, err := utils.CreateToken(configCopy.AccessTokenExpiresIn, user.ID, configCopy.AccessTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	domain := ac.GetDomains(ctx, configCopy)
	ctx.SetCookie("access_token", accessToken, configCopy.AccessTokenMaxAge*60, "/", domain, false, true)
	ctx.SetCookie("logged_in", "true", configCopy.AccessTokenMaxAge*60, "/", domain, false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	configCopy := config.Env
	domain := ac.GetDomains(ctx, configCopy)
	ctx.SetCookie("access_token", "", -1, "/", domain, false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", domain, false, true)
	ctx.SetCookie("logged_in", "", -1, "/", domain, false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

// VerifyEmail [...] Verify Email
func (ac *AuthController) VerifyEmail(ctx *gin.Context) {

	code := ctx.Params.ByName("verificationCode")
	verificationCode := utils.Encode(code)

	var updatedUser models.User
	result := ac.DB.First(&updatedUser, "verification_code = ?", verificationCode)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid verification code or user doesn't exists"})
		return
	}

	if updatedUser.Verified {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User already verified"})
		return
	}

	updatedUser.VerificationCode = ""
	updatedUser.Verified = true
	ac.DB.Save(&updatedUser)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Email verified successfully"})
}

func (ac *AuthController) ForgotPassword(ctx *gin.Context) {
	var payload *models.ForgotPasswordInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := "You will receive a reset email if user with that email exist"

	var user models.User
	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	configCopy := config.Env

	// Generate Verification Code
	resetToken := randstr.String(6)

	passwordResetToken := utils.Encode(resetToken)
	user.PasswordResetToken = passwordResetToken
	user.PasswordResetAt = time.Now().Add(time.Minute * 15)
	ac.DB.Save(&user)

	var firstName = user.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// ? Send Email
	emailData := utils.AccountEmailData{
		URL:              configCopy.ClientOrigin + "/#/resetpassword/" + resetToken,
		VerificationCode: resetToken,
		FirstName:        firstName,
		Subject:          "Your password reset token (valid for 10min)",
	}

	go utils.SendAccountEmail(&user, &emailData, "resetPassword.html")

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}

func (ac *AuthController) ResetPassword(ctx *gin.Context) {
	var payload *models.ResetPasswordInput
	resetToken := ctx.Params.ByName("resetToken")
	configCopy := config.Env

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	if payload.Password != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	hashedPassword, _ := utils.HashPassword(payload.Password)

	passwordResetToken := utils.Encode(resetToken)

	var updatedUser models.User
	result := ac.DB.First(&updatedUser, "password_reset_token = ? AND password_reset_at > ?", passwordResetToken, time.Now())
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "The reset token is invalid or has expired"})
		return
	}

	updatedUser.Password = hashedPassword
	updatedUser.Verified = true
	updatedUser.PasswordResetToken = ""

	// TODO 修改相关的git账号信息，再去插入库
	_, err := gitea.UpdateUserPassword(updatedUser.Name, updatedUser.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "The reset token is invalid or has expired"})
		return
	}

	ac.DB.Save(&updatedUser)
	domain := ac.GetDomains(ctx, configCopy)
	ctx.SetCookie("token", "", -1, "/", domain, false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Password data updated successfully"})
}
