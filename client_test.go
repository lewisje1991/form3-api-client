package client

import "testing"

func TestClientInitialistion(t *testing.T) {

	t.Run("Returns error if host is empty", func(t *testing.T) {
		_, err := NewClient("")
		if err == nil {
			t.Fatal("Wanted an error but did not get one")
		}

		
	})

}
