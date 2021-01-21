package client

import "errors"

type Client struct {
	Host string
}

func NewClient(host string) (*Client, error) {

	return nil, errors.New("host cannot be empty")

}
