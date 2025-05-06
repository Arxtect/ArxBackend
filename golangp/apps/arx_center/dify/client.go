package dify

import (
	"github.com/arxtect/ArxBackend/golangp/common/logger"
	"github.com/arxtect/ArxBackend/golangp/config"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/toheart/functrace"
)

type DifyClientConfig struct {
	Key         string
	Host        string
	HostUrl     string
	ConsoleHost string
	Timeout     int
	SkipTLS     bool
	User        string
}

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

func CreateDifyClient(difyConfig DifyClientConfig) (*DifyClient, error) {
	defer functrace.Trace([]interface {
	}{difyConfig})()
	cnf := config.Env
	key := cnf.DifyKey
	if key == "" {
		return nil, fmt.Errorf("dify API Key is required")
	}

	host := cnf.DifyHost
	if host == "" {
		return nil, fmt.Errorf("dify Host is required")
	}

	consoleURL := host + "/console/api"

	timeout := 0 * time.Second
	if difyConfig.Timeout <= 0 {
		if difyConfig.Timeout < 0 {
			fmt.Println("Timeout should be a positive number, reset to default value: 10s")
		}
		timeout = DEFAULT_TIMEOUT * time.Second
	}

	skipTLS := false
	if difyConfig.SkipTLS {
		skipTLS = true
	}

	difyConfig.User = strings.TrimSpace(difyConfig.User)
	if difyConfig.User == "" {
		difyConfig.User = DEFAULT_USER
	}

	var client *http.Client

	if skipTLS {
		client = &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}
	} else {
		client = &http.Client{}
	}

	if timeout > 0 {
		client.Timeout = timeout
	}

	return &DifyClient{
		Key:         key,
		Host:        host,
		HostUrl:     host + "/api",
		ConsoleHost: consoleURL,
		Timeout:     timeout,
		SkipTLS:     skipTLS,
		Client:      client,
		User:        difyConfig.User,
	}, nil
}

func GetDifyClient() (*DifyClient, error) {
	defer functrace.Trace([]interface {
	}{})()
	client, err := CreateDifyClient(DifyClientConfig{})
	if err != nil {
		logger.Warning("failed to create DifyClient: %v\n", err)
		return nil, err
	}

	fmt.Println(client)

	return client, nil
}
