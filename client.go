package client

import "errors"

type Client struct {
	Host string
}

var (
	// ErrHostEmpty Cannot create a new client with an empty host
	ErrHostEmpty = errors.New("host cannot be empty")
)

// NewClient returns a configured instance of a client
func NewClient(host string) (*Client, error) {

	if host == "" {
		return nil, ErrHostEmpty
	}

	return &Client{
		Host: host,
	}, nil

}
