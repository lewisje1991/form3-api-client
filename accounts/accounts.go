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

	resourceURL string
	baseURL     string
}

func NewClient(baseURL string) (*Client, error) {
	if baseURL == "" {
		return nil, ErrBaseURLEmpty
	}

	if _, err := url.ParseRequestURI(baseURL); err != nil {
		return nil, fmt.Errorf("baseURL %s: %w", baseURL, ErrBaseURLInvalid)
	}

	c := &Client{
		resourceURL: "/organisation/accounts",
		baseURL:     baseURL,
	}

	return c, nil
}

func (c *Client) Fetch(id string) (*Entity, error) {
	req, err := c.buildRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.resourceURL, id), nil)
	if err != nil {
		return nil, err
	}

	var account = &Entity{}
	err = c.do(req, account)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	return account, nil
}

func (c *Client) Create(requestData *RequestData) (*Entity, error) {
	req, err := c.buildRequest(http.MethodPost, c.resourceURL, requestData)
	if err != nil {
		return nil, err
	}

	respObj := &Entity{}

	err = c.do(req, respObj)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	return respObj, nil
}

func (c *Client) List(pageSize, pageNumber int64) (*EntityList, error) {
	req, err := c.buildRequest(http.MethodGet, c.resourceURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("page[number]", strconv.FormatInt(pageNumber, 10))
	q.Add("page[size]", strconv.FormatInt(pageSize, 10))
	req.URL.RawQuery = q.Encode()

	respObj := &EntityList{}

	err = c.do(req, respObj)
	if err != nil {
		return nil, err
	}

	return respObj, nil
}

func (c *Client) Delete(id string, version int64) error {
	req, err := c.buildRequest(http.MethodDelete, fmt.Sprintf("%s/%s", c.resourceURL, id), nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("version", strconv.FormatInt(version, 10))
	req.URL.RawQuery = q.Encode()

	err = c.do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) do(req *http.Request, respObj interface{}) error {
	res, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}

	defer res.Body.Close()

	if err := checkResponseForError(res); err != nil {
		return err
	}

	err = json.NewDecoder(res.Body).Decode(respObj)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}
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
		err = json.Unmarshal(data, message)
		if err != nil {
			return fmt.Errorf("error mashalling body: %w", err)
		}
	}

	return &APIError{
		StatusCode: res.StatusCode,
		Message:    message.ErrorMessage,
	}
}

func (c *Client) buildRequest(method string, resourceURL string, bodyData interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if bodyData != nil {
		bodyData, err := json.Marshal(bodyData)
		if err != nil {
			return nil, fmt.Errorf("error mashalling body: %w", err)
		}

		buf = bytes.NewBuffer(bodyData)
	}

	url := c.buildURL(resourceURL)

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("content-type", "application/vnd.api+json")

	return req, nil
}

func (c *Client) buildURL(path string) string {
	return fmt.Sprintf("%s%s", c.baseURL, path)
}

// todo document process.
// docker tests run initial setup...
