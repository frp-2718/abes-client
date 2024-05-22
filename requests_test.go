package sudoc

import (
	"io"
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
				w.Write([]byte(`
                <?xml version="1.0" encoding="UTF-8" ?>
                    <sudoc service="multiwhere">
                    <error>Found a null xml in result : values={ppn=notfound}, query=select autorites.MULTIWHERE(#ppn#) from dual</error>
                    </sudoc>
                `))
			} else if r.URL.Path == "/empty" {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`
                <?xml version="1.0" encoding="UTF-8" ?>
                    <sudoc service="multiwhere">
                    <error>Invalid char in query string, values={}, query=select autorites.MULTIWHERE(#ppn#) from dual</error>
                    </sudoc>
                `))
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}))
	defer server.Close()

	s := New(nil)
	value, _ := s.do(server.URL + "/ppn")
	body, _ := io.ReadAll(value.Body)
	value.Body.Close()
	if string(body) != "fixed" {
		t.Error("error")
	}
}
