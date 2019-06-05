package git

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/PyramidSystemsInc/go/logger"
)

func (info GitInfo) authHeader() string {
	if info.Token != "" {
		if info.TokenType != "" {
			return strings.Join([]string{info.TokenType, info.Token}, " ")
		} else {
			return strings.Join([]string{"Basic", info.Token}, " ")
		}
	} else {
		logger.Err("Token not found")
		return ""
	}
}

func (info GitInfo) endpoint() string {
	baseUrl := "https://api.github.com"
	// If an org name is provided, use the org endpoint
	// If no org name is provided, use the user endpoint
	if info.OrgName != "" {
		return strings.Join([]string{baseUrl, "orgs", info.OrgName, "repos"}, "/")
	} else {
		return strings.Join([]string{baseUrl, "user", "repos"}, "/")
	}

}

// CreateRepository creates a new git repo with the payload provided to info
func (info GitInfo) CreateRepository(repo *GitRepo) (*GitRepo, error) {
	config, err := repo.config()
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", info.endpoint(), config)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", info.authHeader())
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	gitResult := &GitRepo{}
	err = json.NewDecoder(response.Body).Decode(&gitResult)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	return gitResult, nil
}

func (repo GitRepo) config() (*bytes.Buffer, error) {
	res, err := json.Marshal(repo)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(res), nil
}

func (repo GitRepo) getUrl() string {
	return repo.Url
}

func (repo GitRepo) getSsh() string {
	return repo.SshUrl
}

func (repo GitRepo) getHooksUrl() string {
	return repo.HooksUrl
}
