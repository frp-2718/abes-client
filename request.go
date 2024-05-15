package sudoc

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
)

const (
	multiwhere_url = "https://www.sudoc.fr/services/multiwhere/"
)

// type Fetcher interface {
// 	Fetch(url string) ([]byte, error)
// }

// type HttpFetcher struct {
// 	client *http.Client
// }

// func NewHttpFetcher(client *http.Client) Fetcher {
// 	var fetcher HttpFetcher
// 	if client == nil {
// 		fetcher.client = &http.Client{Timeout: 5 * time.Second}
// 	} else {
// 		fetcher.client = client
// 	}
// 	return fetcher
// }

func (s *Sudoc) Fetch(url string) ([]byte, error) {
	resp, err := s.client.Get(url)
	if err != nil {
		return []byte{}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return []byte{}, errors.New(strconv.Itoa(resp.StatusCode))
	}
	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}
	return data, nil
}
