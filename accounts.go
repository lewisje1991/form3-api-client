package client

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type (
	AccountService struct {
		httpService *HTTPService
		resourceURL string
	}

	Account struct {
		Data struct {
			Attributes struct {
				AccountClassification       string      `json:"account_classification"`
				AccountNumber               string      `json:"account_number"`
				AlternativeBankAccountNames interface{} `json:"alternative_bank_account_names"`
				BankID                      string      `json:"bank_id"`
				BankIDCode                  string      `json:"bank_id_code"`
				BaseCurrency                string      `json:"base_currency"`
				Bic                         string      `json:"bic"`
				Country                     string      `json:"country"`
				CustomerID                  string      `json:"customer_id"`
				Iban                        string      `json:"iban"`
			} `json:"attributes"`
			CreatedOn      time.Time `json:"created_on"`
			ID             string    `json:"id"`
			ModifiedOn     time.Time `json:"modified_on"`
			OrganisationID string    `json:"organisation_id"`
			Type           string    `json:"type"`
			Version        int       `json:"version"`
		} `json:"data"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
	}

	AccountList struct {
		Data  []ResponseData `json:"data"`
		Links ListLinks      `json:"links"`
	}

	AccountCreateRequest struct {
		Data Data `json:"data"`
	}

	Data struct {
		ID             string            `json:"id"`
		OrganisationID string            `json:"organisation_id"`
		Type           string            `json:"type"`
		Attributes     AccountAttributes `json:"attributes"`
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

func (s *AccountService) Fetch(id string) (*Account, error) {
	req, err := s.httpService.BuildRequest(http.MethodGet, fmt.Sprintf("%s/%s", s.resourceURL, id), nil)
	if err != nil {
		return nil, err
	}

	var account = &Account{}
	err = s.httpService.Do(req, account)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	return account, nil
}

func (s *AccountService) Create(request *AccountCreateRequest) (*Account, error) {
	req, err := s.httpService.BuildRequest(http.MethodPost, s.resourceURL, request)
	if err != nil {
		return nil, err
	}

	respObj := &Account{}

	err = s.httpService.Do(req, respObj)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	return respObj, nil
}

func (s *AccountService) List(pageSize, pageNumber int64) (*AccountList, error) {
	req, err := s.httpService.BuildRequest(http.MethodGet, s.resourceURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("page[number]", strconv.FormatInt(pageNumber, 10))
	q.Add("page[size]", strconv.FormatInt(pageSize, 10))
	req.URL.RawQuery = q.Encode()

	respObj := &AccountList{}

	err = s.httpService.Do(req, respObj)
	if err != nil {
		return nil, err
	}

	return respObj, nil
}

func (s *AccountService) Delete(id string, version int64) error {
	req, err := s.httpService.BuildRequest(http.MethodDelete, fmt.Sprintf("%s/%s", s.resourceURL, id), nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("version", strconv.FormatInt(version, 10))
	req.URL.RawQuery = q.Encode()

	err = s.httpService.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

// todo document process.
// docker tests run initial setup...
