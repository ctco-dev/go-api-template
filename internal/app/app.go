package app

import (
	"context"
	"net/http"
	"path"
	"strings"

	mongobeer "github.com/ctco-dev/go-api-template/internal/db"

	"github.com/ctco-dev/go-api-template/internal/joke"
)

//Specification is a container of app config parameters
type Specification struct {
	JokeServiceURL  string `split_words:"true" required:"true" default:"https://api.chucknorris.io/jokes/random"`
	Port            int    `split_words:"true" required:"true" default:"3000"`
	MongoHost       string `required:"true" default:"mongodb://localhost:27017"`
	MongoDatabase   string `required:"true" default:"template"`
	MongoCollection string `required:"true" default:"beer"`
}

// App implements a sample http service
type app struct {
	jokeHandler *jokeHandler
	beerHandler *beerHandler
}

// New creates a new application
func New(ctx context.Context, env Specification) http.Handler {
	return &app{
		jokeHandler: &jokeHandler{client: joke.NewChuckNorrisAPIClient(env.JokeServiceURL)},
		beerHandler: &beerHandler{
			repository: mongobeer.NewRepo(ctx, env.MongoHost, env.MongoDatabase, env.MongoCollection),
		},
	}
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	switch head {
	case "joke":
		a.jokeHandler.ServeHTTP(w, r)
	case "beer":
		a.beerHandler.ServeHTTP(w, r)
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
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
