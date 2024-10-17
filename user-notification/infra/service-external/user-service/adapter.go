package userservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/wesleymassine/swordhealth/user-notification/domain"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"user-service",
	fx.Provide(NewConfig),
	fx.Provide(NewRetryableHTTPClient),
	fx.Provide(NewClient),
	fx.Provide(New),
)

type Adapter struct {
	Config *Config
	Client *Client
}

// New creates a new Adapter instance implementing the BCPort interface
func New(config *Config, client *Client) domain.UserServiceClient {
	return &Adapter{Config: config, Client: client}
}

// GetBalancesMap fetches the balance map for the given walletID
func (a *Adapter) GetUserByTaskID(ctx context.Context, taskID int64) (*domain.User, error) {
	path := fmt.Sprintf("/api/v1/tasks/%v", taskID)

	// Send the request and handle the response
	bodyBytes, err := a.sendRequestAndGetResponseBody(ctx, http.MethodGet, path)
	if err != nil {
		log.WithError(err).Error("failed to get response for user by task ID")
		return nil, err
	}

	// Decode the response body into the User struct
	var userResponse domain.User
	if err := json.Unmarshal(bodyBytes, &userResponse); err != nil {
		log.WithError(err).Error("failed to decode user response")
		return nil, err
	}

	return &userResponse, nil
}

// sendRequestAndGetResponseBody sends an HTTP request and returns the response body as a byte slice
func (a *Adapter) sendRequestAndGetResponseBody(ctx context.Context, method, path string) ([]byte, error) {
	// Create the request
	req, err := a.Client.NewRequest(ctx, method, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for path %s: %w", path, err)
	}

	// Execute the request
	resp, err := a.Client.DoRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for path %s: %w", path, err)
	}

	// Ensure that the response is not nil
	if resp == nil {
		return nil, fmt.Errorf("received nil response for path %s", path)
	}

	// Check if the response status is an error
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("received bad response for path %s, status code: %d", path, resp.StatusCode)
	}

	// Read the response body before closing it
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body for path %s: %w", path, err)
	}

	// Close the response body after reading it
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.WithError(err).Error("failed to close response body")
		}
	}()

	return bodyBytes, nil
}
