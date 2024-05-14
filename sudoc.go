package sudoc

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const BASE_URL = "https://www.sudoc.fr/services/multiwhere/"

type Sudoc struct {
}

type Library struct {
	XMLName   xml.Name `xml:"library"`
	RCR       string   `xml:"rcr"`
	Shortname string   `xml:"shortname"`
	Latitude  string   `xml:"latitude"`
	Longitude string   `xml:"longitude"`
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

func New() *Sudoc {
	return new(Sudoc)
}

func (s *Sudoc) Locations(ppns []string) {
	ppnString := strings.Join(ppns, ",")
	res, _ := http.Get(BASE_URL + ppnString)
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()

	var sr ServiceResult
	xml.Unmarshal(body, &sr)

	for _, query := range sr.Queries {
		fmt.Println(query.PPN)
		for _, library := range query.Result.Libraries {
			fmt.Printf("%s : %s\n", library.RCR, library.Shortname)
		}
		fmt.Println("#######################")
	}
}
