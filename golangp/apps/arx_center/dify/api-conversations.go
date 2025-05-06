package dify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/toheart/functrace"
	"io"
	"net/http"
)

type ConversationsResponse struct {
	Limit   int  `json:"limit"`
	HasMore bool `json:"has_more"`
	Data    []struct {
		ID     string `json:"id"`
		Name   string `json:"name,omitempty"`
		Inputs struct {
			Book   string `json:"book"`
			MyName string `json:"myName"`
		} `json:"inputs,omitempty"`
		Status    string `json:"status,omitempty"`
		CreatedAt int    `json:"created_at,omitempty"`
	} `json:"data"`
}

type RenameConversationsResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Inputs struct {
	} `json:"inputs"`
	Status       string `json:"status"`
	Introduction string `json:"introduction"`
	CreatedAt    int    `json:"created_at"`
}

func (dc *DifyClient) GetConversations(Authorization string, limit int) (result ConversationsResponse, err error) {
	defer functrace.Trace([]interface {
	}{dc, Authorization, limit})()
	if limit == 0 {
		limit = 100
	}

	api := dc.GetAPI(API_CONVERSATIONS)
	api = api + "?limit=" + fmt.Sprintf("%d", limit) + "&pinned=false"

	code, body, err := SetGetAppsAuthorization(dc, api, Authorization)

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

type DeleteConversationsResponse struct {
	Result string `json:"result"`
}

func (dc *DifyClient) DeleteConversations(conversation_id string, Authorization string) (err error) {
	defer functrace.Trace([]interface {
	}{dc, conversation_id, Authorization})()
	if conversation_id == "" {
		return fmt.Errorf("conversation_id is required")
	}

	payloadBody := map[string]string{
		"user": dc.User,
	}

	api := dc.GetAPI(API_CONVERSATIONS_DELETE)
	api = UpdateAPIParam(api, API_PARAM_CONVERSATION_ID, conversation_id)

	buf, err := json.Marshal(payloadBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("DELETE", api, bytes.NewBuffer(buf))

	if err != nil {
		return fmt.Errorf("could not create a new request: %v", err)
	}
	req.Header.Set("Authorization", Authorization)
	req.Header.Set("Content-Type", "application/json")

	resp, err := dc.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode <= 200 || resp.StatusCode >= 300 {
		bodyText, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("status code: %d, could not read the body", resp.StatusCode)
		}
		return fmt.Errorf("status code: %d, %s", resp.StatusCode, bodyText)
	}

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("status code: %d, could not read the body %v", resp.StatusCode, err)
		return err
	}
	fmt.Println(string(bodyText))
	return nil
}

func (dc *DifyClient) RenameConversations(conversation_id string, name string, Authorization string) (result RenameConversationsResponse, err error) {
	defer functrace.Trace([]interface {
	}{dc, conversation_id, name, Authorization})()
	if conversation_id == "" {
		return result, fmt.Errorf("conversation_id is required")
	}

	payload := map[string]string{
		"name": name,
	}

	api := dc.GetAPI(API_CONVERSATIONS_RENAME)
	api = UpdateAPIParam(api, API_PARAM_CONVERSATION_ID, conversation_id)

	code, body, err := SetAppsPostAuthorization(dc, api, payload, Authorization)

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
