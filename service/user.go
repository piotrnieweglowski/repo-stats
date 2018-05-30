package service

// User contains basic user information pulled from service (github.com)
// It contains information about repositories and language
// statistics for repositories as well
type User struct {
	Login        string `json:"login"`
	Email        string `json:"email"`
	Blog         string `json:"blog"`
	AccountType  string `json:"type"`
	AvatarURL    string `json:"avatar_url"`
	GithubURL    string `json:"html_url"`
	Repositories []Repo
	Statistics   []RepoStatistics
}
