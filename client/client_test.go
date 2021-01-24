package client

import (
	"errors"
	"reflect"
	"testing"
)

func TestSetBaseUrl(t *testing.T) {

	t.Run("Returns error if baseUrl is empty", func(t *testing.T) {
		c := Client{}
		got := c.SetBaseUrl("")
		if got == nil {
			t.Fatal("Wanted an error but did not get one")
		}

		want := ErrBaseURLEmpty

		if !errors.Is(got, want) {
			t.Fatalf("wanted %s, but got %s", want, got)
		}
	})

	t.Run("Returns error if baseUrl is invalid", func(t *testing.T) {
		c := Client{}
		got := c.SetBaseUrl("testssdfsdsd")
		if got == nil {
			t.Fatal("Wanted an error but did not get one")
		}

		want := ErrBaseURLInvalid

		if !errors.Is(got, want) {
			t.Fatalf("wanted %s, but got %s", ErrBaseURLInvalid, want)
		}
	})

	t.Run("Returns client if validation passes", func(t *testing.T) {
		c := Client{}
		err := c.SetBaseUrl("https://api.form3.com")
		if err != nil {
			t.Fatal("Didn't want error but got one", err)
		}

		got := c.baseURL
		want := "https://api.form3.com"

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("wanted %+v, but got %+v", want, got)
		}
	})
}
