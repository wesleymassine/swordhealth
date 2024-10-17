package userservice

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"syscall"

	"github.com/hashicorp/go-retryablehttp"
)

type Client struct {
	config *Config
	client *retryablehttp.Client
}

type Config struct {
	Endpoint string `env:"USER_SERVICE_URL" envDefault:"http://localhost:8081"`
	APIKey   string `env:"USER_SERVICE_API_KEY" envDefault:"swordhealth"`
}

type Response struct {
	*http.Response
}

var (
	ErrUserServiceConnectionRefused = errors.New("user-service: connection refused")
	ErrInvalidConfig                = errors.New("invalid configuration: missing endpoint or API key")
	ErrEmptyResponseBody            = errors.New("response body is empty")
)

// NewConfig initializes the configuration for the user service.
func NewConfig() *Config {
	return &Config{
		Endpoint: "http://localhost:8082",
		APIKey:   "swordhealth",
	}
}

// NewRetryableHTTPClient creates and returns a new retryable HTTP client.
func NewRetryableHTTPClient() *retryablehttp.Client {
	client := retryablehttp.NewClient()
	client.RetryMax = 5 // Set the number of retry attempts
	return client
}

// NewClient initializes a new client with the provided configuration and retryable client.
func NewClient(config *Config, retryClient *retryablehttp.Client) *Client {
	// if config.Endpoint == "" || config.APIKey == "" {
	// 	return nil, ErrInvalidConfig
	// }
	return &Client{
		config: config,
		client: retryClient,
	}
}

// NewRequest creates a new HTTP request with appropriate headers and context.
func (c *Client) NewRequest(ctx context.Context, method, path string, body io.Reader) (*retryablehttp.Request, error) {
	url := c.config.Endpoint + path
	req, err := retryablehttp.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.config.APIKey)
	return req, nil
}

// DoRequest executes the HTTP request and handles common network errors.
func (c *Client) DoRequest(ctx context.Context, req *retryablehttp.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)

	// Handle specific errors like connection refused
	if errors.Is(err, syscall.ECONNREFUSED) {
		return nil, ErrUserServiceConnectionRefused
	}

	// Return the response without closing the body here
	return resp, err
}

// Decode decodes the response body into the given payload.
func (r *Response) Decode(payload any) error {
	if r.Body == nil {
		return ErrEmptyResponseBody
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(payload)
}
