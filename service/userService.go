package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// UserService provides method to fetch user data from service (github.com)
type UserService interface {
	GetUser(id string) User
}

type userService struct {
	client http.Client
}

// CreateUserService creates a new instance of UserService
// httpClient is created with default timeout set to 5 sec
func CreateUserService() UserService {
	timeout := time.Duration(5 * time.Second)
	return userService{
		client: http.Client{Timeout: timeout},
	}
}

func (s userService) GetUser(id string) User {
	url := fmt.Sprintf("https://api.github.com/users/%s", id)
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(os.Getenv("GITHUB_ACCOUNT"), os.Getenv("GITHUB_PASS"))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	repoService := CreateRepoService()
	user.Repositories = repoService.GetRepo(id)

	ch := make(chan interface{}, len(user.Repositories))
	var wg sync.WaitGroup

	wg.Add(len(user.Repositories))

	go func() {
		for _, item := range user.Repositories {
			repoService.GetRepoStatistics(id, item.Name, ch)
		}
	}()

	var statistics []RepoStatistics
	go func() {
		for val := range ch {
			if stat, ok := val.([]RepoStatistics); ok {
				for _, s := range stat {
					statistics = append(statistics, s)
				}
			}

			wg.Done()
		}
	}()

	wg.Wait()
	stati := make(map[string]int)
	sumOfBytes := 0

	for _, val := range statistics {
		sumOfBytes += val.Size
		if elem, ok := stati[val.Language]; ok {
			stati[val.Language] = elem + val.Size
		} else {
			stati[val.Language] = val.Size
		}
	}

	for k, v := range stati {
		fraction := float64(v) / float64(sumOfBytes)
		percentage := Round(fraction*100, 0.01)
		user.Statistics = append(user.Statistics, RepoStatistics{Language: k, Size: v, Percentage: percentage})
	}

	return user
}
