package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
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
	timeout := time.Duration(5 * time.Second)
	return repoService{
		client: http.Client{Timeout: timeout},
	}
}

func (s repoService) GetRepo(userID string) []Repo {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", userID)
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(os.Getenv("GITHUB_ACCOUNT"), os.Getenv("GITHUB_PASS"))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := s.client.Do(req)
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
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(os.Getenv("GITHUB_ACCOUNT"), os.Getenv("GITHUB_PASS"))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var objmap map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&objmap); err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	var returnVal []RepoStatistics
	for k, v := range objmap {
		returnVal = append(returnVal, RepoStatistics{Language: k, Size: v})
	}

	ch <- returnVal
}
