package accounts

import "time"

type Entity struct {
	Data  ResponseData `json:"data"`
	Links Links        `json:"links"`
}

type EntityList struct {
	Data  []ResponseData `json:"data"`
	Links ListLinks      `json:"links"`
}

type ResponseData struct {
	DataProperties
	CreatedOn  time.Time `json:"created_on"`
	ModifiedOn time.Time `json:"modified_on"`
	Version    int       `json:"version"`
}

type Links struct {
	Self string `json:"self"`
}

type ListLinks struct {
	Links
	First string `json:"first"`
	Last  string `json:"last"`
}

type RequestData struct {
	Data DataProperties `json:"data"`
}

type DataProperties struct {
	ID             string     `json:"id"`
	OrganisationID string     `json:"organisation_id"`
	Type           string     `json:"type"`
	Attributes     Attributes `json:"attributes"`
}

type Attributes struct {
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
