package gitea

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
	"github.com/toheart/functrace"
)

func CreateUser(username, password, email string) (*gitea.User, error) {
	defer functrace.Trace([]interface {
	}{username, password, email})()
	client, err := GetGiteaAdminClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get admin client: %v", err)
	}

	user := gitea.CreateUserOption{
		Username:           username,
		Password:           password,
		Email:              email,
		MustChangePassword: gitea.OptionalBool(false),
		SendNotify:         false,
	}

	userResp, _, err := client.AdminCreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return userResp, nil
}

func UpdateUserPassword(username string, newPassword string) (*gitea.Response, error) {
	defer functrace.Trace([]interface {
	}{username, newPassword})()
	client, err := GetGiteaAdminClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get admin client: %v", err)
	}

	userEditOptions := gitea.EditUserOption{
		LoginName:          username,
		Password:           newPassword,
		MustChangePassword: gitea.OptionalBool(false),
	}

	userResp, err := client.AdminEditUser(username, userEditOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to update user password: %v", err)
	}

	return userResp, nil
}
