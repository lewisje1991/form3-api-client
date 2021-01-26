package client

import (
	"errors"
	"reflect"
	"testing"
)

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

	t.Run("Returns client if baseUrl is valid", func(t *testing.T) {
		client, err := NewClient("https://api.form3.com")
		assertNoError(t, err)

		got := client.Accounts.httpService.BaseURL
		want := "https://api.form3.com"

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("wanted %+v, but got %+v", want, got)
		}
	})
}
