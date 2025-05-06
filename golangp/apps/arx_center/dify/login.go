package dify

import (
	"github.com/arxtect/ArxBackend/golangp/common/logger"
	"github.com/arxtect/ArxBackend/golangp/config"
	"encoding/json"
	"fmt"

	"github.com/toheart/functrace"
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
	defer functrace.Trace([]interface {
	}{dc, email, password})()
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
	defer functrace.Trace([]interface {
	}{dc})()

	if dc.ConsoleToken != "" {
		return dc.ConsoleToken, nil
	}

	config := config.Env

	result, err := dc.UserLogin(config.DifyConsoleEmail, config.DifyConsolePassword)
	if err != nil {
		logger.Warning("failed to login: %v\n", err)
		return "", err
	}

	dc.ConsoleToken = result.Data
	return result.Data, nil
}
