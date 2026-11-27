package abes

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRecord(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/172258367.xml":
			fmt.Fprint(w, `<record><leader>     cam0 22        450 </leader></record>`)
		case "/999999999.xml":
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><error>Les données bibliographiques sont indéfinies <ppn>999999999</ppn>`)
		case "/error200.xml":
			fmt.Fprint(w, `<error>Les données bibliographiques sont indéfinies</error>`)
		case "/forbidden.xml":
			http.Error(w, "forbidden", http.StatusForbidden)
		}
	}
	ts := newTestServer(handler)
	defer ts.close()
	ac := NewAbesClient(ts.client)
	assert := assert.New(t)

	record, err := ac.UnimarcXML.GetRecord("172258367")
	assert.NoError(err)
	assert.Equal("     cam0 22        450 ", record.Leader)

	record, err = ac.UnimarcXML.GetRecord("999999999")
	assert.Nil(record)
	var nf *NotFoundError
	assert.ErrorAs(err, &nf)
	assert.Equal("Les données bibliographiques sont indéfinies", nf.Message)
	assert.True(NotFound(err))

	_, err = ac.UnimarcXML.GetRecord("error200")
	assert.ErrorAs(err, &nf)

	_, err = ac.UnimarcXML.GetRecord("forbidden")
	var apiErr *APIError
	assert.ErrorAs(err, &apiErr)
	assert.Equal(http.StatusForbidden, apiErr.StatusCode)

	ts.simulateNetworkFailure(true)
	_, err = ac.UnimarcXML.GetRecord("network_error")
	var ne *NetworkError
	assert.ErrorAs(err, &ne)
}
