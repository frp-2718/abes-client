package abes

import "net/http"

type service struct{}

// Abes contains all exposed APIs
type Abes struct {
	client     *http.Client
	Multiwhere *MultiwhereService
}

// NewAbesClient returns a new initialized Abes client.
func NewAbesClient(client *http.Client) *Abes {
	if client == nil {
		client = http.DefaultClient
	}
	abes := new(Abes)
	return abes
}
