package dify

import (
	"github.com/toheart/functrace"
	"net/http"
)

func setAPIAuthorization(dc *DifyClient, req *http.Request) {
	defer functrace.Trace([]interface {
	}{dc, req})()

	req.Header.Set("Content-Type", "application/json")
}

func SendGetRequestToAPI(dc *DifyClient, api string, header map[string]string) (httpCode int, bodyText []byte, err error) {
	defer functrace.Trace([]interface {
	}{dc, api, header})()
	return SendGetRequest(false, dc, api, header)
}

func SendPostRequestToAPI(dc *DifyClient, api string, postBody interface{}, header map[string]string) (httpCode int, bodyText []byte, err error) {
	defer functrace.Trace([]interface {
	}{dc, api, postBody, header})()
	return SendPostRequest(false, dc, api, postBody, header)
}

func SendPutRequestToAPI(dc *DifyClient, api string, putBody interface{}, header map[string]string) (httpCode int, bodyText []byte, err error) {
	defer functrace.Trace([]interface {
	}{dc, api, putBody, header})()
	return SendPutRequest(false, dc, api, putBody, header)
}
