package dify

import (
	"encoding/json"
	"fmt"
	"github.com/toheart/functrace"
	"io"
	"net/http"
)

type SupportTypeResponse struct {
	AllowedExtensions []string `json:"allowed_extensions"`
}

type FileUploadPayload struct {
	Image ImagePayload `json:"image"`
}

type ImagePayload struct {
	Enabled         bool     `json:"enabled"`
	NumberLimits    int      `json:"number_limits"`
	Detail          string   `json:"detail"`
	TransferMethods []string `json:"transfer_methods"`
}

func (dc *DifyClient) GetSupportTypes() (result SupportTypeResponse, err error) {
	defer functrace.Trace([]interface {
	}{dc})()
	api := dc.GetConsoleAPI(CONSOLE_API_SUPPORT_TYPES)

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

func (dc *DifyClient) PreviewFile(file_id string) (body []byte, contentType string, err error) {
	defer functrace.Trace([]interface {
	}{dc, file_id})()
	if file_id == "" {
		return nil, "", fmt.Errorf("file_id is required")
	}
	api := dc.GetConsoleAPI(CONSOLE_API_FILE_PREVIEW)
	api = UpdateAPIParam(api, CONSOLE_API_PARAM_FILE_ID, file_id)

	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return nil, "", err
	}
	setConsoleAuthorization(dc, req)

	resp, err := dc.Client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	code := resp.StatusCode
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	err = CommonRiskForSendRequest(code, err)
	if err != nil {
		return nil, "", err
	}

	contentType = resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return body, contentType, nil
}
