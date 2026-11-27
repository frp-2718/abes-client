package abes

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIsbn2ppnService(t *testing.T) {
	is := newIsbn2ppnService(http.DefaultClient, "isbn_endpoint")

	assert.Equal(t, "isbn_endpoint", is.endpoint)
	assert.Equal(t, MAX_ISBN2PPN_ISBN, is.max_isbns)
	assert.Same(t, http.DefaultClient, is.client)
}

func TestISBNResultString(t *testing.T) {
	tests := []struct {
		name string
		in   ISBNResult
		want string
	}{
		{"both lists", ISBNResult{ISBN: "978123", Holdings: []string{"PPN1", "PPN2"}, NoHolding: []string{"PPN3"}}, "978123 -> [PPN1,PPN2][PPN3]"},
		{"empty lists", ISBNResult{ISBN: "978123"}, "978123 -> [][]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.in.String())
		})
	}
}

func TestGetMultiPPNs(t *testing.T) {
	var requestedPaths []string
	handler := func(w http.ResponseWriter, r *http.Request) {
		requestedPaths = append(requestedPaths, r.URL.Path)
		switch r.URL.Path {
		case "/services/isbn2ppn/missing":
			http.Error(w, "not found", http.StatusNotFound)
		case "/services/isbn2ppn/9780000000001,9780000000002":
			w.Header().Set("Content-Type", "application/xml")
			fmt.Fprint(w, `<sudoc service="isbn2ppn">
                <query><isbn>9780000000001</isbn><result><ppn>PPN1</ppn><ppn>PPN2</ppn></result><resultNoHolding><ppn>PPN3</ppn></resultNoHolding></query>
                <query><isbn>9780000000002</isbn><result><ppn>PPN4</ppn></result></query>
            </sudoc>`)
		default:
			http.NotFound(w, r)
		}
	}
	ts := newTestServer(handler)
	defer ts.close()
	is := newIsbn2ppnService(ts.client, ts.server.URL+"/services/isbn2ppn/")

	result, err := is.GetMultiPPNs([]string{"missing"}, 0)
	var nf *NotFoundError
	assert.ErrorAs(t, err, &nf)
	assert.True(t, NotFound(err))
	assert.Nil(t, result)

	result, err = is.GetMultiPPNs([]string{"9780000000001", "9780000000002"}, 0)
	assert.NoError(t, err)
	assert.Equal(t, map[string]ISBNResult{
		"9780000000001": {ISBN: "9780000000001", Holdings: []string{"PPN1", "PPN2"}, NoHolding: []string{"PPN3"}},
		"9780000000002": {ISBN: "9780000000002", Holdings: []string{"PPN4"}},
	}, result)

	assert.Equal(t, []string{
		"/services/isbn2ppn/missing",
		"/services/isbn2ppn/9780000000001,9780000000002",
	}, requestedPaths)

	ts.simulateNetworkFailure(true)
	result, err = is.GetMultiPPNs([]string{"network_error"}, 0)
	assert.Nil(t, result)
	var ne *NetworkError
	assert.ErrorAs(t, err, &ne)
}

func TestGetMultiPPNsSplitsRequests(t *testing.T) {
	var requestedPaths []string
	handler := func(w http.ResponseWriter, r *http.Request) {
		requestedPaths = append(requestedPaths, r.URL.Path)
		fmt.Fprint(w, `<sudoc service="isbn2ppn"/>`)
	}
	ts := newTestServer(handler)
	defer ts.close()
	is := newIsbn2ppnService(ts.client, ts.server.URL+"/services/isbn2ppn/")

	_, err := is.GetMultiPPNs([]string{"a", "b", "c", "d", "e"}, 2)
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"/services/isbn2ppn/a,b",
		"/services/isbn2ppn/c,d",
		"/services/isbn2ppn/e",
	}, requestedPaths)
}

func TestGetPPNs(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/services/isbn2ppn/978-0134190440" {
			// The Sudoc service can return the normalized ISBN-10 for an ISBN-13 query.
			fmt.Fprint(w, `<sudoc service="isbn2ppn"><query><isbn>0134190440</isbn><result><ppn>191539082</ppn></result></query></sudoc>`)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	}
	ts := newTestServer(handler)
	defer ts.close()
	is := newIsbn2ppnService(ts.client, ts.server.URL+"/services/isbn2ppn/")

	result, err := is.GetPPNs("978-0134190440")
	assert.NoError(t, err)
	assert.Equal(t, ISBNResult{ISBN: "0134190440", Holdings: []string{"191539082"}}, result)

	result, err = is.GetPPNs("missing")
	var nf *NotFoundError
	assert.ErrorAs(t, err, &nf)
	assert.Equal(t, ISBNResult{}, result)

	ts.simulateNetworkFailure(true)
	result, err = is.GetPPNs("network_error")
	var ne *NetworkError
	assert.ErrorAs(t, err, &ne)
	assert.Equal(t, ISBNResult{}, result)
}
