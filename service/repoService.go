package service

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sync"
)

// RepoService provides method to fetch data of the repository from service (github.com)
type RepoService interface {
	GetRepos(user User) []Repo
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

// GetRepos retrieves all public reopsitories for a particular user
func (s repoService) GetRepos(user User) []Repo {
	reposToFetch := user.PublicRepos
	pageSize := 100
	page := 1
	chunksCount := int(math.Ceil(float64(reposToFetch) / float64(pageSize)))
	allRepos := make([]Repo, 0)

	ch := make(chan interface{}, chunksCount)
	var wg sync.WaitGroup
	wg.Add(chunksCount)

	for reposToFetch > 0 {
		url := fmt.Sprintf("https://api.github.com/users/%s/repos?page=%d&per_page=%d", user.Login, page, pageSize)
		go getReposChunk(url, s.client, ch)

		page++
		reposToFetch -= pageSize
	}

	go reciveReposChunk(&allRepos, &wg, ch)
	wg.Wait()

	return allRepos
}

func getReposChunk(url string, client http.Client, ch chan interface{}) {
	resp, err := MakeGetRequest(url, client)
	if err != nil {
		log.Fatal(err)
	}

	repos := make([]Repo, 0)
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	ch <- repos
}

func reciveReposChunk(allRepos *[]Repo, wg *sync.WaitGroup, ch chan interface{}) {
	for val := range ch {
		if repos, ok := val.([]Repo); ok {
			for _, repo := range repos {
				*allRepos = append(*allRepos, repo)
			}
		}

		wg.Done()
	}
}

// GetRepoStatistics returns language statistics for a given repository
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
