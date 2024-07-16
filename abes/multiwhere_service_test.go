package abes

import (
	"encoding/xml"
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
	ac := NewAbesClient(nil)
	ppns := []string{"144089661", "154923206"}

	result := ac.Multiwhere.GetMultiLocations(ppns, 0)
}

// // GetMultiLocations returns a map associating each valid PPN to its locations,
// // represented by a list of libraries.
// func (ms *MultiwhereService) GetMultiLocations(ppns []string, max_ppns int) map[string][]Library {
// 	ppnStrings := ms.concatPPNs(ppns, max_ppns)
// 	result := make(map[string][]Library)

// 	for _, p := range ppnStrings {
// 		// TODO: handle do() errors
// 		res, _ := ms.client.Get(ms.buildURL(ms.endpoint, p))
// 		body, _ := io.ReadAll(res.Body)
// 		res.Body.Close()

// 		var sr serviceResult
// 		xml.Unmarshal(body, &sr)

// 		for _, query := range sr.Queries {
// 			for _, library := range query.Result.Libraries {
// 				result[query.PPN] = append(result[query.PPN], library)
// 			}
// 		}
// 	}
// 	return result
// }
// // Library represents a location.
// type Library struct {
// 	XMLName   xml.Name `xml:"library"`
// 	RCR       string   `xml:"rcr"`
// 	Shortname string   `xml:"shortname"`
// 	Latitude  float64  `xml:"latitude"`
// 	Longitude float64  `xml:"longitude"`
// }

// type result struct {
// 	XMLName   xml.Name  `xml:"result"`
// 	Libraries []Library `xml:"library"`
// }

// type query struct {
// 	XMLName xml.Name `xml:"query"`
// 	PPN     string   `xml:"ppn"`
// 	Result  result   `xml:"result"`
// }

// type serviceResult struct {
// 	XMLName xml.Name `xml:"sudoc"`
// 	Queries []query  `xml:"query"`
// }

// // GetLocations returns the list of the locations of the given PPN.
// func (ms *MultiwhereService) GetLocations(ppn string) []Library {
// 	ppnList := []string{ppn}

// 	res := ms.GetMultiLocations(ppnList, 1)
// 	return res[ppn]
// }

// // GetMultiLocationsWithErrors returns a map associating each valid PPN to its
// // locations - represented by a list of libraries - and a list of the  invalid
// // PPNs among the requested ones.
// func (ms *MultiwhereService) GetMultiLocationsWithErrors(ppns []string, max_ppns int) (map[string][]Library, []string) {
// 	result := ms.GetMultiLocations(ppns, max_ppns)
// 	var invalid_ppns []string
// 	var found_ppns []string
// 	for k := range result {
// 		found_ppns = append(found_ppns, k)
// 	}
// 	for _, ppn := range ppns {
// 		if !slices.Contains(found_ppns, ppn) {
// 			invalid_ppns = append(invalid_ppns, ppn)
// 		}
// 	}
// 	return result, invalid_ppns
// }
