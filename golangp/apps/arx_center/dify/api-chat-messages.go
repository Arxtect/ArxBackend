package dify

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/toheart/functrace"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type ChatMessagesPayload struct {
	Inputs         map[string]interface{}    `json:"inputs"`
	Query          string                    `json:"query"`
	ResponseMode   string                    `json:"response_mode,omitempty"`
	ConversationID string                    `json:"conversation_id,omitempty"`
	User           string                    `json:"user,omitempty"`
	Files          []ChatMessagesPayloadFile `json:"files"`
}

type ChatMessagesPayloadFile struct {
	Type           string `json:"type"`
	TransferMethod string `json:"transfer_method"`
	URL            string `json:"url"`
	UploadFileID   string `json:"upload_file_id"`
}

type ChatMessagesResponse struct {
	Event          string `json:"event"`
	MessageID      string `json:"message_id"`
	ConversationID string `json:"conversation_id"`
	TaskId         string `json:"task_id,omitempty"`
	ID             string `json:"id,omitempty"`
	Mode           string `json:"mode"`
	Answer         string `json:"answer"`
	Metadata       any    `json:"metadata"`
	CreatedAt      int    `json:"created_at"`
}

type ChatMessagesStopResponse struct {
	Result string `json:"result"`
}

type ErrorResponse struct {
	Code    string `json:"status"`
	Message string `json:"message"`
}

func PrepareChatPayload(payload map[string]interface{}) (string, error) {
	defer functrace.Trace([]interface {
	}{payload})()
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (dc *DifyClient) ChatMessages(query string, inputs map[string]interface{}, conversation_id string, files []ChatMessagesPayloadFile, appAuthorization string) (result ChatMessagesResponse, err error) {
	defer functrace.Trace([]interface {
	}{dc, query, inputs, conversation_id, files, appAuthorization})()
	var payload ChatMessagesPayload

	if inputs != nil {
		payload.Inputs = inputs
	} else {
		payload.Inputs = make(map[string]interface{})
	}
	if query == "" {
		return result, fmt.Errorf("query should be a valid JSON string")
	} else {
		payload.Query = query
	}

	payload.ResponseMode = RESPONSE_MODE_BLOCKING
	payload.User = dc.User

	if conversation_id != "" {
		payload.ConversationID = conversation_id
	}

	if len(files) > 0 {
		payload.Files = files
	}

	api := dc.GetAPI(API_CHAT_MESSAGES)

	code, body, err := SetAppsPostAuthorization(dc, api, payload, appAuthorization)

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

func (dc *DifyClient) ChatMessagesStreaming(query string, inputs map[string]interface{}, conversationID string, files []ChatMessagesPayloadFile, appAuthorization string) (<-chan string, error) {
	defer functrace.Trace([]interface {
	}{dc, query, inputs, conversationID, files, appAuthorization})()
	var payload ChatMessagesPayload

	if inputs != nil {
		payload.Inputs = inputs
	} else {
		payload.Inputs = make(map[string]interface{})
	}
	if query == "" {
		return nil, fmt.Errorf("query should be a valid JSON string")
	} else {
		payload.Query = query
	}

	payload.ResponseMode = "streaming"

	if conversationID != "" {
		payload.ConversationID = conversationID
	}

	if len(files) > 0 {
		payload.Files = files
	}

	fmt.Printf("Response Status:", payload.Files)

	api := dc.GetAPI(API_CHAT_MESSAGES)

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", api, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", appAuthorization)

	client := &http.Client{
		Timeout: time.Minute,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	fmt.Println("Response Status:", resp.StatusCode, resp)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
		}

		var errorResponse ErrorResponse
		if err := json.Unmarshal(bodyBytes, &errorResponse); err != nil {
			fmt.Println("Error parsing response body:", err)
			return nil, fmt.Errorf("server error, please try again later response，code: %d", 400)
		}

		if errorResponse.Code == "" {
			errorResponse.Code = "400"
		}

		if errorResponse.Message == "" {
			errorResponse.Message = "server error, please try again later response，code: " + fmt.Sprintf("%d", errorResponse.Code)
		}

		return nil, fmt.Errorf("message: %s, code: %s", errorResponse.Message, errorResponse.Code)
	}

	dataChan := make(chan string)
	go func() {
		defer close(dataChan)
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.TrimSpace(line) == "" {
				continue
			} else {
				jsonStr := strings.TrimPrefix(line, "data: ")

				var message ChatMessagesResponse
				err := json.Unmarshal([]byte(jsonStr), &message)
				if err != nil {
					fmt.Println("Error decoding JSON:", err)

					errorMessage, _ := json.Marshal(map[string]string{
						"error": fmt.Sprintf("error decoding JSON: %v", err),
					})
					dataChan <- string(errorMessage)
					continue
				}

				messageJson, err := json.Marshal(message)
				if err != nil {
					fmt.Println("Error encoding JSON:", err)

					errorMessage, _ := json.Marshal(map[string]string{
						"error": fmt.Sprintf("error encoding JSON: %v", err),
					})
					dataChan <- string(errorMessage)
					continue
				}
				dataChan <- string(messageJson)
			}

		}

		if err := scanner.Err(); err != nil {
			errorMessage, _ := json.Marshal(map[string]string{
				"error": fmt.Sprintf("scanner error: %v", err),
			})
			dataChan <- string(errorMessage)
		}
	}()

	return dataChan, nil
}

func (dc *DifyClient) ChatMessagesStop(task_id string, Authorization string) (result ChatMessagesStopResponse, err error) {
	defer functrace.Trace([]interface {
	}{dc, task_id, Authorization})()
	if task_id == "" {
		return result, fmt.Errorf("task_id is required")
	}

	if Authorization == "" {
		return result, fmt.Errorf("Authorization is required")
	}

	api := dc.GetAPI(API_CHAT_MESSAGES_STOP)
	api = UpdateAPIParam(api, API_PARAM_TASK_ID, task_id)

	payload := map[string]string{}

	code, body, err := SetAppsPostAuthorization(dc, api, payload, Authorization)
	if err != nil {
		return result, fmt.Errorf("failed to send request: %v", err)
	}

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
