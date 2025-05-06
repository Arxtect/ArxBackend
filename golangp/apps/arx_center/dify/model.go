package dify

import (
	"encoding/json"
	"fmt"
	"github.com/toheart/functrace"
)

type GetCurrentWorkspaceLLMDefaultModelResponse struct {
	Data any `json:"data"`
}

func (dc *DifyClient) GetCurrentWorkspaceLLMDefaultModel() (result GetCurrentWorkspaceLLMDefaultModelResponse, err error) {
	defer functrace.Trace([]interface {
	}{dc})()
	api := dc.GetConsoleAPI(CONSOLE_API_CURRENT_WORKSPACE_LLM_MODEL)

	code, body, err := SendGetRequestToConsole(dc, api)

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

func (dc *DifyClient) GetCurrentWorkspaceLLMModel() (result GetCurrentWorkspaceLLMDefaultModelResponse, err error) {
	defer functrace.Trace([]interface {
	}{dc})()
	api := dc.GetConsoleAPI(CONSOLE_API_WORKSPACES_LLM_MODEL)

	code, body, err := SendGetRequestToConsole(dc, api)

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
