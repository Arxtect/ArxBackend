package middleware

import (
	"github.com/arxtect/ArxBackend/golangp/apps/arx_center/models"
	"github.com/arxtect/ArxBackend/golangp/common/initializers"
	"github.com/arxtect/ArxBackend/golangp/common/utils"
	"github.com/arxtect/ArxBackend/golangp/config"
	"fmt"
	"net/http"
	"strings"

	"github.com/toheart/functrace"

	"github.com/gin-gonic/gin"
)

func DeserializeUser() gin.HandlerFunc {
	defer functrace.Trace([]interface {
	}{})()
	return func(ctx *gin.Context) {
		var accessToken string
		cookie, err := ctx.Cookie("access_token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
		} else if err == nil {
			accessToken = cookie
		}

		if accessToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		configCopy := config.Env
		sub, err := utils.ValidateToken(accessToken, configCopy.AccessTokenPublicKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		var user models.User
		result := initializers.DB.First(&user, "id = ?", fmt.Sprint(sub))
		if result.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}
