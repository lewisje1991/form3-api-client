package client

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {

	t.Run("Returns error if baseUrl is empty", func(t *testing.T) {
		_, got := NewClient("")
		if got == nil {
			t.Fatal("Wanted an error but did not get one")
		}

		want := ErrBaseURLEmpty

		if !errors.Is(got, want) {
			t.Fatalf("wanted %s, but got %s", want, got)
		}
	})

	t.Run("Returns error if baseUrl is invalid", func(t *testing.T) {
		_, err := NewClient("testssdfsdsd")
		if err == nil {
			t.Fatal("Wanted an error but did not get one")
		}

		if !errors.Is(err, ErrBaseURLInvalid) {
			t.Fatalf("wanted %s, but got %s", ErrBaseURLInvalid, err)
		}
	})

	t.Run("Returns client if validation passes", func(t *testing.T) {
		got, err := NewClient("https://api.form3.com")
		if err != nil {
			t.Fatal("Didn't want error but got one", err)
		}

		want := &Client{
			BaseURL: "https://api.form3.com",
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("wanted %+v, but got %+v", want, got)
		}
	})

}
