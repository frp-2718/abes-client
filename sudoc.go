package sudoc

import (
	"net/http"
	"regexp"
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

// SetClient allow user to provide a custom HTTP client.
func (s *Sudoc) SetHTTPClient(client *http.Client) {
	if client != nil {
		s.client = client
	}
}

// IsValidPPN checks if the given PPN is well formed.
func IsValidPPN(ppn string) bool {
	matched, _ := regexp.Match(`^[0-9]{8}?([0-9]{1}?|[xX])$`, []byte(ppn))
	return matched
}
