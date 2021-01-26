package accounts

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

var (
	// ErrBaseURLEmpty Cannot create a new client with an empty base url
	ErrBaseURLEmpty = errors.New("BaseURL cannot be empty")
	// ErrBaseURLInvalid Cannot create a client with an invalid url
	ErrBaseURLInvalid = errors.New("BaseURL is invalid")
)

type Client struct {
	http.Client

	path    string
	baseURL string
}

func NewClient(baseURL string) (*Client, error) {
	if baseURL == "" {
		return nil, ErrBaseURLEmpty
	}

	if _, err := url.ParseRequestURI(baseURL); err != nil {
		return nil, fmt.Errorf("baseURL %s: %w", baseURL, ErrBaseURLInvalid)
	}

	c := &Client{
		path:    "/organisation/accounts",
		baseURL: baseURL,
	}

	return c, nil
}

func (c *Client) Fetch(id string) (*Entity, error) {
	resourceURL := c.buildURL(fmt.Sprintf("%s/%s", c.path, id))

	req, err := http.NewRequest(http.MethodGet, resourceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	var account = &Entity{}
	err = c.do(req, account)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	return account, nil
}

func (c *Client) Create(a *RequestData) (*Entity, error) {
	resourceURL := c.buildURL(c.path)

	postBody, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("error mashalling body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, resourceURL, bytes.NewBuffer(postBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	respObj := &Entity{}

	err = c.do(req, respObj)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	return respObj, nil
}

func (c *Client) List(pageSize, pageNumber int64) (*EntityList, error) {
	resourceURL := c.buildURL(c.path)

	req, err := http.NewRequest(http.MethodGet, resourceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	q := req.URL.Query()
	q.Add("page[number]", strconv.FormatInt(pageNumber, 10))
	q.Add("page[size]", strconv.FormatInt(pageSize, 10))
	req.URL.RawQuery = q.Encode()

	respObj := &EntityList{}

	err = c.do(req, respObj)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	return respObj, nil
}

func (c *Client) Delete(id string, version int64) error {
	resourceURL := c.buildURL(fmt.Sprintf("%s/%s", c.path, id))

	req, err := http.NewRequest(http.MethodDelete, resourceURL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	q := req.URL.Query()
	q.Add("version", strconv.FormatInt(version, 10))
	req.URL.RawQuery = q.Encode()

	err = c.do(req, nil)
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}

	return nil
}

func (c *Client) do(req *http.Request, respObj interface{}) error {

	req.Header.Add("content-type", "application/vnd.api+json")

	res, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}

	defer res.Body.Close()

	if err := checkResponseForError(res); err != nil {
		return err
	}

	decErr := json.NewDecoder(res.Body).Decode(respObj)

	if decErr == io.EOF {
		decErr = nil
	}
	if decErr != nil {
		err = decErr
	}
	return nil
}

func checkResponseForError(res *http.Response) error {
	if c := res.StatusCode; http.StatusOK <= c && c < http.StatusMultipleChoices {
		return nil
	}

	message := &BodyError{}

	data, err := ioutil.ReadAll(res.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, message)
	}

	return &APIError{
		StatusCode: res.StatusCode,
		Message:    message.ErrorMessage,
	}
}

func (c *Client) buildURL(path string) string {
	return fmt.Sprintf("%s%s", c.baseURL, path)
}

// todo better error handling e.g include error message.
// todo document process.
// docker tests run initial setup...
