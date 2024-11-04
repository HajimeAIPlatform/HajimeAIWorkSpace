package logging

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var Logger *log.Logger
var once sync.Once

func init() {
	once.Do(func() {
		// 检查环境变量 LOG_FILE_PATH
		logFilePath := os.Getenv("LOG_FILE_PATH")
		if logFilePath != "" {
			// 打开或创建日志文件
			file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("Failed to open log file: %v", err)
			}
			// 初始化日志记录器，输出到文件
			Logger = log.New(file, "INFO", log.Ldate|log.Ltime|log.Lshortfile)
		} else {
			// 初始化日志记录器，输出到标准输出
			Logger = log.New(os.Stdout, "INFO", log.Ldate|log.Ltime|log.Lshortfile)
		}
	})
}

func Info(format string, args ...interface{}) {
	Logger.SetPrefix("[INFO]")
	message := fmt.Sprintf(format, args...)

	Logger.Println(message)
}

func Danger(format string, args ...interface{}) {
	Logger.SetPrefix("[ERROR]")
	message := fmt.Sprintf(format, args...)
	Logger.Fatal(message)
}

func Warning(format string, args ...interface{}) {
	Logger.SetPrefix("[WARNING]")
	message := fmt.Sprintf(format, args...)
	Logger.Println(message)
}

func DeBug(format string, args ...interface{}) {
	Logger.SetPrefix("[DeBug]")
	message := fmt.Sprintf(format, args...)
	Logger.Println(message)
}
