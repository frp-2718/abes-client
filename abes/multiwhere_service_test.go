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
	ac := NewAbesClient(nil)

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
			result := ac.Multiwhere.concatPPNs(test.input, 2)
			if !slices.Equal(result, test.expected) {
				t.Errorf("PPNs %v: got %v, want %v", test.input, result, test.expected)
			}
		})
	}
}

func TestBuildURL(t *testing.T) {
	ac := NewAbesClient(nil)
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
			result := ac.Multiwhere.buildURL(test.base, test.path)
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
			http.Error(w, "not found", http.StatusInternalServerError)
			return
		case r.URL.Path == "/services/multiwhere/111111111,154923206":
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `<sudoc service="multiwhere"><query>
                <ppn>154923206</ppn>
                <result>
                <library>
                <rcr>751052105</rcr>
                <shortname>PARIS-BIS, Fonds général</shortname>
                <latitude>48.8492618</latitude>
                <longitude>2.3433311</longitude>
                </library>
                <library>
                <rcr>751052116</rcr>
                <shortname>PARIS-Bib. Sainte Geneviève</shortname>
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
	assert.NoError(err, "notfound should return a 500 response, got an error")
	assert.Equal(map[string][]Library{}, result)

	ppns := []string{"111111111", "154923206"}
	xml := xml.Name{Local: "library"}
	result, err = ac.Multiwhere.GetMultiLocations(ppns, 0)
	assert.NoError(err, "ppns should return a 200 response, got an error")
	expected := map[string][]Library{
		"154923206": {
			Library{xml, "751052105", "PARIS-BIS, Fonds général", 48.8492618, 2.3433311},
			Library{xml, "751052116", "PARIS-Bib. Sainte Geneviève", 48.8467139, 2.3463854},
		},
	}
	assert.Equal(expected, result)

	ts.simulateNetworkFailure(true)
	result, err = ac.Multiwhere.GetMultiLocations([]string{"network_error"}, 0)
	assert.NotNil(err)
	assert.IsType(&NetworkError{}, err)
}

func TestGetLocations(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/services/multiwhere/notfound":
			http.Error(w, "not found", http.StatusInternalServerError)
			return
		case r.URL.Path == "/services/multiwhere/154923206":
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `<sudoc service="multiwhere"><query>
                <ppn>154923206</ppn>
                <result>
                <library>
                <rcr>751052105</rcr>
                <shortname>PARIS-BIS, Fonds général</shortname>
                <latitude>48.8492618</latitude>
                <longitude>2.3433311</longitude>
                </library>
                <library>
                <rcr>751052116</rcr>
                <shortname>PARIS-Bib. Sainte Geneviève</shortname>
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
	assert.NoError(err, "notfound should return a 500 response, got an error")
	assert.Nil(result, "notfound should return a nil slice")

	ppn := "154923206"
	xml := xml.Name{Local: "library"}
	result, err = ac.Multiwhere.GetLocations(ppn)
	assert.NoError(err, "ppns should return a 200 response, got an error")
	expected := []Library{
		{xml, "751052105", "PARIS-BIS, Fonds général", 48.8492618, 2.3433311},
		{xml, "751052116", "PARIS-Bib. Sainte Geneviève", 48.8467139, 2.3463854},
	}
	assert.Equal(expected, result)

	ts.simulateNetworkFailure(true)
	result, err = ac.Multiwhere.GetLocations("network_error")
	assert.NotNil(err)
	assert.IsType(&NetworkError{}, err)
}

func TestGetMultiLocationsWithErrors(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/services/multiwhere/notfound":
			http.Error(w, "not found", http.StatusInternalServerError)
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
                <shortname>PARIS-BIS, Fonds général</shortname>
                <latitude>48.8492618</latitude>
                <longitude>2.3433311</longitude>
                </library>
                <library>
                <rcr>751052116</rcr>
                <shortname>PARIS-Bib. Sainte Geneviève</shortname>
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
	assert.NoError(err, "notfound should return a 500 response, got an error")
	assert.Equal(map[string][]Library{}, result)
	assert.Equal([]string{"notfound"}, wrong)

	ppns := []string{"154923206"}
	xml := xml.Name{Local: "library"}
	result, wrong, err = ac.Multiwhere.GetMultiLocationsWithErrors(ppns, 0)
	assert.NoError(err, "ppns should return a 200 response, got an error")
	expected := map[string][]Library{
		"154923206": {
			Library{xml, "751052105", "PARIS-BIS, Fonds général", 48.8492618, 2.3433311},
			Library{xml, "751052116", "PARIS-Bib. Sainte Geneviève", 48.8467139, 2.3463854},
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
	assert.NotNil(err)
	assert.IsType(&NetworkError{}, err)
}
