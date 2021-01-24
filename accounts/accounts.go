package accounts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/lewisje1991/f3-accounts-api-client/client"
)

type Client struct {
	client.Client

	path string "/organisation/accounts"
}

func NewClient(baseURL string) (*Client, error) {
	a := &Client{}
	err := a.SetBaseUrl(baseURL)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (c *Client) Fetch(id string) (*Entity, error) {
	resourceURL := c.BuildURL(fmt.Sprintf("%s/%s", c.path, id))

	req, err := http.NewRequest(http.MethodGet, resourceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// fmt.Printf("%+v", req)

	res, err := c.ExecuteWithMiddleware(req, &Entity{})
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	return res.(*Entity), err
}

func (c *Client) Create(a *RequestData) (*Entity, error) {
	resourceURL := c.BuildURL(c.path)

	postBody, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("error mashalling body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, resourceURL, bytes.NewBuffer(postBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	fmt.Printf("%+v", req)

	res, err := c.ExecuteWithMiddleware(req, &Entity{})
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	return res.(*Entity), nil
}

func (c *Client) List(pageSize, pageNumber int64) (*EntityList, error) {
	resourceURL := c.BuildURL(c.path)

	req, err := http.NewRequest(http.MethodGet, resourceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	q := req.URL.Query()
	q.Add("page[number]", strconv.FormatInt(pageNumber, 10))
	q.Add("page[size]", strconv.FormatInt(pageSize, 10))
	req.URL.RawQuery = q.Encode()

	res, err := c.ExecuteWithMiddleware(req, &EntityList{})
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	return res.(*EntityList), nil
}

func (c *Client) Delete(id string, version int64) (*EntityList, error) {
	resourceURL := c.BuildURL(fmt.Sprintf("%s/%s", c.path, id))

	req, err := http.NewRequest(http.MethodDelete, resourceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	q := req.URL.Query()
	q.Add("version", strconv.FormatInt(version, 10))
	req.URL.RawQuery = q.Encode()

	res, err := c.ExecuteWithMiddleware(req, &EntityList{})
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	return res.(*EntityList), nil
}

// todo error handling from API
// todo integration testing of client possibly with dockertest if I can get full api running.
// todo document process.
// todo add example.
// todo find better way to do base func.
