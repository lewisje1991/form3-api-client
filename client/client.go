package client

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

var (
	// ErrBaseURLEmpty Cannot create a new client with an empty base url
	ErrBaseURLEmpty = errors.New("BaseURL cannot be empty")
	// ErrBaseURLInvalid Cannot create a client with an invalid url
	ErrBaseURLInvalid = errors.New("BaseURL is invalid")
)

type Client struct {
	baseURL string

	http.Client
}

func (c *Client) SetBaseUrl(baseURL string) error {
	if baseURL == "" {
		return ErrBaseURLEmpty
	}

	if _, err := url.ParseRequestURI(baseURL); err != nil {
		return fmt.Errorf("baseURL %s: %w", baseURL, ErrBaseURLInvalid)
	}

	c.baseURL = baseURL
	return nil
}

func (c *Client) BuildURL(path string) string {
	return fmt.Sprintf("%s%s", c.baseURL, path)
}
