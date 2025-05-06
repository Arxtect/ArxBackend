package migration_user

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/gitea"

	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/models"
	"github.com/Arxtect/ArxBackend/golangp/common/constants"
	"github.com/Arxtect/ArxBackend/golangp/common/initializers"
	"github.com/Arxtect/ArxBackend/golangp/config"
	"github.com/Arxtect/ArxBackend/golangp/common/utils"
	"gorm.io/gorm"
)

func User(db *gorm.DB) {
	var userList []map[string]string
	for i := 1; i <= 10; i++ {
		user := map[string]string{
			"Name":     "inside_test" + strconv.Itoa(i),
			"Email":    "inside_test" + strconv.Itoa(i) + "@gmail.com",
			"Password": "inside_test" + strconv.Itoa(i) + "@gmail.com",
		}
		userList = append(userList, user)
	}

	now := time.Now()
	for _, user := range userList {
		hashedPassword, err := utils.HashPassword(user["Password"])
		if err != nil {
			log.Printf("Error hashing password for user %s: %v\n", user["Name"], err)
			continue
		}
		payload := models.User{
			Name:      user["Name"],
			Email:     strings.ToLower(user["Email"]),
			Password:  hashedPassword,
			Role:      constants.AppRoleUser,
			Verified:  true,
			Photo:     "test",
			Provider:  "local",
			CreatedAt: now,
			UpdatedAt: now,
		}

		// 将用户插入到数据库中
		if err := db.Create(&payload).Error; err != nil {
			log.Printf("Error creating user %s: %v\n", payload.Name, err)
		}
	}
}

func AutoSyncUser(db *gorm.DB) {
	var userList []models.User

	// 从数据库中获取用户列表
	if err := db.Find(&userList).Error; err != nil {
		fmt.Printf("Error fetching users from database: %v\n", err)
		return
	}

	// 对每个用户调用 CreateUser 函数
	for _, user := range userList {
		giteaUser, err := gitea.CreateUser(user.Name, user.Password, user.Email)
		if err != nil {
			if err.Error() == "user already exists" {
				fmt.Printf("User %s already exists, skipping...\n", user.Name)
				continue
			}
			fmt.Printf("Error creating user %s: %v\n", user.Name, err)
		} else {
			fmt.Printf("User %s created successfully.\n", user.Name)
			fmt.Printf("Gitea User: %s created successfully.\n", giteaUser.UserName)
		}
	}
}

func main() {
	// 初始化数据库连接
	err := config.LoadEnv("config/settings-dev.yml")
	if err != nil {
		log.Fatalf("Load failed: %v\n", err)
		return
	}
	log.Println("🚗 Loading env is success....", config.Env.Mode)

	initializers.ConnectDB(&config.Env)
	defer func() {
		db, _ := initializers.DB.DB()
		db.Close()
	}()

	// 调用 User 函数
	User(initializers.DB)
	AutoSyncUser(initializers.DB)
}
