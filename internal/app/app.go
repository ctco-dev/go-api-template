package app

import (
	"context"
	"ctco-dev/go-api-template/internal/joke"
	"ctco-dev/go-api-template/internal/log"
	"time"
)

type Specification struct {
	JokeServiceURL string `split_words:"true" required:"true" default:"https://api.chucknorris.io/jokes/random"`
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
