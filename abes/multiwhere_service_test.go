package abes

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLibraryString(t *testing.T) {
	xml := xml.Name{}
	tests := []struct {
		input    Library
		expected string
	}{
		{Library{xml, "RCR1", "Library1", 1.1, 1.2}, "[RCR1] Library1 (1.1, 1.2)"},
		{Library{xml, "", "Library1", 1.1, 1.2}, "[] Library1 (1.1, 1.2)"},
		{Library{xml, "", "", 1.1, 1.2}, "[]  (1.1, 1.2)"},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			result := test.input.String()
			if result != test.expected {
				t.Errorf("Library %#v: got %#v, want %#v", test.input, result, test.expected)
			}
		})
	}
}

func TestNewMultiwhereService(t *testing.T) {
	ms := newMultiwhereService(http.DefaultClient, "ms_endpoint")
	assert := assert.New(t)
	assert.Equal(ms.endpoint, "ms_endpoint", "Endpoint mismatch")
	assert.Equal(ms.max_ppns, MAX_MULTIWHERE_PPNS)
	assert.Same(ms.client, http.DefaultClient)
}

func TestConcatPPNs(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{[]string{}, []string{}},
		{[]string{"PPN1"}, []string{"PPN1"}},
		{[]string{"PPN1", "PPN2"}, []string{"PPN1,PPN2"}},
		{[]string{"PPN1", "PPN2", "PPN3"}, []string{"PPN1,PPN2", "PPN3"}},
		{[]string{"PPN1", "PPN2", "PPN3", "PPN4"}, []string{"PPN1,PPN2", "PPN3,PPN4"}},
		{[]string{"PPN1", "PPN2", "PPN3", "PPN4", "PPN5"}, []string{"PPN1,PPN2", "PPN3,PPN4", "PPN5"}},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			result := concatValues(test.input, 2, MAX_MULTIWHERE_PPNS)
			if !slices.Equal(result, test.expected) {
				t.Errorf("PPNs %v: got %v, want %v", test.input, result, test.expected)
			}
		})
	}
}

func TestBuildURL(t *testing.T) {
	tests := []struct {
		base, path string
		expected   string
	}{
		{"", "", "/"},
		{"base", "", "base/"},
		{"", "path", "/path"},
		{"base", "path", "base/path"},
		{"base/", "/path", "base/path"},
		{"base/", "path", "base/path"},
		{"base", "/path", "base/path"},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			result := joinURLPath(test.base, test.path)
			if result != test.expected {
				t.Errorf("(%s,%s): got %s, want %s", test.base, test.path, result, test.expected)
			}
		})
	}
}

func TestGetMultiLocations(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/services/multiwhere/notfound":
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `<sudoc service="multiwhere"><error>Found a null xml in result : values={ppn=notfound}</error></sudoc>`)
			return
		case r.URL.Path == "/services/multiwhere/servererror":
			http.Error(w, "boom", http.StatusInternalServerError)
			return
		case r.URL.Path == "/services/multiwhere/111111111,154923206":
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `<sudoc service="multiwhere"><query>
                <ppn>154923206</ppn>
                <result>
                <library>
                <rcr>751052105</rcr>
                <shortname>PARIS-BIS, Fonds ge\u0301ne\u0301ral</shortname>
                <latitude>48.8492618</latitude>
                <longitude>2.3433311</longitude>
                </library>
                <library>
                <rcr>751052116</rcr>
                <shortname>PARIS-Bib. Sainte Genevie\u0300ve</shortname>
                <latitude>48.8467139</latitude>
                <longitude>2.3463854</longitude>
                </library>
                </result>
                </query>
                </sudoc>`)
		default:
			fmt.Println(r.URL.Path)
		}
	}
	ts := newTestServer(handler)
	defer ts.close()

	assert := assert.New(t)
	ac := NewAbesClient(ts.client)

	result, err := ac.Multiwhere.GetMultiLocations([]string{"notfound"}, 0)
	var nf *NotFoundError
	assert.ErrorAs(err, &nf)
	assert.True(NotFound(err))
	assert.Nil(result)

	result, err = ac.Multiwhere.GetMultiLocations([]string{"servererror"}, 0)
	var apiErr *APIError
	assert.ErrorAs(err, &apiErr)
	assert.Equal(http.StatusInternalServerError, apiErr.StatusCode)
	assert.Nil(result)

	ppns := []string{"111111111", "154923206"}
	xml := xml.Name{Local: "library"}
	result, err = ac.Multiwhere.GetMultiLocations(ppns, 0)
	assert.NoError(err, "ppns should return a 200 response, got an error")
	expected := map[string][]Library{
		"154923206": {
			Library{xml, "751052105", "PARIS-BIS, Fonds ge\u0301ne\u0301ral", 48.8492618, 2.3433311},
			Library{xml, "751052116", "PARIS-Bib. Sainte Genevie\u0300ve", 48.8467139, 2.3463854},
		},
	}
	assert.Equal(expected, result)

	ts.simulateNetworkFailure(true)
	result, err = ac.Multiwhere.GetMultiLocations([]string{"network_error"}, 0)
	var ne *NetworkError
	assert.ErrorAs(err, &ne)
	assert.Nil(result)
}

