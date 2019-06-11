package git

import (
	"bytes"
)

type Git interface {
	authHeader() string
	endpoint() string

	GetRepository(repo GitRepo) (GitRepo, error)
	CreateRepository(info GitInfo, repo GitRepo) (GitRepo, error)
}

type GitInfo struct {
	Url, UserName, OrgName, Name string
	Token, TokenType             string
}

type Repository interface {
	config() *bytes.Buffer
	getUrl(repo string) string
	getSsh(repo string) string
	getHooksUrl(repo string) string
}

type GitRepo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	NodeId      string `json:"node_id"`
	Url         string `json:"url"`
	SshUrl      string `json:"ssh_url"`
	HooksUrl    string `json:"hooks_url"`
}
