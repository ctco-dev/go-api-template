package joke

import (
	"context"
	"testing"
	"time"
)

func Test_ReadJoke(t *testing.T) {

	client := NewChuckNorrisAPIClient("https://api.chucknorris.io/jokes/random")

	ctx := context.Background()
	timeoutCtx, cancelFunc := context.WithTimeout(ctx, time.Second*5)
	defer cancelFunc()

	joke, err := client.GetJoke(timeoutCtx)
	if err != nil {
		t.Fatal(err)
	}

	if len(joke.Value) == 0 {
		t.Errorf("shold be a joke, got nothing")
	}

}
