package dify

import (
	"encoding/json"
	"fmt"
	"github.com/toheart/functrace"
)

type ConsoleWorkspace struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	Plan           string      `json:"plan"`
	Status         string      `json:"status"`
	CreatedAt      int         `json:"created_at"`
	Role           string      `json:"role"`
	InTrial        interface{} `json:"in_trial"`
	TrialEndReason interface{} `json:"trial_end_reason"`
	CustomConfig   interface{} `json:"custom_config"`
}

func (dc *DifyClient) GetConsoleWorkspaceCurrent() (result ConsoleWorkspace, err error) {
	defer functrace.Trace([]interface {
	}{dc})()
	api := dc.GetConsoleAPI(CONSOLE_API_WORKSPACE_CURRENT)

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
