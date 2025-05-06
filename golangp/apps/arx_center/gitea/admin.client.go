package gitea

import (
	"github.com/Arxtect/ArxBackend/golangp/config"
	"fmt"

	"code.gitea.io/sdk/gitea"
	"github.com/toheart/functrace"
)

type GiteaAdminClient struct {
	Host     string
	User     string
	Password string
	Client   *gitea.Client
}

func CreateGiteaAdminClient() (*GiteaAdminClient, error) {
	defer functrace.Trace([]interface {
	}{})()
	configCopy := config.Env

	client, err := gitea.NewClient(configCopy.GiteaHost, gitea.SetBasicAuth(configCopy.GiteaAdminUser, configCopy.GiteaAdminPassword))
	if err != nil {
		return nil, fmt.Errorf("failed to create admin client: %v", err)
	}

	return &GiteaAdminClient{Host: configCopy.GiteaHost, User: configCopy.GiteaAdminUser, Password: configCopy.GiteaAdminPassword, Client: client}, nil
}

func GetGiteaAdminClient() (*gitea.Client, error) {
	defer functrace.Trace([]interface {
	}{})()
	adminClient, err := CreateGiteaAdminClient()
	if err != nil {
		return nil, err
	}

	return adminClient.Client, nil
}
