package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

func (c *Client) ExecuteWithMiddleware(req *http.Request, respObj interface{}) (interface{}, error) {
	req.Header.Add("content-type", "application/vnd.api+json")
	return c.execute(req, respObj)
}

func (c *Client) execute(req *http.Request, respObj interface{}) (interface{}, error) {
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, respObj)
	if err != nil {
		return nil, fmt.Errorf("error unmashalling response: %w", err)
	}

	return respObj, nil
}
