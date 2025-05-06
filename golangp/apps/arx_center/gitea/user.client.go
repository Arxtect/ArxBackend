package gitea

import (
	"github.com/Arxtect/ArxBackend/golangp/config"

	"code.gitea.io/sdk/gitea"
	"github.com/toheart/functrace"
)

type GiteaUserClient struct {
	Host     string
	User     string
	Password string
	Client   *gitea.Client
}

func CreateGiteaUserClient(username string, password string) (*GiteaUserClient, error) {
	defer functrace.Trace([]interface {
	}{username, password})()
	configCopy := config.Env

	client, err := gitea.NewClient(configCopy.GiteaHost, gitea.SetBasicAuth(username, password))
	if err != nil {
		return nil, err
	}

	return &GiteaUserClient{
		Host:     configCopy.GiteaHost,
		User:     username,
		Password: password,
		Client:   client,
	}, nil
}
