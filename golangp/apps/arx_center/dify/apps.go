package dify

import (
	"encoding/json"
	"fmt"
	"github.com/toheart/functrace"
)

type CreateAppsPayload struct {
	Name           string `json:"name,omitempty"`
	Icon           string `json:"icon,omitempty"`
	IconBackground string `json:"icon_background,omitempty"`
	Mode           string `json:"mode,omitempty"`
	Description    string `json:"description"`
}

type FileUpload struct {
	Image struct {
		Enabled         bool     `json:"enabled"`
		NumberLimits    int      `json:"number_limits"`
		Detail          string   `json:"detail"`
		TransferMethods []string `json:"transfer_methods"`
	} `json:"image"`
}

type CreateAppsResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Mode           string `json:"mode"`
	Icon           string `json:"icon"`
	IconBackground string `json:"icon_background"`
	EnableSite     bool   `json:"enable_site"`
	EnableAPI      bool   `json:"enable_api"`
	ModelConfig    struct {
		OpeningStatement              any   `json:"opening_statement"`
		SuggestedQuestions            []any `json:"suggested_questions"`
		SuggestedQuestionsAfterAnswer struct {
			Enabled bool `json:"enabled"`
		} `json:"suggested_questions_after_answer"`
		SpeechToText struct {
			Enabled bool `json:"enabled"`
		} `json:"speech_to_text"`
		TextToSpeech struct {
			Enabled bool `json:"enabled"`
		} `json:"text_to_speech"`
		RetrieverResource struct {
			Enabled bool `json:"enabled"`
		} `json:"retriever_resource"`
		AnnotationReply struct {
			Enabled bool `json:"enabled"`
		} `json:"annotation_reply"`
		MoreLikeThis struct {
			Enabled bool `json:"enabled"`
		} `json:"more_like_this"`
		SensitiveWordAvoidance struct {
			Enabled bool   `json:"enabled"`
			Type    string `json:"type"`
			Configs []any  `json:"configs"`
		} `json:"sensitive_word_avoidance"`
		ExternalDataTools []any `json:"external_data_tools"`
		Model             struct {
			Provider         string `json:"provider"`
			Name             string `json:"name"`
			Mode             string `json:"mode"`
			CompletionParams struct {
			} `json:"completion_params"`
		} `json:"model"`
		UserInputForm        []any `json:"user_input_form"`
		DatasetQueryVariable any   `json:"dataset_query_variable"`
		PrePrompt            any   `json:"pre_prompt"`
		AgentMode            struct {
			Enabled  bool  `json:"enabled"`
			Strategy any   `json:"strategy"`
			Tools    []any `json:"tools"`
			Prompt   any   `json:"prompt"`
		} `json:"agent_mode"`
		PromptType       string `json:"prompt_type"`
		ChatPromptConfig struct {
		} `json:"chat_prompt_config"`
		CompletionPromptConfig struct {
		} `json:"completion_prompt_config"`
		DatasetConfigs struct {
			RetrievalModel string `json:"retrieval_model"`
		} `json:"dataset_configs"`
		FileUpload FileUpload `json:"file_upload"`
		CreatedAt  int        `json:"created_at"`
	} `json:"model_config"`
	Tracing   any `json:"tracing"`
	CreatedAt int `json:"created_at"`
}

type GetAppsResponse struct {
	Page    int  `json:"page"`
	Limit   int  `json:"limit"`
	Total   int  `json:"total"`
	HasMore bool `json:"has_more"`
	Data    []struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		Description    string `json:"description"`
		Mode           string `json:"mode"`
		Icon           string `json:"icon"`
		IconBackground string `json:"icon_background"`
		ModelConfig    struct {
			Model struct {
				Provider         string `json:"provider"`
				Name             string `json:"name"`
				Mode             string `json:"mode"`
				CompletionParams struct {
				} `json:"completion_params"`
			} `json:"model"`
			PrePrompt any `json:"pre_prompt"`
		} `json:"model_config"`
		CreatedAt int   `json:"created_at"`
		Tags      []any `json:"tags"`
	} `json:"data"`
}

type DatasetArray struct {
	Enabled bool   `json:"enabled"`
	ID      string `json:"id"`
}

type DatasetConfigs struct {
	RetrievalModel string `json:"retrieval_model"`
	Datasets       struct {
		Datasets []struct {
			Dataset DatasetArray `json:"dataset"`
		} `json:"datasets"`
	} `json:"datasets"`
}

