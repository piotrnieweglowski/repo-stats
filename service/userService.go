package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
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
	return userService{
		client: CreateNewClient(),
	}
}

// GetUser returns basic user data for the user with
// particular id. Repositories data are included as well
func (s userService) GetUser(id string) User {
	url := fmt.Sprintf("https://api.github.com/users/%s", id)
	resp, err := MakeGetRequest(url, s.client)
	if err != nil {
		log.Fatal(err)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	getRepositoriesData(&user)
	return user
}

func getRepositoriesData(user *User) {
	repoService := CreateRepoService()
	user.Repositories = repoService.GetRepo(user.Login)

	statistics := make([]RepoStatistics, 0)
	ch := make(chan interface{}, len(user.Repositories))
	var wg sync.WaitGroup
	wg.Add(len(user.Repositories))
	go getRepositoryStatistics(repoService, user, ch)
	go reciveRepositoryStatisticsData(&statistics, &wg, ch)

	wg.Wait()
	fillStatistics(user, statistics)
}

func getRepositoryStatistics(repoService RepoService, user *User, ch chan interface{}) {
	for _, item := range user.Repositories {
		repoService.GetRepoStatistics(user.Login, item.Name, ch)
	}
}

func reciveRepositoryStatisticsData(statistics *[]RepoStatistics, wg *sync.WaitGroup, ch chan interface{}) {
	for val := range ch {
		if stat, ok := val.([]RepoStatistics); ok {
			for _, s := range stat {
				*statistics = append(*statistics, s)
			}
		}

		wg.Done()
	}
}

func fillStatistics(user *User, statistics []RepoStatistics) {
	sumOfBytes := getTotalSumOfBytes(statistics)
	groupedStatistics := groupStatisticsByLanguage(statistics)

	for k, v := range groupedStatistics {
		fraction := float64(v) / float64(sumOfBytes)
		percentage := Round(fraction*100, 0.01)
		user.Statistics = append(user.Statistics, RepoStatistics{Language: k, Size: v, Percentage: percentage})
	}
}
