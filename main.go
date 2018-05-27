package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/v1/users/", handleUsers)
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("templates", "index.html")
	http.ServeFile(w, r, fp)
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/users/")
	message := fmt.Sprintf("Passed parameter id is: %s", id)
	data, err := json.Marshal(message)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
