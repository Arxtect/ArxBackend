package dify

import (
	"fmt"
	"github.com/toheart/functrace"
	"net/http"
)

func setConsoleAuthorization(dc *DifyClient, req *http.Request) {
	defer functrace.Trace([]interface {
	}{dc, req})()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", dc.ConsoleToken))
	req.Header.Set("Content-Type", "application/json")
}

func SendGetRequestToConsole(dc *DifyClient, api string) (httpCode int, bodyText []byte, err error) {
	defer functrace.Trace([]interface {
	}{dc, api})()
	return SendGetRequest(true, dc, api, nil)
}

func SendPostRequestToConsole(dc *DifyClient, api string, postBody interface{}) (httpCode int, bodyText []byte, err error) {
	defer functrace.Trace([]interface {
	}{dc, api, postBody})()
	return SendPostRequest(true, dc, api, postBody, nil)
}

func SendPutRequestToConsole(dc *DifyClient, api string, putBody interface{}) (httpCode int, bodyText []byte, err error) {
	defer functrace.Trace([]interface {
	}{dc, api, putBody})()
	return SendPutRequest(true, dc, api, putBody, nil)
}

func SendDeleteRequestToConsole(dc *DifyClient, api string) (httpCode int, bodyText []byte, err error) {
	defer functrace.Trace([]interface {
	}{dc, api})()
	return SendDeleteRequest(true, dc, api)
}
