package utils

import (
	"log"
	"os"
	"testing"

	"github.com/Arxtect/ArxBackend/golangp/common/utils"
	"github.com/Arxtect/ArxBackend/golangp/config"
)

func TestSendAccountEmail(t *testing.T) {
	err := os.Chdir("/home/dechang/Einstein")
	if err != nil {
		panic(err)
	}

	err = config.LoadEnv("config/settings-dev.yml")
	if err != nil {
		log.Fatalf("Load failed: %v\n", err)
		return
	}

	emailData := utils.AccountEmailData{
		URL:              "https://arxtect.com",
		VerificationCode: "123456",
		FirstName:        "test",
		Subject:          "Welcome to our service for test",
	}

	// Call the function
	utils.SendEmail("zhengdevin10@gmail.com", emailData.Subject, emailData, "resetPassword.html")
}
