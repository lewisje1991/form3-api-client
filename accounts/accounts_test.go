package accounts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

const ENDPOINT = "http://accountapi:8080/v1"

func TestCreate(t *testing.T) {

	client, err := NewClient(ENDPOINT)
	assertNoError(t, err)

	t.Run("Successful create returns created account", func(t *testing.T) {

		organisationID, accountID, got := createAccount(t, client)

		respData, err := loadTestData("valid-create-response.json", &Entity{})
		assertNoError(t, err)

		want := respData.(*Entity)
		want.Data.OrganisationID = organisationID
		want.Data.ID = accountID
		want.Links.Self = fmt.Sprintf("%s/%s", "/v1/organisation/accounts/", accountID)

		// Set times to match so that we can deep compare... is there a better way?

		now := time.Now()
		got.Data.CreatedOn = now
		want.Data.CreatedOn = now

		got.Data.ModifiedOn = now
		want.Data.ModifiedOn = now

		if reflect.DeepEqual(want, got) {
			t.Fatalf("wanted %+v, but got %+v", got, want)
		}
	})

	t.Run("Invalid request create returns error", func(t *testing.T) {
		testReqData, err := loadTestData("valid-create-request.json", &RequestData{})
		assertNoError(t, err)

		requestData := testReqData.(*RequestData)
		requestData.Data.OrganisationID = ""

		_, got := client.Create(requestData)
		assertError(t, got)
	})
}

func TestFetch(t *testing.T) {

	client, err := NewClient(ENDPOINT)
	assertNoError(t, err)

	t.Run("When account exists, returns entity model", func(t *testing.T) {
		_, accountID, _ := createAccount(t, client)

		want, err := client.Fetch(accountID)
		assertNoError(t, err)

		fmt.Println(want)
		// todo compare response model

	})

	t.Run("When account does not exist, returns error", func(t *testing.T) {
		ID := uuid.New().String()
		_, got := client.Fetch(ID)
		assertError(t, got)

		want := "404 Not Found"

		if got.Error() != want {
			t.Fatalf("wanted %s, but got %s", got, want)
		}
	})
}

func TestDelete(t *testing.T) {
	client, err := NewClient(ENDPOINT)
	assertNoError(t, err)

	t.Run("Delete returns error when resource not found", func(t *testing.T) {
		t.Skip()
		got := client.Delete("invalid resource", 0)
		assertError(t, got)
	})

	t.Run("Delete returns error when version invalid", func(t *testing.T) {
		_, accountID, _ := createAccount(t, client)
		got := client.Delete(accountID, 2)
		assertNoError(t, got)
	})

	t.Run("Delete returns no error when resource deleted", func(t *testing.T) {
		_, accountID, _ := createAccount(t, client)
		got := client.Delete(accountID, 0)
		assertNoError(t, got)
	})

	//todo additional errors bad request with invalid guid.
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

func createAccount(t *testing.T, client *Client) (string, string, *Entity) {
	testReqData, err := loadTestData("valid-create-request.json", &RequestData{})
	assertNoError(t, err)

	requestData := testReqData.(*RequestData)

	orgID := uuid.New().String()
	accountID := uuid.New().String()

	requestData.Data.OrganisationID = orgID
	requestData.Data.ID = accountID

	account, err := client.Create(requestData)
	assertNoError(t, err)

	return orgID, accountID, account
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

// func TestList(t *testing.T) {
// 	client, err := NewClient(ENDPOINT)
// 	if err != nil {
// 		t.Fatal("Didn't want error but got one", err)
// 	}

// 	acc, err := client.List(1, 2)

// 	if err != nil {
// 		t.Fatal("Didn't want error but got one", err)
// 	}

// 	fmt.Printf("%+v\n", acc)
// }
