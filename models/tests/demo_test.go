package tests

import "testing"

func TestListStatus(t *testing.T) {
	err := ListStatus("")
	if err != nil {
		t.Fatal(err)
	}
}
