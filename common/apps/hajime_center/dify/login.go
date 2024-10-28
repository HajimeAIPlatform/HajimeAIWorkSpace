package dify

import (
	"encoding/json"
	"fmt"
	"HajimeAIWorkSpace/common/apps/hajime_center/logger"
	"HajimeAIWorkSpace/common/apps/hajime_center/initializers"
)

type UserLoginParams struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"remember_me"`
}

type UserLoginResponse struct {
	Result string `json:"result"`
	Data   string `json:"data"`
}

func (dc *DifyClient) UserLogin(email string, password string) (result UserLoginResponse, err error) {
	var payload = UserLoginParams{
		Email:      email,
		Password:   password,
		RememberMe: true,
	}

	api := dc.GetConsoleAPI(CONSOLE_API_LOGIN)

	code, body, err := SendPostRequestToConsole(dc, api, payload)

	err = CommonRiskForSendRequest(code, err)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal the response: %v", err)
	}

	return result, nil
}

func (dc *DifyClient) GetUserToken() (string, error) {
	// 检查是否已经存在 ACCESS_TOKEN
	if dc.ConsoleToken != "" {
		return dc.ConsoleToken, nil
	}

	config, _ := initializers.LoadEnv(".")

	// 不存在 ACCESS_TOKEN，进行登录操作
	result, err := dc.UserLogin(config.DifyConsoleEmail, config.DifyConsolePassword)
	if err != nil {
		logger.Warning("failed to login: %v\n", err)
		return "", err
	}

	// 更新 ACCESS_TOKEN
	dc.ConsoleToken = result.Data
	return result.Data, nil
}
