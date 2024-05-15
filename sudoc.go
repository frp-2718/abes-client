package sudoc

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const BASE_URL = "https://www.sudoc.fr/services/multiwhere/"

type Sudoc struct {
	client *http.Client
}

type Library struct {
	XMLName   xml.Name `xml:"library"`
	RCR       string   `xml:"rcr"`
	Shortname string   `xml:"shortname"`
	Latitude  string   `xml:"latitude"`
	Longitude string   `xml:"longitude"`
}

func (l Library) String() string {
	return fmt.Sprintf("%s : %s", l.RCR, l.Shortname)
}

type Result struct {
	XMLName   xml.Name  `xml:"result"`
	Libraries []Library `xml:"library"`
}

type Query struct {
	XMLName xml.Name `xml:"query"`
	PPN     string   `xml:"ppn"`
	Result  Result   `xml:"result"`
}

type ServiceResult struct {
	XMLName xml.Name `xml:"sudoc"`
	Queries []Query  `xml:"query"`
}

func NewSudoc(client *http.Client) *Sudoc {
	sudoc := new(Sudoc)
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	sudoc.client = client
	return sudoc
}

func (s *Sudoc) do(url string) (*http.Response, error) {
	return http.Get(url)
}

func (s *Sudoc) Locations(ppns []string) map[string][]Library {
	ppnString := strings.Join(ppns, ",")
	res, _ := s.do(BASE_URL + ppnString)
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()

	var sr ServiceResult
	xml.Unmarshal(body, &sr)

	result := make(map[string][]Library)

	for _, query := range sr.Queries {
		for _, library := range query.Result.Libraries {
			result[query.PPN] = append(result[query.PPN], library)
		}
	}
	return result
}
