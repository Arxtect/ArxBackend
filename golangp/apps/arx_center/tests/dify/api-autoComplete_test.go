package dify

import (
	"net/http"
	"testing"

	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/dify"
)

/*
type DifyClient struct {
	Key          string
	Host         string
	HostUrl      string
	ConsoleHost  string
	ConsoleToken string
	Timeout      time.Duration
	SkipTLS      bool
	Client       *http.Client
	User         string
}
*/

func TestAutoComplete(t *testing.T) {
	payload := dify.AutoCompletePayload{}
	difyClient := &dify.DifyClient{
		Key:          "test",
		Host:         "http://network.jancsitech.net:1510/",
		HostUrl:      "http://network.jancsitech.net:1510/api",
		ConsoleHost:  "http://network.jancsitech.net:1510/console/api",
		ConsoleToken: "test",
		Timeout:      10,
		SkipTLS:      false,
		Client:       &http.Client{},
		User:         "test",
	}
	payload.ResponseMode = "blocking"
	payload.Inputs = map[string]interface{}{
		"filename": "resume.tex",
	}
	payload.AcessToken = "Bearer " + dify.API_AUTOCOMPLETE_ACCESSTOKEN
	result, err := difyClient.AutoComplete(&payload)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("AutoComplete Response Suggesion: %v", result.Data.Outputs.Suggestion)
}
