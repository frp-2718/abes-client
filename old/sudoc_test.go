package sudocclient

import (
	"errors"
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
	if sudoc.maxAttempts != defaultMaxAttempts {
		t.Errorf("maxAttempts = %d, expected %d", sudoc.maxAttempts, defaultMaxAttempts)
	}
	sudoc.SetMaxAttempts(0)
	if sudoc.maxAttempts != defaultMaxAttempts {
		t.Errorf("maxAttempts = %d, expected %d", sudoc.maxAttempts, defaultMaxAttempts)
	}
	sudoc.SetMaxAttempts(5)
	if sudoc.maxAttempts != 5 {
		t.Errorf("maxAttempts = %d, expected 5", sudoc.maxAttempts)
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

func TestBuildRequest(t *testing.T) {
	var expected []request
	expected = append(expected, request{http.NewRequest("GET", sudocBaseURL+"service/biblio/123456", nil), 0})
	expected[0].Header.Set("Accept", "text/xml")
	expected = append(expected, request{http.NewRequest("GET", sudocBaseURL+"service/biblio/123456", nil), 0})
	expected[1].Header.Set("Accept", "text/json")
	expected = append(expected, request{http.NewRequest("GET", sudocBaseURL+"service/biblio/123456", nil), 0})
	expected[1].Header.Set("Accept", "text/json")
	var tests = []struct {
		serv service
		ppns []string
		fmt  responseFormat
		want request
		err  error
	}{
		{biblio, []string{"123456"}, xml, expected[0], nil},
		{null, []string{"123456"}, xml, nil, errors.New("no service 'null'")},
		{biblio, []string{"123456"}, json, expected[1], nil},
		{null, []string{"123456"}, json, nil, errors.New("no service 'null'")},
		{biblio, []string{"123456"}, null, nil, errors.New("unknown 'null' format")},
		{null, []string{"123456"}, null, nil, erros.New("no service 'null'")},
		{biblio, []string{"123456", "234567"}, xml, nil, errors.New("service 'biblio' does not accept multiple ppns")},
		{null, []string{"123456", "234567"}, xml, nil, errors.New("no service 'null'")},
		{biblio, []string{"123456", "234567"}, json, nil, errors.New("service 'biblio' does not accept multiple ppns")},
		{null, []string{"123456", "234567"}, json, nil, errors.New("no service 'null'")},
		{biblio, []string{"123456", "234567"}, null, nil, errors.New("unknown 'null' format")},
		{null, []string{"123456", "234567"}, null, nil, errors.New("no service 'null'")},
		{biblio, []string{"123456", ""}, xml, expected[0], nil},
		{null, []string{"123456", ""}, xml, nil, errors.New("no service 'null'")},
		{biblio, []string{"123456", ""}, json, expected[1], nil},
		{null, []string{"123456", ""}, null, nil, errors.New("no service 'null'")},
		{biblio, []string{"123456", ""}, null, nil, errors.New("unknown 'null' format")},
		{biblio, []string{""}, xml, nil, errors.New("ppn(s) missing")},
		{null, []string{""}, xml, nil, errors.New("no service 'null'")},
		{biblio, []string{""}, json, nil, errors.New("ppn(s) missing")},
		{null, []string{""}, json, nil, errors.New("no service 'null'")},
		{biblio, []string{""}, null, nil, errors.New("unknown 'null' format")},
		{null, []string{""}, null, nil, errors.New("no service 'null'")},
		{biblio, []string{}, xml, nil, errors.New("ppn(s) missing")},
		{null, []string{}, xml, nil, errors.New("no service 'null'")},
		{biblio, []string{}, json, nil, errors.New("ppn(s) missing")},
		{null, []string{}, json, nil, errors.New("no service 'null'")},
		{biblio, []string{}, null, nil, errors.New("unknownn 'null' format")},
		{null, []string{}, null, nil, errors.New("no service 'null'")},
		{biblio, nil, xml, nil, errors.New("ppn(s) missing")},
		{null, nil, xml, nil, errors.New("no service 'null'")},
		{biblio, nil, json, nil, errors.New("ppn(s) missing")},
		{null, nil, json, nil, errors.New("no service 'null'")},
		{biblio, nil, null, nil, errors.New("unknown 'null' format")},
		{null, nil, null, nil, errors.New("no service 'null'")},
	}
	for i, test := range tests {
		got := sudoc.buildRequest(test.serv, test.ppns, test.fmt)
		if !equalRequests(got, test.want) {
			t.Errorf("test %d: %v, %v, %v got %v want %v", i, test.serv, test.ppns,
				test.fmt, got, test.want)
		}
	}
}

func equalRequests(r1, r2 request) bool {
	return r1.Method == r2.Method &&
		r1.URL.String() == r2.URL.String() &&
		r1.Header.Get("Accept") == r2.Header.Get("Accept") &&
		r1.attempts == r2.attempts
}
