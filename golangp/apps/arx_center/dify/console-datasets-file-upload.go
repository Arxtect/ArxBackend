package dify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/toheart/functrace"
	"io"
	"mime/multipart"
	"net/http"
)

type FileUploadResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Size      int    `json:"size"`
	Extension string `json:"extension"`
	MimeType  string `json:"mime_type"`
	CreatedBy string `json:"created_by"`
	CreatedAt int    `json:"created_at"`
}

type DatasetFileResponse struct {
	ID             string `json:"id,omitempty"`
	Position       int    `json:"position,omitempty"`
	DataSourceType string `json:"data_source_type,omitempty"`
	DataSourceInfo struct {
		UploadFileID string `json:"upload_file_id,omitempty"`
	} `json:"data_source_info,omitempty"`
	DataSourceDetailDict struct {
		UploadFile struct {
			ID        string  `json:"id,omitempty"`
			Name      string  `json:"name,omitempty"`
			Size      int     `json:"size,omitempty"`
			Extension string  `json:"extension,omitempty"`
			MimeType  string  `json:"mime_type,omitempty"`
			CreatedBy string  `json:"created_by,omitempty"`
			CreatedAt float64 `json:"created_at,omitempty"`
		} `json:"upload_file,omitempty"`
	} `json:"data_source_detail_dict,omitempty"`
	DatasetProcessRuleID string `json:"dataset_process_rule_id,omitempty"`
	Name                 string `json:"name,omitempty"`
	CreatedFrom          string `json:"created_from,omitempty"`
	CreatedBy            string `json:"created_by,omitempty"`
	CreatedAt            int    `json:"created_at,omitempty"`
	Tokens               int    `json:"tokens,omitempty"`
	IndexingStatus       string `json:"indexing_status,omitempty"`
	Error                any    `json:"error,omitempty"`
	Enabled              bool   `json:"enabled,omitempty"`
	DisabledAt           any    `json:"disabled_at,omitempty"`
	DisabledBy           any    `json:"disabled_by,omitempty"`
	Archived             bool   `json:"archived,omitempty"`
	DisplayStatus        string `json:"display_status,omitempty"`
	WordCount            int    `json:"word_count,omitempty"`
	HitCount             int    `json:"hit_count,omitempty"`
	DocForm              string `json:"doc_form,omitempty"`
}

type GetDatasetFileListResponse struct {
	Data    []DatasetFileResponse `json:"data,omitempty"`
	HasMore bool                  `json:"has_more,omitempty"`
	Limit   int                   `json:"limit,omitempty"`
	Total   int                   `json:"total,omitempty"`
	Page    int                   `json:"page,omitempty"`
}

func (dc *DifyClient) DatasetsFileUpload(file multipart.File, fileName string) (result FileUploadResponse, err error) {
	defer functrace.Trace([]interface {
	}{dc, file, fileName})()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return result, fmt.Errorf("error creating form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return result, fmt.Errorf("error copying file: %v", err)
	}

	_ = writer.WriteField("user", dc.User)
	err = writer.Close()
	if err != nil {
		return result, fmt.Errorf("error closing writer: %v", err)
	}

	req, err := http.NewRequest("POST", dc.GetConsoleAPI(CONSOLE_API_FILE_UPLOAD_DATASETS), body)
	if err != nil {
		return result, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", dc.ConsoleToken))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return result, fmt.Errorf("status code: %d, create file failed", resp.StatusCode)
	}

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("could not read the body: %v", err)
	}

	err = json.Unmarshal(bodyText, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal the response: %v", err)
	}
	return result, nil
}

func (dc *DifyClient) DatasetsFileUploadChat(file io.Reader, fileName string, appAuthorization string) (result FileUploadResponse, err error) {
	defer functrace.Trace([]interface {
	}{dc, file, fileName, appAuthorization})()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return result, fmt.Errorf("error creating form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return result, fmt.Errorf("error copying file: %v", err)
	}

	err = writer.WriteField("user", dc.User)
	if err != nil {
		return result, fmt.Errorf("error writing user field: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return result, fmt.Errorf("error closing writer: %v", err)
	}

	req, err := http.NewRequest("POST", dc.GetAPI(CONSOLE_API_FILE_UPLOAD), body)
	if err != nil {
		return result, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", appAuthorization)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return result, fmt.Errorf("status code: %d, create file failed", resp.StatusCode)
	}

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("could not read the body: %v", err)
	}

	err = json.Unmarshal(bodyText, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal the response: %v", err)
	}
	return result, nil
}

func (dc *DifyClient) GetDatasetsFileList(datasets_id string, limit int, page int) (result GetDatasetFileListResponse, err error) {
	defer functrace.Trace([]interface {
	}{dc, datasets_id, limit, page})()
	if limit == 0 {
		limit = 500
	}

	if page == 0 {
		page = 1
	}

	api := dc.GetConsoleAPI(CONSOLE_API_DATASETS_UPDATE_DATASETS)
	api = UpdateAPIParam(api, CONSOLE_API_PARAM_DATASETS_ID, datasets_id)
	api = api + "?limit=" + fmt.Sprintf("%d", limit) + "&page=" + fmt.Sprintf("%d", page) + "&keyword=&fetch="

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

func (dc *DifyClient) DeleteDocumentForDatasets(datasets_id string, document_id string) (ok bool, err error) {
	defer functrace.Trace([]interface {
	}{dc, datasets_id, document_id})()
	if datasets_id == "" || document_id == "" {
		return false, fmt.Errorf("datasets_id or document_id is required")
	}

	api := dc.GetConsoleAPI(CONSOLE_API_DATASETS_DELETE_FILE)
	api = UpdateAPIParam(api, CONSOLE_API_PARAM_DATASETS_ID, datasets_id)
	api = UpdateAPIParam(api, CONSOLE_API_PARAM_DOCUMENT_ID, document_id)

	fmt.Println(api)

	req, err := http.NewRequest("DELETE", api, nil)
	if err != nil {
		return false, fmt.Errorf("could not create a new request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", dc.ConsoleToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := dc.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		bodyText, err := io.ReadAll(resp.Body)
		if err != nil {
			return false, fmt.Errorf("status code: %d, could not read the body", resp.StatusCode)
		}
		return false, fmt.Errorf("status code: %d, %s", resp.StatusCode, bodyText)
	}

	return true, nil
}

func (dc *DifyClient) RenameDocumentForDatasets(datasets_id string, document_id string, name string) (result DatasetFileResponse, err error) {
	defer functrace.Trace([]interface {
	}{dc, datasets_id, document_id, name})()
	if datasets_id == "" || document_id == "" {
		return result, fmt.Errorf("datasets_id or document_id is required")
	}

	payload := map[string]string{
		"name": name,
	}

	api := dc.GetConsoleAPI(CONSOLE_API_DATASETS_RENAME_FILE)
	api = UpdateAPIParam(api, CONSOLE_API_PARAM_DATASETS_ID, datasets_id)
	api = UpdateAPIParam(api, CONSOLE_API_PARAM_DOCUMENT_ID, document_id)

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
