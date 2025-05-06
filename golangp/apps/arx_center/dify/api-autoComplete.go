package dify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/toheart/functrace"
	"io/ioutil"
	"net/http"
	"time"
)

type AutoCompletePayload struct {
	ResponseMode string                    `json:"response_mode,omitempty"`
	Inputs       map[string]interface{}    `json:"inputs"`
	User         string                    `json:"user,omitempty"`
	Files        []ChatMessagesPayloadFile `json:"files"`
	AcessToken   string                    `json:"access_token"`
}

type AutoCompleteResponse struct {
	WorkflowRunID string `json:"workflow_run_id"`
	TaskID        string `json:"task_id"`
	Data          Data   `json:"data"`
}

type Data struct {
	ID          string   `json:"id"`
	WorkflowID  string   `json:"workflow_id"`
	Status      string   `json:"status"`
	Outputs     *Outputs `json:"outputs,omitempty"`
	Error       *string  `json:"error,omitempty"`
	ElapsedTime *float64 `json:"elapsed_time,omitempty"`
	TotalTokens *int     `json:"total_tokens,omitempty"`
	TotalSteps  int      `json:"total_steps"`
	CreatedAt   int64    `json:"created_at"`
	FinishedAt  int64    `json:"finished_at"`
}

type Outputs struct {
	Suggestion string `json:"suggestion"`
}

type AutoCompleteErrorResponse struct {
	Code    string `json:"status"`
	Message string `json:"message"`
}

func (dc *DifyClient) AutoComplete(payload *AutoCompletePayload) (AutoCompleteResponse, error) {
	defer functrace.Trace([]interface {
	}{dc, payload})()
	var result AutoCompleteResponse

	payload.User = dc.User

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return result, fmt.Errorf("failed to marshal request payload: %v", err)
	}

	api := dc.Host + API_AUTOCOMPLETE

	req, err := http.NewRequest("POST", api, bytes.NewBuffer(reqBody))
	if err != nil {
		return result, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", payload.AcessToken)

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	fmt.Print("AutoComplete Request: ", req)
	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		var errorResponse AutoCompleteErrorResponse
		if err := json.Unmarshal(bodyBytes, &errorResponse); err != nil {
			return result, fmt.Errorf("server error, response code: %d", resp.StatusCode)
		}

		return result, fmt.Errorf("message: %s, code: %s", errorResponse.Message, errorResponse.Code)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("failed to read response body: %v", err)
	}

	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return result, nil
}
