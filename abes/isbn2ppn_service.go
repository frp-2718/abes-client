package abes

import (
	"encoding/xml"
	"net/http"
	"strings"
)

// Official documentation of this service strongly advise to sanitize the ISBNs
// before using the webservice (see https://documentation.abes.fr/sudoc/manuels/administration/aidewebservices/index.html#isbn2ppn :
// "En d'autres termes, pour éviter d'avoir des réponses inconsistantes,
// le plus sûr est de valider ses ISBN préalablement à l'interrogation du
// webservice .")

// ISBN can be given with or without dashes.

// Isbn2ppnService wraps www.sudoc.fr/services/isbn2ppn/
type Isbn2ppnService struct {
	service
	endpoint  string
	max_isbns int
}

const MAX_ISBN2PPN_ISBN = 50

// Webservice XML response
type SudocResponse struct {
	XMLName xml.Name     `xml:"sudoc"`
	Service string       `xml:"service,attr"`
	Queries []ISBNResult `xml:"query"`
}

type ISBNResult struct {
	ISBN      string   `xml:"isbn"`
	Holdings  []string `xml:"result>ppn"`
	NoHolding []string `xml:"resultNoHolding>ppn"`
}

func (i ISBNResult) String() string {
	var sb strings.Builder
	sb.WriteString(i.ISBN)
	sb.WriteString(" -> [")

	for i, ppn := range i.Holdings {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(ppn)
	}

	sb.WriteString("][")

	for i, ppn := range i.NoHolding {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(ppn)
	}

	sb.WriteString("]")
	return sb.String()
}

// newIsbn2ppnService returnes a configured Isbn2ppnService instance.
func newIsbn2ppnService(client *http.Client, endpoint string) *Isbn2ppnService {
	ms := new(Isbn2ppnService)
	ms.client = client
	ms.endpoint = endpoint
	ms.max_isbns = MAX_ISBN2PPN_ISBN
	return ms
}

// GetPPNs returns two lists of PPNs for the given ISBN : first one is for
// PPNs with holdings and second one for PPNs with no holdings.
func (is *Isbn2ppnService) GetPPNs(isbn string) (ISBNResult, error) {
	isbnList := []string{isbn}

	res, err := is.GetMultiPPNs(isbnList, 1)
	if result, ok := res[isbn]; ok {
		return result, err
	}
	// The service may return a normalized ISBN (for example, ISBN-10 when
	// queried with its ISBN-13 equivalent), so use the sole result for this
	// single-ISBN request when the original key is absent.
	for _, result := range res {
		return result, err
	}
	return ISBNResult{}, err
}

// GetMultiPPNs returns a map associating each valid ISBN to its PPNs,
// represented by two lists, possibly empty : with holdings/without holdings. If
// an ISBN is not found, it is ignored.
func (is *Isbn2ppnService) GetMultiPPNs(isbns []string, max_isbns int) (map[string]ISBNResult, error) {
	isbnStrings := concatValues(isbns, max_isbns, is.max_isbns)
	result := make(map[string]ISBNResult)

	for _, p := range isbnStrings {
		url := joinURLPath(is.endpoint, p)
		body, err := is.do("Isbn2ppn.GetMultiPPNs", url, false)
		if err != nil {
			return nil, err
		}

		var sr SudocResponse
		if err := xml.Unmarshal(body, &sr); err != nil {
			return nil, &SystemError{Op: "Isbn2ppn.GetMultiPPNs", URL: url, Err: err}
		}

		for _, query := range sr.Queries {
			result[query.ISBN] = query
		}
	}
	return result, nil
}
