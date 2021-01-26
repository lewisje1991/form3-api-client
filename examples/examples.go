package main

import (
	"log"

	"github.com/google/uuid"
	"github.com/lewisje1991/f3-accounts-api-client/accounts"
)

func main() {
	client, err := accounts.NewClient("http://localhost:8080/v1")
	if err != nil {
		log.Fatal("error creating client", err)
	}

	account, err := client.Create(&accounts.RequestData{
		Data: accounts.Data{
			ID:             uuid.NewString(),
			OrganisationID: uuid.NewString(),
			Type:           "accounts",
			Attributes: accounts.Attributes{
				Country:               "GB",
				AccountClassification: "Personal",
			},
		},
	})

	if err != nil {
		log.Fatal("error creating account", err)
	}

	log.Printf("Account %+v", account)

}
