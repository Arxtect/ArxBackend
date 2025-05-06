package motest

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/models"
	"github.com/Arxtect/ArxBackend/golangp/common/constants"
	"github.com/Arxtect/ArxBackend/golangp/common/initializers"
	"github.com/Arxtect/ArxBackend/golangp/common/logger"
	"github.com/Arxtect/ArxBackend/golangp/common/utils"
	"github.com/Arxtect/ArxBackend/golangp/config"
	"gorm.io/gorm"
)

// TestModeMigrate test 环境需要迁移表结构
func TestModeMigrate() {
	err := initializers.DB.AutoMigrate(&models.User{}, &models.BillingHistory{}, &models.Post{}, &models.DocumentMigrate{})
	if err != nil {
		logger.Danger("🚀 Could not migrate User model", err)
		return
	}
	//SetupAdmin(initializers.DB)
	// 执行原始 SQL 创建索引
	result := initializers.DB.Exec(`
CREATE TEXT SEARCH CONFIGURATION simple_pdfs ( COPY = simple );
CREATE INDEX documents_content_index ON documents USING gin(to_tsvector('simple_pdfs', content));
`)
	if result.Error != nil {
		log.Println("Error creating index:", result.Error)
	} else {
		log.Println("Index created successfully")
	}
	//model.Random_create_user()
	//model.Random_create_tags()
	RecoveryDataByProd(initializers.DB)
	fmt.Println("👍 Migration all complete by test mode")
}

func SetupAdmin(DB *gorm.DB) {
	adminPassword := config.Env.AdminPassword

	hashedPassword, err := utils.HashPassword(adminPassword)
	if err != nil {
		logger.Danger("Error hashing password %s", err.Error())
	}

	removeAllAdmins(DB)

	for _, adminEmail := range config.Env.AdminEmail {
		now := time.Now()
		newUser := models.User{
			Name:      "Admin Admin",
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
