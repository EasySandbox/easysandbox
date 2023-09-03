package virtproviders_test

import (
	"testing"
)

func TestGetProviders(t *testing.T) {
	providers, err := GetProviders(".")
	if err != nil {
		t.Fatal(err)
	}

	testFunc, err := providers["test.so"].Lookup("TestString")
	if err != nil {
		t.Fatal(err)
	}

	testStringFunc := testFunc.(func() string)
	if got := testStringFunc(); got != "Test string" {
		t.Errorf("expected 'Test string', got %q", got)
	}
}
