package main

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/lewisje1991/form3-api-client"
)

func main() {
	c, err := form3.NewClient("http://localhost:8080/v1")
	if err != nil {
		log.Fatal("error creating client", err)
	}

	account, err := c.Accounts.Create(&form3.AccountCreateRequest{
		Data: form3.AccountCreateRequestData{
			ID:             uuid.NewString(),
			OrganisationID: uuid.NewString(),
			Type:           "accounts",
			Attributes: form3.AccountAttributes{
				Country:               "GB",
				AccountClassification: "Personal",
				AccountNumber:         "10000004",
				BankID:                "400302",
				BankIDCode:            "GBDSC",
				BaseCurrency:          "GBP",
				CustomerID:            "234",
				Iban:                  "GB28NWBK40030212764204",
			},
		},
	})

	if err != nil {
		log.Fatal("error creating account", err)
	}

	j, _ := json.Marshal(account)
	log.Printf("Account %+v", string(j))
}
