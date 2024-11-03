package main

import (
	"fmt"
	"github.com/goccy/go-json"
	"hajime/golangp/libs"
)

func main() {
	fmt.Println("Demo says:", libs.Hello())

	data := map[string]interface{}{
		"message": "Hello, World!",
	}

	// 使用 go-json 库进行编码
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	fmt.Println("JSON Output:", string(jsonData))
}
