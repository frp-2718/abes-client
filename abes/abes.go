package abes

import "net/http"

const (
	sudocBaseURL       = "https://www.sudoc.fr/"
	servicesEndpoint   = "services/"
	multiwhereEndpoint = sudocBaseURL + servicesEndpoint + "multiwhere/"
	isbn2ppnEndpooint = sudocBaseURL + servicesEndpoint + "isbn2ppn/"
)

type service struct {
	client *http.Client
}

// Abes contains all exposed APIs
type Abes struct {
	Isbn2ppn *Isbn2ppnService
	Multiwhere *MultiwhereService
	UnimarcXML *UnimarcXMLService
}

// NewAbesClient returns a new initialized Abes client.
func NewAbesClient(client *http.Client) *Abes {
	if client == nil {
		client = http.DefaultClient
	}
	abes := new(Abes)
	abes.Isbn2ppn = newIsbn2ppnService(client, isbn2ppnEndpooint)
	abes.Multiwhere = newMultiwhereService(client, multiwhereEndpoint)
	abes.UnimarcXML = newUnimarcXMLService(client, sudocBaseURL)
	return abes
}
