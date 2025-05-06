package dify

import (
	"encoding/json"
	"fmt"
	"github.com/toheart/functrace"
)

type PassportResponse struct {
	AccessToken string `json:"access_token"`
}

func (dc *DifyClient) GetAppsAccessToken(appsAccessToken string) (result PassportResponse, err error) {
	defer functrace.Trace([]interface {
	}{dc, appsAccessToken})()
	api := dc.GetAPI(API_PASSPORT)

	header := map[string]string{"X-App-Code": appsAccessToken}

	code, body, err := SendGetRequestToAPI(dc, api, header)

	fmt.Println(code, api, err)

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

func SetGetAppsAuthorization(dc *DifyClient, api string, Authorization string) (code int, body []byte, err error) {
	defer functrace.Trace([]interface {
	}{dc, api, Authorization})()
	header := map[string]string{"Authorization": Authorization}

	code, body, err = SendGetRequestToAPI(dc, api, header)

	return code, body, err
}

func SetAppsPostAuthorization(dc *DifyClient, api string, postBody interface{}, Authorization string) (code int, body []byte, err error) {
	defer functrace.Trace([]interface {
	}{dc, api, postBody, Authorization})()
	header := map[string]string{"Authorization": Authorization}

	code, body, err = SendPostRequestToAPI(dc, api, postBody, header)

	return code, body, err
}

func SetAppsPutAuthorization(dc *DifyClient, api string, postBody interface{}, Authorization string) (code int, body []byte, err error) {
	defer functrace.Trace([]interface {
	}{dc, api, postBody, Authorization})()
	header := map[string]string{"Authorization": Authorization}

	code, body, err = SendPutRequestToAPI(dc, api, postBody, header)

	return code, body, err
}
