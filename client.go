package client

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL string

	httpClient *http.Client
}

var (
	// ErrBaseURLEmpty Cannot create a new client with an empty base url
	ErrBaseURLEmpty = errors.New("BaseURL cannot be empty")
	// ErrBaseURLInvalid Cannot create a client with an invalid url
	ErrBaseURLInvalid = errors.New("BaseURL is invalid")
)

// NewClient returns a configured instance of a client
func NewClient(baseURL string) (*Client, error) {

	if baseURL == "" {
		return nil, ErrBaseURLEmpty
	}

	if _, err := url.ParseRequestURI(baseURL); err != nil {
		return nil, fmt.Errorf("baseURL %s: %w", baseURL, ErrBaseURLInvalid)
	}

	return &Client{
		BaseURL: baseURL,
	}, nil
}
