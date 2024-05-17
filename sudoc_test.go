package sudoc

import (
	"net/http"
	"slices"
	"testing"
)

func TestNewSudoc(t *testing.T) {
	http_client := &http.Client{}
	sudoc := NewSudoc(nil)
	if sudoc == nil {
		t.Fatal("sudoc was not initialized")
	}
	if sudoc.client == nil {
		t.Fatal("sudoc HTTP client was not set")
	}
	sudoc = NewSudoc(http_client)
	if sudoc.client != http_client {
		t.Fatalf("sudoc HTTP client was not set: have %v, want %v", sudoc.client, http_client)
	}
}

// func TestDo(t *testing.T) {
// 	server := httptest.NewServer(http.HandlerFunc(
// 		func(w http.ResponseWriter, r *http.Request) {
// 			if r.URL.Path != "ppn" {
// 				t.Errorf("Expected to request /ppn, got %s", r.URL.Path)
// 			}
// 			w.WriteHeader(http.StatusOK)
// 			w.Write([]byte(`{"value":"fixed"}`))
// 		}))
// 	defer server.Close()

// 	s := NewSudoc(nil)
// 	value, _ := s.do(server.URL)
// 	body, _ := io.ReadAll(value.Body)
// 	value.Body.Close()
// 	if string(body) != "fixed" {
// 		t.Error("error")
// 	}
// }

func TestConcatPPNs(t *testing.T) {
	s := NewSudoc(nil)
	max_ppns := 3

	var tests = []struct {
		in       []string
		expected []string
	}{
		{[]string{}, []string{}},
		{[]string{"a"}, []string{"a"}},
		{[]string{"a", "b", "c"}, []string{"a,b,c"}},
		{[]string{"a", "b", "c", "d"}, []string{"a,b,c", "d"}},
		{[]string{"a", "b", "c", "d", "e", "f", "g", "h", "i"},
			[]string{"a,b,c", "d,e,f", "g,h,i"}},
	}

	for _, test := range tests {
		got := s.concatPPNs(test.in, max_ppns)
		if slices.Compare(got, test.expected) != 0 {
			t.Errorf("got %#v, want %#v", got, test.expected)
		}
	}
}
