package abes

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

// MultiwhereService wraps www.sudoc.fr/services/multiwhere/
type MultiwhereService struct {
	service
	endpoint string
	max_ppns int
}

const MAX_MULTIWHERE_PPNS = 50

// Library represents a location.
type Library struct {
	XMLName   xml.Name `xml:"library"`
	RCR       string   `xml:"rcr"`
	Shortname string   `xml:"shortname"`
	Latitude  float64  `xml:"latitude"`
	Longitude float64  `xml:"longitude"`
}

type result struct {
	XMLName   xml.Name  `xml:"result"`
	Libraries []Library `xml:"library"`
}

type query struct {
	XMLName xml.Name `xml:"query"`
	PPN     string   `xml:"ppn"`
	Result  result   `xml:"result"`
}

type serviceResult struct {
	XMLName xml.Name `xml:"sudoc"`
	Queries []query  `xml:"query"`
}

func (l Library) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	sb.WriteString(l.RCR)
	sb.WriteString("] ")
	sb.WriteString(l.Shortname)
	sb.WriteString(" (")
	sb.WriteString(strconv.FormatFloat(l.Latitude, 'f', -1, 64))
	sb.WriteString(", ")
	sb.WriteString(strconv.FormatFloat(l.Longitude, 'f', -1, 64))
	sb.WriteString(")")
	return sb.String()
}

// newMultiwhereService returnes a configured MultiwhereService instance.
func newMultiwhereService(client *http.Client, endpoint string) *MultiwhereService {
	ms := new(MultiwhereService)
	ms.client = client
	ms.endpoint = endpoint
	ms.max_ppns = MAX_MULTIWHERE_PPNS
	return ms
}

// GetLocations returns the list of the locations of the given PPN.
func (ms *MultiwhereService) GetLocations(ppn string) []Library {
	ppnList := []string{ppn}

	res := ms.GetMultiLocations(ppnList, 1)
	return res[ppn]
}

// GetMultiLocations returns a map associating each valid PPN to its locations,
// represented by a list of libraries.
func (ms *MultiwhereService) GetMultiLocations(ppns []string, max_ppns int) map[string][]Library {
	ppnStrings := ms.concatPPNs(ppns, max_ppns)
	result := make(map[string][]Library)

	for _, p := range ppnStrings {
		// TODO: handle do() errors
		res, err := ms.client.Get(ms.buildURL(ms.endpoint, p))
		if err != nil {
			fmt.Println(err)
			return nil
		}
		body, _ := io.ReadAll(res.Body)
		res.Body.Close()

		var sr serviceResult
		xml.Unmarshal(body, &sr)

		for _, query := range sr.Queries {
			for _, library := range query.Result.Libraries {
				result[query.PPN] = append(result[query.PPN], library)
			}
		}
	}
	return result
}

// GetMultiLocationsWithErrors returns a map associating each valid PPN to its
// locations - represented by a list of libraries - and a list of the  invalid
// PPNs among the requested ones.
func (ms *MultiwhereService) GetMultiLocationsWithErrors(ppns []string, max_ppns int) (map[string][]Library, []string) {
	result := ms.GetMultiLocations(ppns, max_ppns)
	var invalid_ppns []string
	var found_ppns []string
	for k := range result {
		found_ppns = append(found_ppns, k)
	}
	for _, ppn := range ppns {
		if !slices.Contains(found_ppns, ppn) {
			invalid_ppns = append(invalid_ppns, ppn)
		}
	}
	return result, invalid_ppns
}

// concatPPNs returns a list of what will be parameters for the multiwhere
// request, ie a list of concatenated PPNs.
func (ms *MultiwhereService) concatPPNs(ppns []string, max_ppns int) []string {
	if max_ppns < 1 {
		max_ppns = 1
	} else if max_ppns > ms.max_ppns {
		max_ppns = ms.max_ppns
	}
	res := []string{}
	for len(ppns) > max_ppns {
		res = append(res, strings.Join(ppns[:max_ppns], ","))
		ppns = ppns[max_ppns:]
	}
	if len(ppns) > 0 {
		res = append(res, strings.Join(ppns, ","))
	}
	return res
}

func (ms *MultiwhereService) buildURL(base, path string) string {
	if !strings.HasSuffix(base, "/") && !strings.HasPrefix(path, "/") {
		return base + "/" + path
	} else if strings.HasSuffix(base, "/") && strings.HasPrefix(path, "/") {
		return base + path[1:]
	} else {
		return base + path
	}
}
