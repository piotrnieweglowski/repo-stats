package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/piotrnieweglowski/repo-stats/service"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/js/app.js", js)
	http.HandleFunc("/v1/users/", handleUser)
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("templates", "index.html")
	http.ServeFile(w, r, fp)
}

func js(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("templates", "js", "app.js")
	http.ServeFile(w, r, fp)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/users/")
	userService := service.CreateUserService()
	data, err := json.Marshal(userService.GetUser(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
