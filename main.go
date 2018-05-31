package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/piotrnieweglowski/repo-stats/service"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/js/app.js", js)
	http.HandleFunc("/css/app.css", css)
	http.HandleFunc("/v1/users/", handleUser)
	listenAndServe()
}

func listenAndServe() {
	port := os.Getenv("REPO_STATS_PORT")
	addr := fmt.Sprintf(":%s", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
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

func css(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("templates", "css", "app.css")
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
