package form3

import (
	"encoding/json"
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

func TestCreate(t *testing.T) {
	client, err := NewClient(ENDPOINT)
	assertNoError(t, err)

	t.Run("Successful create returns created account", func(t *testing.T) {

		got := createAccount(t, client)

		respData, err := loadTestData("valid-create-response.json", &Account{})
		assertNoError(t, err)

		want := respData.(*Account)

		compareAccountModels(t, got, want)

		cleanUpAccounts(t, []string{got.Data.ID}, client)
	})

	t.Run("Validation errors are returned", func(t *testing.T) {
		testReqData, err := loadTestData("valid-create-request.json", &AccountCreateRequest{})
		assertNoError(t, err)

		requestData := testReqData.(*AccountCreateRequest)
		requestData.Data.OrganisationID = ""

		_, got := client.Accounts.Create(requestData)
		assertErrorContains(t, "organisation_id in body is required", got)
	})
}

func TestFetch(t *testing.T) {
	client, err := NewClient(ENDPOINT)
	assertNoError(t, err)

	t.Run("When account existsaaaaa, returns account model", func(t *testing.T) {
		account := createAccount(t, client)

		got, err := client.Accounts.Fetch(account.Data.ID)
		assertNoError(t, err)

		respData, err := loadTestData("valid-create-response.json", &Account{})
		assertNoError(t, err)

		want := respData.(*Account)

		compareAccountModels(t, got, want)

		cleanUpAccounts(t, []string{account.Data.ID}, client)
	})

	t.Run("When account does not exist, returns error", func(t *testing.T) {
		ID := uuid.New().String()
		_, got := client.Accounts.Fetch(ID)
		assertErrorContains(t, "status code 404", got)
	})
}

func TestDelete(t *testing.T) {
	client, err := NewClient(ENDPOINT)
	assertNoError(t, err)

	t.Run("Delete returns error when resource not found", func(t *testing.T) {
		got := client.Accounts.Delete("invalid resource uuid invalid", 0)
		assertErrorContains(t, "status code 400", got)
	})

	t.Run("Delete returns error when version invalid", func(t *testing.T) {
		account := createAccount(t, client)
		got := client.Accounts.Delete(account.Data.ID, 2)
		assertErrorContains(t, "status code 404", got)
		cleanUpAccounts(t, []string{account.Data.ID}, client)
	})

	t.Run("Delete returns no error when resource deleted", func(t *testing.T) {
		account := createAccount(t, client)
		got := client.Accounts.Delete(account.Data.ID, 0)
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
		got, err := client.Accounts.List(3, 0)
		assertNoError(t, err)

		if len(got.Data) != 3 {
			t.Fatalf("wanted %d, but got %d", 3, len(got.Data))
		}
	})

	t.Run("3 Per page, page 2 returns 2 items ", func(t *testing.T) {
		got, err := client.Accounts.List(3, 1)
		assertNoError(t, err)

		if len(got.Data) != 2 {
			t.Fatalf("wanted %d, but got %d", 2, len(got.Data))
		}
	})

	cleanUpAccounts(t, accounts, client)
}

func cleanUpAccounts(t *testing.T, accounts []string, client *Client) {
	for _, accountID := range accounts {
		err := client.Accounts.Delete(accountID, 0)
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

func createAccount(t *testing.T, client *Client) *Account {
	testReqData, err := loadTestData("valid-create-request.json", &AccountCreateRequest{})
	assertNoError(t, err)

	requestData := testReqData.(*AccountCreateRequest)
	requestData.Data.ID = uuid.New().String()

	account, err := client.Accounts.Create(requestData)
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

func compareAccountModels(t *testing.T, got, want *Account) {
	want.Data.ID = got.Data.ID
	want.Links.Self = fmt.Sprintf("%s/%s", "/v1/organisation/accounts/", got.Data.ID)

	// Set times to match so that we can deep compare... is there a better way?
	now := time.Now()
	got.Data.CreatedOn = now
	got.Data.ModifiedOn = now
	want.Data.ModifiedOn = now
	want.Data.CreatedOn = now

	if reflect.DeepEqual(want, got) {
		t.Fatalf("wanted %+v, but got %+v", got, want)
	}
}
