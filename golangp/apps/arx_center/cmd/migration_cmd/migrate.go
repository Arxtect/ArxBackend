package migration_cmd

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/models"
	"github.com/Arxtect/ArxBackend/golangp/common/constants"
	"github.com/Arxtect/ArxBackend/golangp/common/initializers"
	"github.com/Arxtect/ArxBackend/golangp/common/logger"
	"github.com/Arxtect/ArxBackend/golangp/config"
	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/migrate/migration_user"
	"github.com/Arxtect/ArxBackend/golangp/common/utils"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var configYml string

// Add MigrateCmd
var MigrateCmd = &cobra.Command{
	Use:          "migrate",
	Short:        "Run database migrations",
	Example:      "ArxBackend migrate -c config/settings-dev.yml",
	SilenceUsage: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		setupMigrate()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return migrate()
	},
}

func init() {
	MigrateCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings-dev.yml", "Run migrations with provided configuration file")
}

func setupMigrate() {
	log.Println("🚗 Loading configuration for migration...")
	err := config.LoadEnv(configYml)
	if err != nil {
		log.Printf("❌ Failed to load configuration: %v", err)
		return
	}
	log.Printf("✅ Configuration loaded successfully in %s mode", config.Env.Mode)

	initializers.ConnectDB(&config.Env)
	log.Printf("✅ Database connection established in %s mode", config.Env.Mode)
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

func migrate() error {
	log.Println("🚀 Starting database migrations...")

	// Run your migrations here
	err := initializers.DB.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.S3File{},
		&models.ProjectMember{},
	)
	if err != nil {
		log.Printf("❌ Migration failed: %v", err)
		return err
	}

	err = initializers.DB.SetupJoinTable(&models.Project{}, "Members", &models.ProjectMember{})
	if err != nil {
		log.Printf("❌ SetupJoinTable failed: %v", err)
		return err
	}

	err = initializers.DB.SetupJoinTable(&models.User{}, "SharedProjects", &models.ProjectMember{})
	if err != nil {
		log.Printf("❌ SetupJoinTable failed: %v", err)
		return err
	}

	SetupAdmin(initializers.DB)
	migration_user.User(initializers.DB)
	log.Println("✅ Database migrations completed successfully")
	return nil
}
