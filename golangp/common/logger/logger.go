package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

var Logger *log.Logger
var once sync.Once

func init() {
	once.Do(func() {
		// 检查环境变量 LOG_FILE_PATH
		logFilePath := os.Getenv("LOG_FILE_PATH")
		fmt.Println("logFilePath:", logFilePath)
		if logFilePath != "" {
			// 打开或创建日志文件
			file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("Failed to open log file: %v", err)
			}
			// 初始化日志记录器，输出到文件
			Logger = log.New(file, "", 0)
		} else {
			// 初始化日志记录器，输出到标准输出
			Logger = log.New(os.Stdout, "", 0)
		}
	})
}

func logWithCaller(prefix, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	message := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	Logger.Printf("%-8s %s %s:%d: %s", prefix, timestamp, file, line, message)
}

func Info(format string, args ...interface{}) {
	logWithCaller("[INFO]", format, args...)
}

func Danger(format string, args ...interface{}) {
	logWithCaller("[ERROR]", format, args...)
	message := fmt.Sprintf("[ERROR] "+format, args...)
	Logger.Fatal(message)
}

func Warning(format string, args ...interface{}) {
	logWithCaller("[WARNING]", format, args...)
}

func DeBug(format string, args ...interface{}) {
	logWithCaller("[DeBug]", format, args...)
}
