package app

import (
	"ctco-dev/go-api-template/internal/joke"
	"net/http"
	"path"
	"strings"
)

//Specification is a container of app config parameters
type Specification struct {
	JokeServiceURL string `split_words:"true" required:"true" default:"https://api.chucknorris.io/jokes/random"`
	Port           int    `split_words:"true" required:"true" default:"3000"`
}

// App implements a sample http service
type app struct {
	jokeHandler *jokeHandler
}

// New creates a new application
func New(env Specification) http.Handler {
	return &app{jokeHandler: &jokeHandler{
		client: joke.NewChuckNorrisAPIClient(env.JokeServiceURL),
	}}
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	if head == "joke" {
		a.jokeHandler.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Not Found", http.StatusNotFound)
}

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
