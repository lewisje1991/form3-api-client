package main

import (
	"log"

	"github.com/google/uuid"
	client "github.com/lewisje1991/form3-api-client"
)

func main() {
	c, err := client.NewClient("http://localhost:8080/v1")
	if err != nil {
		log.Fatal("error creating client", err)
	}

	account, err := c.Accounts.Create(&client.AccountCreateRequest{
		Data: client.Data{
			ID:             uuid.NewString(),
			OrganisationID: uuid.NewString(),
			Type:           "accounts",
			Attributes: client.AccountAttributes{
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
