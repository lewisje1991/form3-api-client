package client

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

type Client struct {
	Accounts *AccountService
}

func NewClient(baseURL string) (*Client, error) {

	if baseURL == "" {
		return nil, ErrBaseURLEmpty
	}

	if _, err := url.ParseRequestURI(baseURL); err != nil {
		return nil, fmt.Errorf("baseURL %s: %w", baseURL, ErrBaseURLInvalid)
	}

	httpService := &HTTPService{
		BaseURL: baseURL,
	}

	return &Client{
		Accounts: &AccountService{
			httpService: httpService,
			resourceURL: "/organisation/accounts",
		},
	}, nil
}
