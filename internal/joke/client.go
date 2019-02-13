package joke

import (
	"context"
	"github.com/ctco-dev/go-api-template/internal/log"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type (
	Client interface {
		GetJoke(ctx context.Context) (Response, error)
	}

	Response struct {
		ID       string   `json:"id"`
		Category []string `json:"category"`
		IconURL  string   `json:"icon_url"`
		URL      string   `json:"url"`
		Value    string   `json:"value"`
	}

	ChuckNorrisAPIClient struct {
		url        string
		httpClient http.Client
	}
)

func NewChuckNorrisAPIClient(url string) Client {

	chuckNorrisAPI := &ChuckNorrisAPIClient{
		url:        url,
		httpClient: http.Client{},
	}

	return chuckNorrisAPI
}

func (c *ChuckNorrisAPIClient) GetJoke(ctx context.Context) (jokeResp Response, err error) {

	req, err := http.NewRequest(http.MethodGet, c.url, nil)
	if err != nil {
		return jokeResp, errors.Wrap(err, "failed to create a request")
	}

	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return jokeResp, errors.Wrap(err, "failed to make a call")
	}

	log.WithCtx(ctx).WithField("status", resp.Status).Info("Got response")

	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return jokeResp, errors.Wrap(err, "failed to read response")
		}

		err = json.Unmarshal(body, &jokeResp)
		if err != nil {
			return jokeResp, errors.Wrapf(err, "can't decode body: %v", body)
		}

		return jokeResp, nil

	}

	return jokeResp, errors.Errorf("got bad response %s", resp.Status)
}
