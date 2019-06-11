package git

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/PyramidSystemsInc/go/errors"
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
	gitResult, err := info.gitRequest(repo, http.MethodGet)
	if err != nil {
		return nil, err
	}
	if gitResult.Name != "" {
		logger.Warn("Git repository already existed. Confirm that this is the correct repo prior to pushing commits.")
	} else {
		gitResult, err = info.gitRequest(repo, http.MethodPost)
		if err != nil {
			return nil, err
		}
	}
	return gitResult, nil
}

func (info GitInfo) gitRequest(repo *GitRepo, action string) (*GitRepo, error) {
	config, err := repo.config()
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(action, info.endpoint(), config)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", info.authHeader())
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode > 400 {
		return nil, errors.New(response.Status)
	}

	// git results - get returns an array; post/create returns single json element; separate out commands
	gitResult := &GitRepo{}
	dec := json.NewDecoder(response.Body)

	if action == http.MethodGet {
		// get returns an array of repositories
		_, err = dec.Token()
		if err != nil {
			logger.Err("error finding token")
		}
		for dec.More() {
			err = dec.Decode(&gitResult)
			if err != nil {
				logger.Err("error", err.Error())
				return nil, err
			}
			if gitResult.Name == repo.Name {
				return gitResult, nil
			}
		}
	} else if action == http.MethodPost {
		err = json.NewDecoder(response.Body).Decode(&gitResult)
		if err != nil {
			logger.Err("error", err.Error())
			return nil, err
		}
		return gitResult, nil

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
