package sudocclient

import (
	"net/http"
	"regexp"
	"time"
)

const (
	defaultMaxAttempts = 10
	sudocBaseURL       = "https://www.sudoc.fr/"
	idrefBaseURL       = "https://www.idref.fr/"
)

type service int

const (
	biblio service = iota
)

type responseFormat int

const (
	json responseFormat = iota
	xml
)

// Sudoc is the only access point to all ABES APIs.
type Sudoc struct {
	client      *http.Client
	Bibs        BibService
	maxAttempts int
}

type request struct {
	*http.Request
	attempts int
}

// New returns an initialized Sudoc struct.
func New() *Sudoc {
	sudoc := new(Sudoc)
	sudoc.client = &http.Client{
		Timeout: time.Second * 10,
	}
	sudoc.Bibs = BibService{client: sudoc}
	sudoc.maxAttempts = defaultMaxAttempts
	return sudoc
}

// SetClient allows user to provide a custom HTTP client.
func (s *Sudoc) SetHTTPClient(client *http.Client) {
	if client != nil {
		s.client = client
	}
}

// SetMaxAttempts adjusts the maximum number of attempts for an HTTP request
// before aborting. If 0, maxAttempts is set to default.
func (s *Sudoc) SetMaxAttempts(n int) {
	if n != 0 {
		s.maxAttempts = n
	} else {
		s.maxAttempts = defaultMaxAttempts
	}
}

// IsValidPPN checks if the given PPN is well formed.
func IsValidPPN(ppn string) bool {
	matched, _ := regexp.Match(`^[0-9]{8}?([0-9]{1}?|[xX])$`, []byte(ppn))
	return matched
}
