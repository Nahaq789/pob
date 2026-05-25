package sync

import (
	"net/http"
	"net/url"
	"time"
)

type ApiClient struct {
	baseURL *url.URL
	client  *http.Client
}

func NewApiClient(baseURL string) (*ApiClient, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 10
	t.MaxConnsPerHost = 10
	t.MaxIdleConnsPerHost = 10
	c := &ApiClient{
		baseURL: u,
		client: &http.Client{
			Timeout:   10 * time.Second,
			Transport: t,
		},
	}
	return c, nil
}

func (a *ApiClient) GetClient() *http.Client {
	return a.client
}
