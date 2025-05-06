package dify

import (
	"encoding/json"
	"fmt"
	"github.com/toheart/functrace"
	"strings"
)

type TagsModel struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type TagsResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	BindingCount string `json:"binding_count"`
}

type TagsBindingPayload struct {
	TagIds   []string `json:"tag_ids"`
	TargetID string   `json:"target_id"`
	Type     string   `json:"type"`
}

func (dc *DifyClient) GetTagsList() (result []TagsResponse, err error) {
	defer functrace.Trace([]interface {
	}{dc})()
	api := dc.GetConsoleAPI(CONSOLE_API_APPS_TAGS_GET)

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

func (dc *DifyClient) CreateTag(name string) (result TagsResponse, err error) {
	defer functrace.Trace([]interface {
	}{dc, name})()
	payload := TagsModel{
		Name: name,
		Type: "app",
	}

	api := dc.GetConsoleAPI(CONSOLE_API_APPS_TAGS_CREATE)
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

func (dc *DifyClient) InitTag() {
	defer functrace.Trace([]interface {
	}{dc})()

	existingTags, err := dc.GetTagsList()
	if err != nil {
		fmt.Println("error fetching tags:", err)
		return
	}

	existingTagNames := make(map[string]bool)
	for _, tag := range existingTags {
		existingTagNames[tag.Name] = true
	}

	var tagList = []string{"assistant", "knowledge"}

	for _, tag := range tagList {

		if existingTagNames[tag] {
			fmt.Println("tag already exists:", tag)
			continue
		}

		result, err := dc.CreateTag(tag)
		if err != nil {
			fmt.Println("error creating tag:", err)
		} else {
			fmt.Println("created tag:", result)
		}
	}
}

func (dc *DifyClient) HandleTagsBinding(appID string, appType string) (err error) {
	defer functrace.Trace([]interface {
	}{dc, appID, appType})()

	existingTags, err := dc.GetTagsList()
	if err != nil {
		return fmt.Errorf("error fetching tags: %v", err)
	}

	var tagID string
	for _, tag := range existingTags {
		if strings.EqualFold(tag.Name, appType) {
			tagID = tag.ID
			break
		}
	}
	if tagID == "" {
		return fmt.Errorf("tag not found for appType: %s", appType)
	}

	payload := TagsBindingPayload{
		TagIds:   []string{tagID},
		TargetID: appID,
		Type:     "app",
	}

	api := dc.GetConsoleAPI(CONSOLE_API_APPS_TAGS_BINDINGS)
	code, body, err := SendPostRequestToConsole(dc, api, payload)
	if err != nil {
		return fmt.Errorf("error sending bind request: %v", err)
	}

	err = CommonRiskForSendRequest(code, err)
	if err != nil {
		fmt.Println("error: ", string(body))
		return err
	}

	if code != 200 {
		return fmt.Errorf("failed to bind tags, received status code: %d", code)
	}

	fmt.Println("tags successfully bound")
	return nil
}
