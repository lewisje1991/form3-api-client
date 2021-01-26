package form3

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type httpService struct {
	BaseURL string

	http.Client
}

type apiError struct {
	StatusCode int
	Message    string
}

type bodyError struct {
	ErrorMessage string `json:"error_message"`
}

func (e *apiError) Error() string {
	return fmt.Sprintf("status code %d: error %v", e.StatusCode, e.Message)
}

func (h *httpService) buildRequest(method string, resourceURL string, bodyData interface{}) (*http.Request, error) {
	var buf io.ReadWriter

	if bodyData != nil {
		bodyData, err := json.Marshal(bodyData)
		if err != nil {
			return nil, fmt.Errorf("error mashalling body: %w", err)
		}

		buf = bytes.NewBuffer(bodyData)
	}

	url := h.buildURL(resourceURL)

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("content-type", "application/vnd.api+json")

	return req, nil
}

func (h *httpService) do(req *http.Request, respObj interface{}) error {
	res, err := h.Client.Do(req)
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
			log.Println("geersdf")
			return err
		}
	}

	return nil
}

func (h *httpService) buildURL(path string) string {
	return fmt.Sprintf("%s%s", h.BaseURL, path)
}

func checkResponseForError(res *http.Response) error {
	if c := res.StatusCode; http.StatusOK <= c && c < http.StatusMultipleChoices {
		return nil
	}

	message := &bodyError{}

	data, err := ioutil.ReadAll(res.Body)
	if err == nil && data != nil {
		err = json.Unmarshal(data, message)
		if err != nil {
			return fmt.Errorf("error mashalling body: %w", err)
		}
	}

	return &apiError{
		StatusCode: res.StatusCode,
		Message:    message.ErrorMessage,
	}
}
