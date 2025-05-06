package gitea

import (
	"code.gitea.io/sdk/gitea"
	"fmt"
	"github.com/toheart/functrace"
	"strconv"
)

func (uc *GiteaUserClient) CreateRepo(repoName, description string, private bool) (*gitea.Repository, error) {
	defer functrace.Trace([]interface {
	}{uc, repoName, description, private})()

	repoOption := gitea.CreateRepoOption{
		Name:        repoName,
		Description: description,
		Private:     private,
	}

	repo, _, err := uc.Client.CreateRepo(repoOption)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %v", err)
	}

	fmt.Printf("Created new repository: %+v\n", repo)
	return repo, nil
}

func (uc *GiteaUserClient) ListUserRepos(page int, perPage int) ([]*gitea.Repository, int, error) {
	defer functrace.Trace([]interface {
	}{uc, page, perPage})()
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	repos, resp, err := uc.Client.ListMyRepos(gitea.ListReposOptions{
		ListOptions: gitea.ListOptions{
			Page:     page,
			PageSize: perPage,
		},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list repositories: %v", err)
	}

	totalReposStr := resp.Header.Get("X-Total-Count")
	totalRepos, err := strconv.Atoi(totalReposStr)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse total repository count: %v", err)
	}

	fmt.Printf("Fetched %d repositories on page %d out of %d total repositories\n", len(repos), page, totalRepos)
	return repos, totalRepos, nil
}

func (uc *GiteaUserClient) DeleteRepo(owner, repoName string) error {
	defer functrace.Trace([]interface {
	}{uc, owner, repoName})()

	resp, err := uc.Client.DeleteRepo(owner, repoName)
	if err != nil {
		return fmt.Errorf("failed to delete repository: %v", err)
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}

	fmt.Printf("Deleted repository: %s/%s\n", owner, repoName)
	return nil
}
