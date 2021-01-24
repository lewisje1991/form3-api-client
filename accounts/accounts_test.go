package accounts

import (
	"fmt"
	"testing"
)

const ENDPOINT = "http://accountapi:8080/v1"

func TestFetchAccount(t *testing.T) {
	client, err := NewClient(ENDPOINT)
	if err != nil {
		t.Fatal("Didn't want error but got one", err)
	}

	acc, err := client.Fetch("e06fa40b-872d-442e-a1a3-ff0ae5d2b866")
	if err != nil {
		t.Fatal("Didn't want error but got one", err)
	}

	fmt.Printf("%+v\n", acc)
}

func TestCreateAccount(t *testing.T) {
	client, err := NewClient(ENDPOINT)
	if err != nil {
		t.Fatal("Didn't want error but got one", err)
	}

	acc, err := client.Create(&RequestData{
		Data: DataProperties{
			Attributes: Attributes{
				Country:               "GB",
				BaseCurrency:          "GBP",
				BankID:                "400302",
				BankIDCode:            "GBDSC",
				AccountNumber:         "10000004",
				CustomerID:            "234",
				Iban:                  "GB28NWBK40030212764204",
				Bic:                   "NWBKGB42",
				AccountClassification: "Personal",
			},
			ID:             "bf5950f8-e120-49f4-ba70-c317180fb878",
			OrganisationID: "8aba657d-1e98-45e4-af35-6d91a1f21b39",
			Type:           "accounts",
		},
	})

	if err != nil {
		t.Fatal("Didn't want error but got one", err)
	}

	fmt.Printf("%+v\n", acc)
}

func TestList(t *testing.T) {
	client, err := NewClient(ENDPOINT)
	if err != nil {
		t.Fatal("Didn't want error but got one", err)
	}

	acc, err := client.List(1, 2)

	if err != nil {
		t.Fatal("Didn't want error but got one", err)
	}	

	fmt.Printf("%+v\n", acc)
}
