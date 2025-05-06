package dify

import (
	"github.com/toheart/functrace"
	"strings"
)

const (
	API_COMPLETION_MESSAGES      = "/completion-messages"
	API_COMPLETION_MESSAGES_STOP = "/completion-messages/:task_id/stop"

	API_CHAT_MESSAGES      = "/chat-messages"
	API_CHAT_MESSAGES_STOP = "/chat-messages/:task_id/stop"
	API_AUTOCOMPLETE       = "/v1/workflows/run"

	API_MESSAGES           = "/messages"
	API_MESSAGES_SUGGESTED = "/messages/:message_id/suggested"
	API_MESSAGES_FEEDBACKS = "/messages/:message_id/feedbacks"

	API_CONVERSATIONS        = "/conversations"
	API_CONVERSATIONS_DELETE = "/conversations/:conversation_id"
	API_CONVERSATIONS_RENAME = "/conversations/:conversation_id/name"

	API_FILE_UPLOAD = "/files/upload"
	API_PARAMETERS  = "/parameters"
	API_META        = "/meta"

	API_AUDIO_TO_TEXT = "/audio-to-text"
	API_TEXT_TO_AUDIO = "/text-to-audio"

	API_PASSPORT                 = "/passport"
	API_AUTOCOMPLETE_ACCESSTOKEN = "app-JeubTEmCifRMISuNC5ZXqT7d"

	API_PARAM_TASK_ID             = ":task_id"
	API_PARAM_MESSAGE_ID          = ":message_id"
	API_PARAM_CONVERSATION_ID     = ":conversation_id"
	API_PARAM_APP_ID              = ":app_id"
	CONSOLE_API_PARAM_DOCUMENT_ID = ":document_id"
	CONSOLE_API_PARAM_DATASETS_ID = ":datasets_id"
	CONSOLE_API_PARAM_FILE_ID     = ":file_id"

	CONSOLE_API_FILE_UPLOAD_DATASETS = "/files/upload?source=datasets"
	CONSOLE_API_FILE_UPLOAD          = "/files/upload"
	CONSOLE_API_FILE_PREVIEW         = "/files/:file_id/preview"
	CONSOLE_API_LOGIN                = "/login"
	CONSOLE_API_SUPPORT_TYPES        = "/files/support-type"

	CONSOLE_API_DATASETS_CREATE          = "/datasets"
	CONSOLE_API_DATASETS_LIST            = "/datasets"
	CONSOLE_API_DATASETS_DELETE          = "/datasets/:datasets_id"
	CONSOLE_API_DATASETS_GET             = "/datasets/:datasets_id"
	CONSOLE_API_DATASETS_INIT            = "/datasets/init"
	CONSOLE_API_DATASETS_INIT_STATUS     = "/datasets/:datasets_id/indexing-status"
	CONSOLE_API_DATASETS_UPDATE_DATASETS = "/datasets/:datasets_id/documents"
	CONSOLE_API_DATASETS_DELETE_FILE     = "/datasets/:datasets_id/documents/:document_id"
	CONSOLE_API_DATASETS_RENAME_FILE     = "/datasets/:datasets_id/documents/:document_id/rename"
	CONSOLE_API_DOCUMENT_STATUS          = "/datasets/:datasets_id/documents/:document_id/indexing-status"

	CONSOLE_API_APPS_CREATE              = "/apps"
	CONSOLE_API_APPS_GET                 = "/apps/:app_id"
	CONSOLE_APU_APPS_UPDATE_MODEL_CONFIG = "/apps/:app_id/model-config"

	CONSOLE_API_APPS_TAGS_CREATE   = "/tags"
	CONSOLE_API_APPS_TAGS_GET      = "/tags?type=app"
	CONSOLE_API_APPS_TAGS_BINDINGS = "/tag-bindings/create"

	CONSOLE_API_WORKSPACE_CURRENT              = "/workspaces/current"
	CONSOLE_API_WORKSPACES_RERANK_MODEL        = "/workspaces/current/models/model-types/rerank"
	CONSOLE_API_CURRENT_WORKSPACE_RERANK_MODEL = "/workspaces/current/default-model?model_type=rerank"

	CONSOLE_API_CURRENT_WORKSPACE_LLM_MODEL = "/workspaces/current/default-model?model_type=llm"
	CONSOLE_API_WORKSPACES_LLM_MODEL        = "/workspaces/current/models/model-types/llm"
)

func (dc *DifyClient) GetAPI(api string) string {
	defer functrace.Trace([]interface {
	}{dc, api})()
	return dc.HostUrl + api
}

func (dc *DifyClient) GetConsoleAPI(api string) string {
	defer functrace.Trace([]interface {
	}{dc, api})()
	return dc.ConsoleHost + api
}

func UpdateAPIParam(api, key, value string) string {
	defer functrace.Trace([]interface {
	}{api, key, value})()
	return strings.ReplaceAll(api, key, value)
}
