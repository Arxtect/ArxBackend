package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	models "github.com/Arxtect/ArxBackend/golangp/apps/arx_center/models"
	"github.com/Arxtect/ArxBackend/golangp/common/constants"
	"github.com/Arxtect/ArxBackend/golangp/common/initializers"
	"github.com/Arxtect/ArxBackend/golangp/common/logger"
	"github.com/Arxtect/ArxBackend/golangp/config"
	"github.com/Arxtect/ArxBackend/golangp/common/utils"
	"gorm.io/gorm"
)

func init() {
	err := config.LoadEnv("config/settings-dev.yml")
	if err != nil {
		logger.Danger("🚀 Could not load environment variables %s", err.Error())
	}

	initializers.ConnectDB(&config.Env)
}

func removeAllAdmins(DB *gorm.DB) {
	var adminUsers []models.User
	res := DB.Find(&adminUsers, "role = ?", constants.AppRoleAdmin)
	if res.Error != nil {
		logger.Warning("Error finding admin users %s", res.Error.Error())
	}

	for _, adminUser := range adminUsers {
		userToDelete := adminUser.Email
		res := DB.Delete(&adminUser)
		if res.Error != nil {
			logger.Warning("Error deleting admin user %s", res.Error.Error())
		}
		logger.Info("Previous admin user %s deleted successfully", userToDelete)
	}
}

func SetupAdmin(DB *gorm.DB) {
	// kong.Parse(&openai_config.CLI)
	// OpenaiConfig := openai_config.LoadOpenAIConfig()
	// adminPassword := OpenaiConfig.AdminPassword
	config := config.Env
	adminPassword := config.AdminPassword

	hashedPassword, err := utils.HashPassword(adminPassword)
	if err != nil {
		logger.Danger("Error hashing password %s", err.Error())
	}
	removeAllAdmins(DB)

	for index, adminEmail := range config.AdminEmail {
		now := time.Now()
		newUser := models.User{
			Name:      "Admin" + strconv.Itoa(index),
			Email:     strings.ToLower(adminEmail),
			Password:  hashedPassword,
			Role:      constants.AppRoleAdmin,
			Verified:  true,
			Photo:     "test",
			Provider:  "local",
			CreatedAt: now,
			UpdatedAt: now,
		}

		var adminUser models.User
		res := DB.First(&adminUser, "email = ?", adminEmail)
		if res.Error != nil {
			logger.Info("Admin user %s does not exist, creating one", adminEmail)
		} else {
			res := DB.Delete(&adminUser)
			if res.Error != nil {
				logger.Warning("Error deleting exist admin user %s", res.Error.Error())
			}
			logger.Info("Existing Admin user deleted successfully")
		}

		result := DB.Create(&newUser)

		if result.Error != nil && strings.Contains(result.Error.Error(), "duplicated key not allowed") {
			logger.Warning("Admin email already exists")
			return
		} else if result.Error != nil {
			logger.Danger("Error creating admin user", result.Error)
		}

		logger.Info("Admin user %s created successfully", adminEmail)
	}

}

func main() {
	err := initializers.DB.AutoMigrate(&models.User{}, &models.BillingHistory{}, &models.Post{}, &models.Project{}, &models.S3File{})
	if err != nil {
		logger.Danger("🚀 Could not migrate User model", err)
	}
	SetupAdmin(initializers.DB)

	//迁移表结构 (再次执行会覆盖之前的表结构)
	//model.MigrateDocuments()

	_ = initializers.DB.AutoMigrate(&models.Prompt{})

	fmt.Println("👍 Migration all complete")

	err = initializers.DB.AutoMigrate(&models.Yroom{})
	if err != nil {
		logger.Danger("🚀 Could not migrate Yroom model %v", err)
	}
}
