package sudoc

import (
	"net/http"
	"time"
)

// Sudoc is the only access point to all ABES APIs.
type Sudoc struct {
	client *http.Client
	Bibs   BibService
}

// New returns an initialized Sudoc struct.
func New() *Sudoc {
	sudoc := new(Sudoc)
	sudoc.client = &http.Client{
		Timeout: time.Second * 10,
	}
	sudoc.Bibs = BibService{client: sudoc}
	return sudoc
}
