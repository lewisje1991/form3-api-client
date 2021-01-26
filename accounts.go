package form3

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type (
	accountService struct {
		httpService *httpService
		resourceURL string
	}

	// Account respresents an account resource
	Account struct {
		Data  AccountData `json:"data"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
	}

	AccountData struct {
		Attributes     AccountAttributes `json:"attributes"`
		CreatedOn      time.Time         `json:"created_on"`
		ID             string            `json:"id"`
		ModifiedOn     time.Time         `json:"modified_on"`
		OrganisationID string            `json:"organisation_id"`
		Type           string            `json:"type"`
		Version        int               `json:"version"`
	}

	AccountList struct {
		Data  []ResponseData `json:"data"`
		Links ListLinks      `json:"links"`
	}

	AccountCreateRequest struct {
		Data AccountCreateRequestData `json:"data"`
	}

	AccountCreateRequestData struct {
		Attributes     AccountAttributes `json:"attributes"`
		ID             string            `json:"id"`
		OrganisationID string            `json:"organisation_id"`
		Type           string            `json:"type"`
	}

	ResponseData struct {
		Data struct {
			ID             string            `json:"id"`
			OrganisationID string            `json:"organisation_id"`
			Type           string            `json:"type"`
			Attributes     AccountAttributes `json:"attributess"`
			CreatedOn      time.Time         `json:"created_on"`
			ModifiedOn     time.Time         `json:"modified_on"`
			Version        int               `json:"version"`
		} `json:"data"`
	}

	ListLinks struct {
		Self  string `json:"self"`
		First string `json:"first"`
		Last  string `json:"last"`
	}

	AccountAttributes struct {
		Country               string `json:"country"`
		BaseCurrency          string `json:"base_currency"`
		BankID                string `json:"bank_id"`
		BankIDCode            string `json:"bank_id_code"`
		AccountNumber         string `json:"account_number"`
		CustomerID            string `json:"customer_id"`
		Iban                  string `json:"iban"`
		Bic                   string `json:"bic"`
		AccountClassification string `json:"account_classification"`
	}
)

// Fetch retrieves a single account by ID
func (s *accountService) Fetch(id string) (*Account, error) {
	req, err := s.httpService.buildRequest(http.MethodGet, fmt.Sprintf("%s/%s", s.resourceURL, id), nil)
	if err != nil {
		return nil, err
	}

	var account = &Account{}
	err = s.httpService.do(req, account)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	return account, nil
}

// Create creates a new account
func (s *accountService) Create(request *AccountCreateRequest) (*Account, error) {
	req, err := s.httpService.buildRequest(http.MethodPost, s.resourceURL, request)
	if err != nil {
		return nil, err
	}

	respObj := &Account{}

	err = s.httpService.do(req, respObj)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	return respObj, nil
}

// List returns multiple accounts and can be paginated
func (s *accountService) List(pageSize, pageNumber int64) (*AccountList, error) {
	req, err := s.httpService.buildRequest(http.MethodGet, s.resourceURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("page[number]", strconv.FormatInt(pageNumber, 10))
	q.Add("page[size]", strconv.FormatInt(pageSize, 10))
	req.URL.RawQuery = q.Encode()

	respObj := &AccountList{}

	err = s.httpService.do(req, respObj)
	if err != nil {
		return nil, err
	}

	return respObj, nil
}

// Delete removes an account by id and version
func (s *accountService) Delete(id string, version int64) error {
	req, err := s.httpService.buildRequest(http.MethodDelete, fmt.Sprintf("%s/%s", s.resourceURL, id), nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("version", strconv.FormatInt(version, 10))
	req.URL.RawQuery = q.Encode()

	err = s.httpService.do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

// todo document process.
// docker tests run initial setup...
