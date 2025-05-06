package gitea

import (
	"github.com/arxtect/ArxBackend/golangp/common/constants"
	"fmt"
	"time"

	"code.gitea.io/sdk/gitea"
	"github.com/toheart/functrace"
)

type CreateRepoOption struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}

func (uc *GiteaUserClient) CreateAccessToken() (string, string, error) {
	defer functrace.Trace([]interface {
	}{uc})()
	tokens, _, err := uc.Client.ListAccessTokens(gitea.ListAccessTokensOptions{})
	if err != nil {
		return "", "", fmt.Errorf("failed to list access tokens: %v", err)
	}

	for _, token := range tokens {
		if token.Name == constants.GiteaUserTokenName {
			fmt.Printf("Found existing access token: %+v\n", token)

			_, err = uc.Client.DeleteAccessToken(token.ID)
			if err != nil {
				return "", "", fmt.Errorf("failed to delete existing access token: %v", err)
			}
			fmt.Printf("Deleted existing access token: %s\n", token.Name)
			break
		}
	}

	scopes := []gitea.AccessTokenScope{gitea.AccessTokenScopeAll}

	token, _, err := uc.Client.CreateAccessToken(gitea.CreateAccessTokenOption{
		Name:   constants.GiteaUserTokenName,
		Scopes: scopes,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to create access token: %v", err)
	}

	fmt.Printf("Created new access token: %+v\n", token)
	if token.Token == "" {
		return uc.Host, "", fmt.Errorf("created token but token value is empty")
	}

	go func() {
		time.Sleep(constants.GiteaUserTokenExpireTime * time.Minute)
		_, err := uc.Client.DeleteAccessToken(token.ID)
		if err != nil {
			fmt.Printf("Failed to delete access token: %v\n", err)
		} else {
			fmt.Printf("Deleted access token: %s\n", token.Name)
		}
	}()

	return uc.Host, token.Token, nil
}

func (uc *GiteaUserClient) ValidateAccessToken(token string) (bool, error) {
	defer functrace.Trace([]interface {
	}{uc, token})()
	client, err := gitea.NewClient(uc.Host, gitea.SetToken(token))
	if err != nil {
		return false, fmt.Errorf("failed to create gitea client: %v", err)
	}

	user, _, err := client.GetMyUserInfo()
	if err != nil {
		return false, fmt.Errorf("failed to validate access token: %v", err)
	}

	fmt.Printf("Access token is valid for user: %s\n", user.UserName)
	return true, nil
}
