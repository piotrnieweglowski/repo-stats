package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// RepoService provides method to fetch data of the repository from service (github.com)
type RepoService interface {
	GetRepo(userID string) []Repo
	GetRepoStatistics(userID string, repoID string, ch chan interface{})
}

type repoService struct {
	client http.Client
}

// CreateRepoService creates a new instance of RepoService
// httpClient is created with default timeout set to 5 sec
func CreateRepoService() RepoService {
	return repoService{
		client: CreateNewClient(),
	}
}

func (s repoService) GetRepo(userID string) []Repo {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", userID)
	resp, err := MakeGetRequest(url, s.client)
	if err != nil {
		log.Fatal(err)
	}

	repos := make([]Repo, 0)
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	return repos
}

func (s repoService) GetRepoStatistics(userID string, repoID string, ch chan interface{}) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/languages", userID, repoID)
	resp, err := MakeGetRequest(url, s.client)
	if err != nil {
		log.Fatal(err)
	}

	var statisticsMap map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&statisticsMap); err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	var statistics []RepoStatistics
	for k, v := range statisticsMap {
		statistics = append(statistics, RepoStatistics{Language: k, Size: v})
	}

	ch <- statistics
}
