package sudoc

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

func TestConcatPPNs(t *testing.T) {
	s := New(nil)
	max_ppns := 3

	var tests = []struct {
		in       []string
		expected []string
	}{
		{[]string{}, []string{}},
		{[]string{"a"}, []string{"a"}},
		{[]string{"a", "b", "c"}, []string{"a,b,c"}},
		{[]string{"a", "b", "c", "d"}, []string{"a,b,c", "d"}},
		{[]string{"a", "b", "c", "d", "e", "f", "g", "h"},
			[]string{"a,b,c", "d,e,f", "g,h"}},
	}

	for _, test := range tests {
		got := s.concatPPNs(test.in, max_ppns)
		if slices.Compare(got, test.expected) != 0 {
			t.Errorf("got %#v, want %#v", got, test.expected)
		}
	}
}

func TestDo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/ok" {
				w.WriteHeader(http.StatusOK)
			} else if r.URL.Path == "/not_found" {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}))
	defer server.Close()

	tests := []struct {
		name string
		url  string
		want int
	}{
		{"ok", "/ok", http.StatusOK},
		{"not_found", "/not_found", http.StatusInternalServerError},
		{"unexpected", "/unexpected", http.StatusNotFound},
	}

	s := New(nil)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _ := s.do(server.URL + test.url)
			if test.want != got.StatusCode {
				t.Errorf("want %v status code, got %v", test.want, got.StatusCode)
			}
		})
	}
}

func TestBuildURL(t *testing.T) {
	tests := []struct {
		name string
		base string
		path string
		want string
	}{
		{"empty base and path", "", "", "/"},
		{"empty base", "", "/path", "/path"},
		{"empty path", "base/", "", "base/"},
		{"whithout separator", "base", "path", "base/path"},
		{"with base separator", "base/", "path", "base/path"},
		{"with path separator", "base", "/path", "base/path"},
		{"with both separator", "base/", "/path", "base/path"},
	}

	s := New(nil)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := s.buildURL(test.base, test.path)
			if test.want != got {
				t.Errorf("want %v, got %v", test.want, got)
			}
		})
	}
}
