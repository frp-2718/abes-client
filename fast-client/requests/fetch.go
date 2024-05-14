// Requests package builds and handles HTTP requests.
package requests

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	multiwhere_url = "https://www.sudoc.fr/services/multiwhere/"
	iln2rcr_url    = "https://www.idref.fr/services/iln2rcr/"
	marcxml_url    = "https://www.sudoc.fr/"
)

type Fetcher interface {
	Fetch(url string) ([]byte, error)
	FetchMarc(ppn string) ([]byte, error)
}

type HttpFetcher struct {
	client *http.Client
}

func NewHttpFetch(client *http.Client) Fetcher {
	var fetcher HttpFetcher
	if client == nil {
		fetcher.client = &http.Client{Timeout: 5 * time.Second}
	} else {
		fetcher.client = client
	}
	return fetcher
}

// Fetch returns the xml record corresponding to the given URL, or nil if
// unsucessful.
func (f HttpFetcher) Fetch(url string) ([]byte, error) {
	resp, err := f.client.Get(url)
	if err != nil {
		// if request time out, just ignore
		// TODO: delay and request again
		// TODO: handle other url.errors
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

func (f HttpFetcher) FetchMarc(ppn string) ([]byte, error) {
	return f.Fetch(marcxml_url + ppn + ".xml")
}

func buildURLs(params []string, max_params int) []string {
	var urls []string
	for len(params) > max_params {
		newUrl := multiwhere_url + strings.Join(params[:max_params], ",")
		urls = append(urls, newUrl)
		params = params[max_params:]
	}
	urls = append(urls, multiwhere_url+strings.Join(params, ","))
	return urls
}
