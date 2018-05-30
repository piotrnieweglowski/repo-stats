package service

// Repo contains basic information about repository fetched from github.com
type Repo struct {
	Name      string `json:"name"`
	GithubURL string `json:"html_url"`
}
