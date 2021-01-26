package accounts

import "time"

type Entity struct {
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

type EntityList struct {
	Data  []ResponseData `json:"data"`
	Links ListLinks      `json:"links"`
}

type RequestData struct {
	Data Data `json:"data"`
}

type Data struct {
	ID             string     `json:"id"`
	OrganisationID string     `json:"organisation_id"`
	Type           string     `json:"type"`
	Attributes     Attributes `json:"attributes"`
}

type ResponseData struct {
	Data struct {
		ID             string     `json:"id"`
		OrganisationID string     `json:"organisation_id"`
		Type           string     `json:"type"`
		Attributes     Attributes `json:"attributess"`
		CreatedOn      time.Time  `json:"created_on"`
		ModifiedOn     time.Time  `json:"modified_on"`
		Version        int        `json:"version"`
	} `json:"data"`
}

type ListLinks struct {
	Self  string `json:"self"`
	First string `json:"first"`
	Last  string `json:"last"`
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

type BodyError struct {
	ErrorMessage string `json:"error_message"`
}