type UpdateAppsModelConfigPayload struct {
	PrePrompt        string `json:"pre_prompt"`
	PromptType       string `json:"prompt_type"`
	ChatPromptConfig struct {
	} `json:"chat_prompt_config"`
	CompletionPromptConfig struct {
	} `json:"completion_prompt_config"`
	UserInputForm        []any  `json:"user_input_form"`
	DatasetQueryVariable string `json:"dataset_query_variable"`
	OpeningStatement     string `json:"opening_statement"`
	SuggestedQuestions   []any  `json:"suggested_questions"`
	MoreLikeThis         struct {
		Enabled bool `json:"enabled"`
	} `json:"more_like_this"`
	SuggestedQuestionsAfterAnswer struct {
		Enabled bool `json:"enabled"`
	} `json:"suggested_questions_after_answer"`
	SpeechToText struct {
		Enabled bool `json:"enabled"`
	} `json:"speech_to_text"`
	TextToSpeech struct {
		Enabled  bool   `json:"enabled"`
		Voice    string `json:"voice"`
		Language string `json:"language"`
	} `json:"text_to_speech"`
	RetrieverResource struct {
		Enabled bool `json:"enabled"`
	} `json:"retriever_resource"`
	SensitiveWordAvoidance struct {
		Enabled bool   `json:"enabled"`
		Type    string `json:"type"`
		Configs []any  `json:"configs"`
	} `json:"sensitive_word_avoidance"`
	AgentMode struct {
		Enabled      bool   `json:"enabled"`
		MaxIteration int    `json:"max_iteration"`
		Strategy     string `json:"strategy"`
		Tools        []any  `json:"tools"`
	} `json:"agent_mode"`
	Model          ListDatasetsModelConfigDetail `json:"model"`
	DatasetConfigs DatasetConfigs                `json:"dataset_configs"`
	FileUpload     FileUploadPayload             `json:"file_upload"`
}

type SiteStruct struct {
	AccessToken    string `json:"access_token"`
	Code           string `json:"code"`
	Title          string `json:"title"`
	Icon           string `json:"icon"`
	IconBackground string `json:"icon_background"`

	DefaultLanguage        string `json:"default_language"`
	ChatColorTheme         any    `json:"chat_color_theme"`
	ChatColorThemeInverted bool   `json:"chat_color_theme_inverted"`
	CustomizeDomain        any    `json:"customize_domain"`
	Copyright              any    `json:"copyright"`
	PrivacyPolicy          any    `json:"privacy_policy"`
	CustomDisclaimer       any    `json:"custom_disclaimer"`
	CustomizeTokenStrategy string `json:"customize_token_strategy"`
	PromptPublic           bool   `json:"prompt_public"`
	AppBaseURL             string `json:"app_base_url"`
	ShowWorkflowSteps      bool   `json:"show_workflow_steps"`
}

type GetAppsViaAppIdResponse struct {
	ID             string                       `json:"id"`
	Name           string                       `json:"name"`
	Description    string                       `json:"description"`
	Mode           string                       `json:"mode"`
	Icon           string                       `json:"icon"`
	IconBackground string                       `json:"icon_background"`
	EnableSite     bool                         `json:"enable_site"`
	EnableAPI      bool                         `json:"enable_api"`
	ModelConfig    UpdateAppsModelConfigPayload `json:"model_config"`
	APIBaseURL     string                       `json:"api_base_url"`
	CreatedAt      int                          `json:"created_at"`
	DeletedTools   []any                        `json:"deleted_tools"`
	Site           SiteStruct                   `json:"site"`
}

func (dc *DifyClient) CreateAppsHandler(name string, description string, model string, prePrompt string) (result GetAppsViaAppIdResponse, err error) {
	defer functrace.Trace([]interface {
	}{dc, name, description, model, prePrompt})()
	payload := CreateAppsPayload{
		Name:           name,
		Description:    description,
		Mode:           "chat",
		Icon:           "robot_face",
		IconBackground: "#E4FBCC",
	}
	api := dc.GetConsoleAPI(CONSOLE_API_APPS_CREATE)

	code, body, err := SendPostRequestToConsole(dc, api, payload)

	fmt.Println("code: ", code, "body: ", string(body), "err: ", err)

	err = CommonRiskForSendRequest(code, err)
	if err != nil {
		fmt.Println("error: ", string(body))
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal the response: %v", err)
	}

	updateResult, err := dc.UpdateAppsInfo(result.ID, name, description, []string{}, model, prePrompt)
	if err != nil {
		return result, err
	}

	return updateResult, nil
}

