package udl

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/time/rate"
)

const (
	DefaultClientTimeout      = 3 * time.Second
	DefaultRateLimitPerMinute = 10 // limits to 10 requests per minute
)

type Client struct {
	baseURL     *url.URL
	client      http.Client
	rateLimiter *rate.Limiter
}

func New(baseURL string) (Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return Client{}, err
	}
	c := http.Client{Timeout: DefaultClientTimeout}
	r := rate.NewLimiter(rate.Every(60*time.Second), DefaultRateLimitPerMinute)
	return Client{
		baseURL:     u,
		client:      c,
		rateLimiter: r,
	}, nil
}

func (c Client) Post(ctx context.Context, path string, body []byte) (Response, error) {

	return Response{}, nil
}