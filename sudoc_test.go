package sudoc

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

// Ad-hoc experimentation.
func assertNotNil(object interface{}, t *testing.T, label string) {
	if object == nil || (reflect.ValueOf(object).Kind() == reflect.Ptr && reflect.ValueOf(object).IsNil()) {
		t.Errorf("%s is nil", label)
	}
}

func TestSudoc(t *testing.T) {
	sudoc := New()
	assertNotNil(sudoc, t, "New()")
	if sudoc != nil {
		assertNotNil(sudoc.client, t, ".client")
		assertNotNil(sudoc.Bibs, t, ".Bibs")
	}
	myClient := &http.Client{
		Timeout: time.Second * 5,
	}
	sudoc.SetHTTPClient(myClient)
	assertNotNil(sudoc.client, t, ".client")
	if sudoc.client != myClient {
		t.Errorf("expected 'myClient' HTTP client, found %v", sudoc.client)
	}
	sudoc.SetHTTPClient(nil)
	assertNotNil(sudoc.client, t, ".client")
}

func TestIsValidPPN(t *testing.T) {
	var tests = []struct {
		ppn  string
		want bool
	}{
		{"123456789", true},
		{"90873728x", true},
		{"73826351X", true},
		{"000000000", true},
		{"", false},
		{"1234567890", false},
		{"a11222333", false},
		{"11122233%", false},
		{"         ", false},
		{"32O25O11O", false},
	}
	for _, test := range tests {
		if got := IsValidPPN(test.ppn); got != test.want {
			t.Errorf("%s returned %v ; want %v", test.ppn, got, test.want)
		}
	}
}
