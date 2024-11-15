package dify

import (
	"encoding/json"
	"fmt"
	"hajime/golangp/apps/hajime_center/constants"
	"hajime/golangp/apps/hajime_center/initializers"
	"hajime/golangp/common/logging"
)

type UserLoginParams struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"remember_me"`
}

type Data struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserLoginResponse struct {
	Result string `json:"result"`
	Data   Data   `json:"data"`
}

func (dc *DifyClient) UserLogin(email string, password string) (result UserLoginResponse, err error) {
	var payload = UserLoginParams{
		Email:      email,
		Password:   password,
		RememberMe: true,
	}
	api := dc.GetConsoleAPI(CONSOLE_API_LOGIN)

	code, body, err := SendPostRequestToConsole(dc, api, payload)

	logging.Warning("code: %d, body: %s\n", code, string(body))
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

func (dc *DifyClient) GetUserToken(role string) (token string, err error) {
	// 检查是否已经存在 ACCESS_TOKEN
	if dc.ConsoleToken != "" {
		return dc.ConsoleToken, nil
	}
	// 如果 role 为空，则设置默认值
	if role == "" {
		role = constants.RoleUser
	}

	config, _ := initializers.LoadEnv(".")

	result := UserLoginResponse{}

	if role == constants.RoleAdmin {
		// 不存在 ACCESS_TOKEN，进行登录操作
		result, err = dc.UserLogin(config.DifyConsoleEmail, config.DifyConsolePassword)
		if err != nil {
			logging.Warning("failed to login: %v\n", err)
			return "", err
		}
	} else if role == constants.RoleEditor {
		result, err = dc.UserLogin(config.DifyEditorEmail, config.DifyEditorPassword)
		if err != nil {
			logging.Warning("failed to login: %v\n", err)
			return "", err
		}
	} else {
		result, err = dc.UserLogin(config.DifyUserEmail, config.DifyUserPassword)
		if err != nil {
			logging.Warning("failed to login: %v\n", err)
			return "", err
		}
	}

	// 更新 ACCESS_TOKEN
	dc.ConsoleToken = result.Data.AccessToken
	return result.Data.AccessToken, nil
}
