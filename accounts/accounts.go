package accounts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/lewisje1991/f3-accounts-api-client/client"
)

type Client struct {
	client.Client

	path string
}

func NewClient(baseURL string) (*Client, error) {
	a := &Client{
		path: "/organisation/accounts",
	}
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

	res, err := c.ExecuteWithMiddleware(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	respObj := &Entity{}

	err = json.Unmarshal(body, respObj)
	if err != nil {
		return nil, fmt.Errorf("error unmashalling response: %w", err)
	}

	return respObj, nil
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

	res, err := c.ExecuteWithMiddleware(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	respObj := &Entity{}

	err = json.Unmarshal(body, respObj)
	if err != nil {
		return nil, fmt.Errorf("error unmashalling response: %w", err)
	}

	return respObj, nil
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

	res, err := c.ExecuteWithMiddleware(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	respObj := &EntityList{}

	err = json.Unmarshal(body, respObj)
	if err != nil {
		return nil, fmt.Errorf("error unmashalling response: %w", err)
	}

	return respObj, nil
}

func (c *Client) Delete(id string, version int64) error {
	resourceURL := c.BuildURL(fmt.Sprintf("%s/%s", c.path, id))

	req, err := http.NewRequest(http.MethodDelete, resourceURL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	q := req.URL.Query()
	q.Add("version", strconv.FormatInt(version, 10))
	req.URL.RawQuery = q.Encode()

	res, err := c.ExecuteWithMiddleware(req)
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}

	if res.StatusCode != http.StatusNoContent {
		return err
	}

	return nil
}

func (c *Client) ExecuteWithMiddleware(req *http.Request) (*http.Response, error) {
	req.Header.Add("content-type", "application/vnd.api+json")
	return c.Do(req)
}

// todo error handling from API
// todo integration testing of client possibly with dockertest if I can get full api running.
// todo document process.
// todo add example.
// todo find better way to do base func.
