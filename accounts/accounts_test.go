package accounts

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

const ENDPOINT = "http://accountapi:8080/v1"

func TestNewClient(t *testing.T) {

	t.Run("Returns error if baseUrl is empty", func(t *testing.T) {
		_, got := NewClient("")
		assertError(t, got)

		want := ErrBaseURLEmpty

		if !errors.Is(got, want) {
			t.Fatalf("wanted %s, but got %s", want, got)
		}
	})

	t.Run("Returns error if baseUrl is invalid", func(t *testing.T) {
		_, got := NewClient("trdtest")
		assertError(t, got)

		want := ErrBaseURLInvalid

		if !errors.Is(got, want) {
			t.Fatalf("wanted %s, but got %s", ErrBaseURLInvalid, want)
		}
	})

	t.Run("Returns client if validation passes", func(t *testing.T) {
		client, err := NewClient("https://api.form3.com")
		assertNoError(t, err)

		got := client.baseURL
		want := "https://api.form3.com"

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("wanted %+v, but got %+v", want, got)
		}
	})
}

func TestCreate(t *testing.T) {
	client, err := NewClient(ENDPOINT)
	assertNoError(t, err)

	t.Run("Successful create returns created account", func(t *testing.T) {

		got := createAccount(t, client)

		respData, err := loadTestData("valid-create-response.json", &Entity{})
		assertNoError(t, err)

		want := respData.(*Entity)
		want.Data.ID = got.Data.ID
		want.Links.Self = fmt.Sprintf("%s/%s", "/v1/organisation/accounts/", got.Data.ID)

		// Set times to match so that we can deep compare... is there a better way?

		now := time.Now()
		got.Data.CreatedOn = now
		want.Data.CreatedOn = now

		got.Data.ModifiedOn = now
		want.Data.ModifiedOn = now

		if reflect.DeepEqual(want, got) {
			t.Fatalf("wanted %+v, but got %+v", got, want)
		}

		cleanUpAccounts(t, []string{got.Data.ID}, client)
	})

	t.Run("Validation errors are returned", func(t *testing.T) {
		testReqData, err := loadTestData("valid-create-request.json", &RequestData{})
		assertNoError(t, err)

		requestData := testReqData.(*RequestData)
		requestData.Data.OrganisationID = ""

		_, got := client.Create(requestData)
		assertErrorContains(t, "organisation_id in body is required", got)
	})
}

func TestFetch(t *testing.T) {
	client, err := NewClient(ENDPOINT)
	assertNoError(t, err)

	t.Run("When account exists, returns entity model", func(t *testing.T) {
		account := createAccount(t, client)

		want, err := client.Fetch(account.Data.ID)
		assertNoError(t, err)

		fmt.Println(want)
		// todo compare response model

		cleanUpAccounts(t, []string{account.Data.ID}, client)
	})

	t.Run("When account does not exist, returns error", func(t *testing.T) {
		ID := uuid.New().String()
		_, got := client.Fetch(ID)
		assertErrorContains(t, "status code 404", got)
	})
}

func TestDelete(t *testing.T) {
	client, err := NewClient(ENDPOINT)
	assertNoError(t, err)

	t.Run("Delete returns error when resource not found", func(t *testing.T) {
		got := client.Delete("invalid resource uuid invalid", 0)
		assertErrorContains(t, "status code 400", got)
	})

	t.Run("Delete returns error when version invalid", func(t *testing.T) {
		account := createAccount(t, client)
		got := client.Delete(account.Data.ID, 2)
		assertErrorContains(t, "status code 404", got)
		cleanUpAccounts(t, []string{account.Data.ID}, client)
	})

	t.Run("Delete returns no error when resource deleted", func(t *testing.T) {
		account := createAccount(t, client)
		got := client.Delete(account.Data.ID, 0)
		assertNoError(t, got)
	})
}

func TestList(t *testing.T) {

	client, err := NewClient(ENDPOINT)
	assertNoError(t, err)

	var accounts []string

	for i := 0; i < 5; i++ {
		account := createAccount(t, client)
		accounts = append(accounts, account.Data.ID)
	}

	t.Run("3 Per page, page 1 returns 3 items", func(t *testing.T) {
		acc, err := client.List(3, 0)
		assertNoError(t, err)

		if len(acc.Data) != 3 {
			t.Fatalf("wanted %d, but got %d", 3, len(acc.Data))
		}

		fmt.Printf("%d", len(acc.Data))
	})

	t.Run("3 Per page, page 2 returns 2 items ", func(t *testing.T) {
		acc, err := client.List(3, 1)
		assertNoError(t, err)

		if len(acc.Data) != 2 {
			t.Fatalf("wanted %d, but got %d", 2, len(acc.Data))
		}

		fmt.Printf("%d", len(acc.Data))
	})

	cleanUpAccounts(t, accounts, client)
}

func cleanUpAccounts(t *testing.T, accounts []string, client *Client) {
	for _, accountID := range accounts {
		err := client.Delete(accountID, 0)
		assertNoError(t, err)
	}
}

func loadTestData(fileName string, model interface{}) (interface{}, error) {
	path := filepath.Join("testdata", fileName)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return model, err
	}

	err = json.Unmarshal(bytes, model)
	if err != nil {
		return model, err
	}

	return model, nil
}

func createAccount(t *testing.T, client *Client) *Entity {
	testReqData, err := loadTestData("valid-create-request.json", &RequestData{})
	assertNoError(t, err)

	requestData := testReqData.(*RequestData)
	requestData.Data.ID = uuid.New().String()

	account, err := client.Create(requestData)
	assertNoError(t, err)

	return account
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Helper()
		t.Fatal("Didn't want error but got one", err)
	}
}

func assertError(t *testing.T, err error) {
	if err == nil {
		t.Helper()
		t.Fatal("wanted error but did not get one", err)
	}
}

func assertErrorContains(t *testing.T, value string, err error) {
	if err == nil {
		t.Helper()
		t.Fatal("wanted error but did not get one", err)
	}

	if !strings.Contains(err.Error(), value) {
		t.Helper()
		t.Fatalf("could not find occurance of %q in %q", value, err.Error())
	}
}
