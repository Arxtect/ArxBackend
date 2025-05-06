package main

import (
	"fmt"
	"log"

	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/gitea"
	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/models"
	"github.com/Arxtect/ArxBackend/golangp/common/initializers"
	"github.com/Arxtect/ArxBackend/golangp/config"
	"gorm.io/gorm"
)

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
	// 调用 AutoSyncUser 函数
	AutoSyncUser(initializers.DB)
}

//ALTER TABLE user
//ADD CONSTRAINT unique_name UNIQUE (name);
