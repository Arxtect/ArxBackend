package model

import (
	"fmt"
	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/models"
	"github.com/Arxtect/ArxBackend/golangp/common/initializers"
	"github.com/Arxtect/ArxBackend/golangp/common/utils"
	"math/rand"
	"time"
)

// Random_create_user 随机生成用户
func Random_create_user() {
	//err := config.LoadEnv("config/settings-test.yml")
	//if err != nil {
	//	fmt.Printf("🚀 Could not load environment variables %s", err.Error())
	//}
	//initializers.ConnectDB(&config.Env)
	//
	//// Migrate the schema
	//_ = initializers.DB.AutoMigrate(&models.User{})

	// Create 100 users
	hashedPassword, _ := utils.HashPassword("jancsitech")
	for i := 0; i < 5; i++ {

		user := models.User{
			Name:      generateRandomString(10),
			Email:     fmt.Sprintf("%s@example.com", generateRandomString(10)),
			Password:  hashedPassword,
			Role:      "user",
			Provider:  "local",
			Photo:     "",
			Verified:  true,
			Balance:   0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := initializers.DB.Create(&user).Error
		if err != nil {
			fmt.Printf("🚀 Could not create user %s", err.Error())
		}
	}
}

func generateRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// Random_create_tags 随机生成tags
func Random_create_tags() {
	//err := config.LoadEnv("config/settings-test.yml")
	//if err != nil {
	//	fmt.Printf("🚀 Could not load environment variables %s", err.Error())
	//}
	//initializers.ConnectDB(&config.Env)
	//
	//// Migrate the schema
	//_ = initializers.DB.AutoMigrate(&models.Tag{})

	// Create 6 tags
	for i := 0; i < 6; i++ {
		tag := models.Tag{
			Name: generateRandomString(4),
		}
		initializers.DB.Create(&tag)
	}
}
