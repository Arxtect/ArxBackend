package initializers

import (
	"github.com/Arxtect/ArxBackend/golangp/common/logger"
	logging "github.com/Arxtect/ArxBackend/golangp/common/logger"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var DBDify *gorm.DB

func ConnectDB(config *Config) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai", config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort, config.DBSslMode)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logging.Danger("Failed to connect to the Database")
	}
	fmt.Println("🚀 Connected Successfully to the Database")
}

func ConnectDBDify(config *Config) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", config.DBHostDify, config.DBUserNameDify, config.DBUserPasswordDify, config.DBNameDify, config.DBPortDify)

	DBDify, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logging.Danger("Failed to connect to the Dify Database")
	}
	fmt.Println("🚀 Connected Successfully to the Dify Database")
}

func TestConnectDb() {
	var err error
	dns := `host=10.10.101.123 user=chatCRO password=chatCRO dbname=chatCRO port=5432 sslmode=disable TimeZone=Asia/Shanghai`
	DB, err = gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		logger.Danger("Failed to connect to the Database")
	}
	fmt.Println("🚗  Connected Successfully to the Database")
}