func TestGetLocations(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/services/multiwhere/notfound":
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `<sudoc service="multiwhere"><error>Found a null xml in result : values={ppn=notfound}</error></sudoc>`)
			return
		case r.URL.Path == "/services/multiwhere/servererror":
			http.Error(w, "boom", http.StatusInternalServerError)
			return
		case r.URL.Path == "/services/multiwhere/154923206":
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `<sudoc service="multiwhere"><query>
                <ppn>154923206</ppn>
                <result>
                <library>
                <rcr>751052105</rcr>
                <shortname>PARIS-BIS, Fonds ge\u0301ne\u0301ral</shortname>
                <latitude>48.8492618</latitude>
                <longitude>2.3433311</longitude>
                </library>
                <library>
                <rcr>751052116</rcr>
                <shortname>PARIS-Bib. Sainte Genevie\u0300ve</shortname>
                <latitude>48.8467139</latitude>
                <longitude>2.3463854</longitude>
                </library>
                </result>
                </query>
                </sudoc>`)
		default:
			fmt.Println(r.URL.Path)
		}
	}
	ts := newTestServer(handler)
	defer ts.close()

	assert := assert.New(t)
	ac := NewAbesClient(ts.client)

	result, err := ac.Multiwhere.GetLocations("notfound")
	var nf *NotFoundError
	assert.ErrorAs(err, &nf)
	assert.Nil(result, "notfound should return a nil slice")

	ppn := "154923206"
	xml := xml.Name{Local: "library"}
	result, err = ac.Multiwhere.GetLocations(ppn)
	assert.NoError(err, "ppns should return a 200 response, got an error")
	expected := []Library{
		{xml, "751052105", "PARIS-BIS, Fonds ge\u0301ne\u0301ral", 48.8492618, 2.3433311},
		{xml, "751052116", "PARIS-Bib. Sainte Genevie\u0300ve", 48.8467139, 2.3463854},
	}
	assert.Equal(expected, result)

	ts.simulateNetworkFailure(true)
	result, err = ac.Multiwhere.GetLocations("network_error")
	var ne *NetworkError
	assert.ErrorAs(err, &ne)
}

func TestGetMultiLocationsWithErrors(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/services/multiwhere/notfound":
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `<sudoc service="multiwhere"><error>Found a null xml in result : values={ppn=notfound}</error></sudoc>`)
			return
		case r.URL.Path == "/services/multiwhere/servererror":
			http.Error(w, "boom", http.StatusInternalServerError)
			return
		case r.URL.Path == "/services/multiwhere/111111111,154923206" ||
			r.URL.Path == "/services/multiwhere/154923206":
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `<sudoc service="multiwhere"><query>
                <ppn>154923206</ppn>
                <result>
                <library>
                <rcr>751052105</rcr>
                <shortname>PARIS-BIS, Fonds ge\u0301ne\u0301ral</shortname>
                <latitude>48.8492618</latitude>
                <longitude>2.3433311</longitude>
                </library>
                <library>
                <rcr>751052116</rcr>
                <shortname>PARIS-Bib. Sainte Genevie\u0300ve</shortname>
                <latitude>48.8467139</latitude>
                <longitude>2.3463854</longitude>
                </library>
                </result>
                </query>
                </sudoc>`)
		default:
			fmt.Println(r.URL.Path)
		}
	}
	ts := newTestServer(handler)
	defer ts.close()

	assert := assert.New(t)
	ac := NewAbesClient(ts.client)

	result, wrong, err := ac.Multiwhere.GetMultiLocationsWithErrors([]string{"notfound"}, 0)
	var nf *NotFoundError
	assert.ErrorAs(err, &nf)
	assert.Nil(result)
	assert.Nil(wrong)

	ppns := []string{"154923206"}
	xml := xml.Name{Local: "library"}
	result, wrong, err = ac.Multiwhere.GetMultiLocationsWithErrors(ppns, 0)
	assert.NoError(err, "ppns should return a 200 response, got an error")
	expected := map[string][]Library{
		"154923206": {
			Library{xml, "751052105", "PARIS-BIS, Fonds ge\u0301ne\u0301ral", 48.8492618, 2.3433311},
			Library{xml, "751052116", "PARIS-Bib. Sainte Genevie\u0300ve", 48.8467139, 2.3463854},
		},
	}
	assert.Equal(expected, result)
	assert.Equal([]string{}, wrong)

	ppns = []string{"111111111", "154923206"}
	result, wrong, err = ac.Multiwhere.GetMultiLocationsWithErrors(ppns, 0)
	assert.NoError(err, "ppns should return a 200 response, got an error")
	assert.Equal(expected, result)
	assert.Equal([]string{"111111111"}, wrong)

	ts.simulateNetworkFailure(true)
	result, wrong, err = ac.Multiwhere.GetMultiLocationsWithErrors([]string{"network_error"}, 0)
	var ne *NetworkError
	assert.ErrorAs(err, &ne)
}
