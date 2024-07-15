package abes

import "net/http"

const (
	sudocBaseURL       = "https://www.sudoc.fr/"
	servicesEndpoint   = "services/"
	multiwhereEndpoint = sudocBaseURL + servicesEndpoint + "multiwhere/"
)

type service struct {
	client *http.Client
}

// Abes contains all exposed APIs
type Abes struct {
	multiwhere *MultiwhereService
}

// NewAbesClient returns a new initialized Abes client.
func NewAbesClient(client *http.Client) *Abes {
	if client == nil {
		client = http.DefaultClient
	}
	abes := new(Abes)
	abes.multiwhere = new(MultiwhereService)
	abes.multiwhere.client = client
	abes.multiwhere.endpoint = multiwhereEndpoint
	abes.multiwhere.max_ppns = 50
	return abes
}
