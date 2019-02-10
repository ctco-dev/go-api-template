package app

import (
	"context"
	"github.com/ctco-dev/go-api-template/internal/log"
	"time"

	"github.com/ctco-dev/go-api-template/internal/joke"
)

type Specification struct {
	JokeServiceURL string `split_words:"true" required:"true" default:"https://api.chucknorris.io/jokes/random"`
	Port           int    `split_words:"true" required:"true" default:"3000"`
}

type (
	App struct {
		client joke.Client
	}

	SomeResult struct {
		Value string
		Time  string
	}
)

func NewApp(env Specification) App {

	apiClient := joke.NewChuckNorrisAPIClient(env.JokeServiceURL)

	return App{
		client: apiClient,
	}

}

func (a *App) DoSomething(ctx context.Context) (result SomeResult, err error) {

	log.WithCtx(ctx).Info("I'm trying to get new joke")

	jokeResp, err := a.client.GetJoke(ctx)
	if err != nil {
		return result, err
	}

	log.WithCtx(ctx).Info("I'm done")

	return SomeResult{
		Value: jokeResp.Value,
		Time:  time.Now().String(),
	}, nil
}
