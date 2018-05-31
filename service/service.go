package service

import (
	"log"
	"net/http"
	"os"
	"time"
)

// CreateNewClient is responsible for creating
// a new http.Client with specified timeout
func CreateNewClient() http.Client {
	timeout := time.Duration(5 * time.Second)
	return http.Client{Timeout: timeout}
}

// MakeGetRequest creates a GET request to resurce specified with the url parameter
func MakeGetRequest(url string, client http.Client) (*http.Response, error) {
	return makeRequest(url, client, "GET")
}

func makeRequest(url string, client http.Client, method string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	req.SetBasicAuth(os.Getenv("GITHUB_ACCOUNT"), os.Getenv("GITHUB_PASS"))
	if err != nil {
		log.Fatal(err)
	}

	return client.Do(req)
}
