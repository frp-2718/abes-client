package sudoc

import "net/http"

// Sudoc is the only access point to all ABES APIs.
type Sudoc struct {
	client *http.Client
}

// New returns an initialized Sudoc struct.
func New(client *http.Client) *Sudoc {
	if client == nil {
		client = http.DefaultClient
	}
	sudoc := new(Sudoc)
	sudoc.client = client
	return sudoc
}
