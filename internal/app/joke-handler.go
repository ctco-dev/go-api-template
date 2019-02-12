package app

import (
	"ctco-dev/go-api-template/internal/joke"
	"ctco-dev/go-api-template/internal/log"
	"encoding/json"
	"net/http"
	"time"
)

type jokeResult struct {
	Value string
	Time  string
}

// JokeHandler handles /joke route
type jokeHandler struct {
	client joke.Client
}

func (j *jokeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, _ = ShiftPath(r.URL.Path)

	if head != "" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		j.handleGet(w, r)
	default:
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
	}
}

func (j *jokeHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.WithCtx(ctx).Info("I'm trying to get new joke")
	jokeResp, err := j.client.GetJoke(ctx)
	if err != nil {
		http.Error(w, "Can't get error", http.StatusInternalServerError)
		return
	}

	result := jokeResult{
		Value: jokeResp.Value,
		Time:  time.Now().String(),
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Can't encode response.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)

	log.WithCtx(ctx).Info("I'm done")
}
