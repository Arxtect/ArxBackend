package dify

import (
	"encoding/json"
	"fmt"
	"github.com/toheart/functrace"
)

type MessagesResponse struct {
	Limit   int  `json:"limit"`
	HasMore bool `json:"has_more"`
	Data    []struct {
		ID             string `json:"id"`
		ConversationID string `json:"conversation_id"`
		Inputs         struct {
			Name string `json:"name"`
		} `json:"inputs"`
		Query              string `json:"query"`
		Answer             string `json:"answer"`
		MessageFiles       []any  `json:"message_files"`
		Feedback           any    `json:"feedback"`
		RetrieverResources []struct {
			Position     int     `json:"position"`
			DatasetID    string  `json:"dataset_id"`
			DatasetName  string  `json:"dataset_name"`
			DocumentID   string  `json:"document_id"`
			DocumentName string  `json:"document_name"`
			SegmentID    string  `json:"segment_id"`
			Score        float64 `json:"score"`
			Content      string  `json:"content"`
		} `json:"retriever_resources"`
		AgentThoughts []any `json:"agent_thoughts"`
		CreatedAt     int   `json:"created_at"`
	} `json:"data"`
}

func (dc *DifyClient) GetConversationsMessages(Authorization string, lastId string, conversationId string, limit int) (result MessagesResponse, err error) {
	defer functrace.Trace([]interface {
	}{dc, Authorization, lastId, conversationId, limit})()
	if conversationId == "" {
		return result, fmt.Errorf("conversation_id is required")
	}

	if limit == 0 {
		limit = 20
	}

	api := dc.GetAPI(API_MESSAGES)
	api = api + "?limit=" + fmt.Sprintf("%d", limit) + "&last_id=" + lastId + "&conversation_id=" + conversationId

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
