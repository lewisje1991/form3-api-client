package form3

import (
	"errors"
	"fmt"
	"net/url"
)

var (
	// ErrBaseURLEmpty Cannot create a new client with an empty base url
	ErrBaseURLEmpty = errors.New("BaseURL cannot be empty")
	// ErrBaseURLInvalid Cannot create a client with an invalid url
	ErrBaseURLInvalid = errors.New("BaseURL is invalid")
)

// Client container for services to interact with form3 api's
type Client struct {
	Accounts *accountService
}

// NewClient returns a new instance of the API client for interacting with form3 api's
func NewClient(baseURL string) (*Client, error) {
	if baseURL == "" {
		return nil, ErrBaseURLEmpty
	}

	if _, err := url.ParseRequestURI(baseURL); err != nil {
		return nil, fmt.Errorf("baseURL %s: %w", baseURL, ErrBaseURLInvalid)
	}

	httpService := &httpService{
		BaseURL: baseURL,
	}

	return &Client{
		Accounts: &accountService{
			httpService: httpService,
			resourceURL: "/organisation/accounts",
		},
	}, nil
}
