package abes

import (
	"encoding/xml"
	"strconv"
	"strings"
)

// MultiwhereService wraps www.sudoc.fr/services/multiwhere/
type MultiwhereService service

type Library struct {
	XMLName   xml.Name `xml:"library"`
	RCR       string   `xml:"rcr"`
	Shortname string   `xml:"shortname"`
	Latitude  float64  `xml:"latitude"`
	Longitude float64  `xml:"longitude"`
}

type result struct {
	XMLName   xml.Name  `xml:"result"`
	Libraries []Library `xml:"library"`
}

type query struct {
	XMLName xml.Name `xml:"query"`
	PPN     string   `xml:"ppn"`
	Result  result   `xml:"result"`
}

type serviceResult struct {
	XMLName xml.Name `xml:"sudoc"`
	Queries []query  `xml:"query"`
}

func (l Library) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	sb.WriteString(l.RCR)
	sb.WriteString("] ")
	sb.WriteString(l.Shortname)
	sb.WriteString(" (")
	sb.WriteString(strconv.FormatFloat(l.Latitude, 'f', -1, 64))
	sb.WriteString(", ")
	sb.WriteString(strconv.FormatFloat(l.Longitude, 'f', -1, 64))
	sb.WriteString(")")
	return sb.String()
}