func (dc *DifyClient) UpdateAppsInfo(appId string, name string, description string, document_ids []string, modelName string, prePrompt string) (result GetAppsViaAppIdResponse, err error) {
	defer functrace.Trace([]interface {
	}{dc, appId, name, description, document_ids, modelName, prePrompt})()

	results, err := dc.GetAppsHandler(appId)
	fmt.Printf("results: %+v\n", results)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	payload := CreateAppsPayload{
		Name:           name,
		Description:    description,
		Icon:           "robot_face",
		IconBackground: "#E4FBCC",
	}

	api := dc.GetConsoleAPI(CONSOLE_API_APPS_GET)
	api = UpdateAPIParam(api, API_PARAM_APP_ID, appId)
	code, body, err := SendPutRequestToConsole(dc, api, payload)
	err = CommonRiskForSendRequest(code, err)
	if err != nil {
		fmt.Println("error: ", string(body))
		return result, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal the response: %v", err)
	}

	modelConfig := result.ModelConfig

	configResult, err := dc.UpdateAppModelConfig(appId, document_ids, modelName, prePrompt, modelConfig)

	if err != nil {
		fmt.Println("error: ", err, "configResult: ", configResult)
		return
	}

	return result, nil
}

type UpdateModelResponse struct {
	Result string `json:"result"`
}

func (dc *DifyClient) UpdateAppModelConfig(appId string, document_ids []string, modelName string, prePrompt string, modelConfig UpdateAppsModelConfigPayload) (result UpdateModelResponse, err error) {
	defer functrace.Trace([]interface {
	}{dc, appId, document_ids, modelName, prePrompt, modelConfig})()

	datasetList := make([]struct {
		Dataset DatasetArray `json:"dataset"`
	}, len(document_ids))

	for i, id := range document_ids {
		datasetList[i] = struct {
			Dataset DatasetArray `json:"dataset"`
		}{
			Dataset: DatasetArray{
				Enabled: true,
				ID:      id,
			},
		}
	}

	modelConfig.DatasetConfigs.Datasets.Datasets = datasetList
	modelConfig.PrePrompt = prePrompt
	modelConfig.Model.Name = modelName
	modelConfig.FileUpload.Image.Enabled = true
	modelConfig.FileUpload.Image.TransferMethods = []string{"remote_url", "local_file"}

	payload := modelConfig

	fmt.Printf("payload: %+v\n", payload, modelName)

	api := dc.GetConsoleAPI(CONSOLE_APU_APPS_UPDATE_MODEL_CONFIG)
	api = UpdateAPIParam(api, API_PARAM_APP_ID, appId)

	fmt.Println("api: ", api)

	code, body, err := SendPostRequestToConsole(dc, api, payload)
	err = CommonRiskForSendRequest(code, err)
	if err != nil {
		fmt.Println("error: ", string(body))
		return result, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal the response: %v", err)
	}
	return result, nil
}

func (dc *DifyClient) GetAppsHandler(appId string) (result GetAppsViaAppIdResponse, err error) {
	defer functrace.Trace([]interface {
	}{dc, appId})()
	api := dc.GetConsoleAPI(CONSOLE_API_APPS_GET)
	api = UpdateAPIParam(api, API_PARAM_APP_ID, appId)
	code, body, err := SendGetRequestToConsole(dc, api)
	err = CommonRiskForSendRequest(code, err)
	if err != nil {
		fmt.Println("error: ", string(body))
		return result, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal the response: %v", err)
	}
	return result, nil
}

func (dc *DifyClient) DeleteAppsHandler(appId string) (err error) {
	defer functrace.Trace([]interface {
	}{dc, appId})()
	api := dc.GetConsoleAPI(CONSOLE_API_APPS_GET)
	api = UpdateAPIParam(api, API_PARAM_APP_ID, appId)
	code, body, err := SendDeleteRequestToConsole(dc, api)
	err = CommonRiskForSendRequest(code, err)
	if err != nil {
		fmt.Println("error: ", string(body))
		return err
	}
	return nil
}
