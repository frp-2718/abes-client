// Package marcxml provides some functions to extract informations from
// a MARCXML record.
package marc

import (
	"encoding/xml"
	"errors"
)

type Record struct {
	XMLName       xml.Name       `xml:"record"`
	Leader        string         `xml:"leader"`
	Controlfields []Controlfield `xml:"controlfield"`
	Datafields    []Datafield    `xml:"datafield"`
}

type Controlfield struct {
	XMLName xml.Name `xml:"controlfield"`
	Tag     string   `xml:"tag,attr"`
	Value   string   `xml:",chardata"`
}

type Datafield struct {
	XMLName   xml.Name   `xml:"datafield"`
	Tag       string     `xml:"tag,attr"`
	Ind1      string     `xml:"ind1,attr"`
	Ind2      string     `xml:"ind2,attr"`
	Subfields []Subfield `xml:"subfield"`
}

type Subfield struct {
	XMLName xml.Name `xml:"subfield"`
	Code    string   `xml:"code,attr"`
	Value   string   `xml:",chardata"`
}

// Discriminated union type. A Field includes Controlfield and Datafield.
type Field interface {
	GetValue(code string) []string
}

// NewRecord converts xml data into a Record, assuming that the provided
// xml data is a valid mono-record marcXML.
func NewRecord(xmlData []byte) (*Record, error) {
	if xmlData == nil {
		return nil, errors.New("NewRecord: can't process nil data")
	}
	var r Record
	err := xml.Unmarshal(xmlData, &r)
	if err != nil {
		//log.Printf("NewRecord: %s", err)
		return nil, err
	}
	return &r, nil
}

// Indicators returns a list of pairs of indicators for a given tag. One entry
// of the list corresponds to a repeated field. The list is nil if the field is
// a controlfield.
func (r *Record) Indicators(tag string) [][2]string {
	for _, field := range r.Controlfields {
		if field.Tag == tag {
			return nil
		}
	}
	var result [][2]string
	for _, field := range r.Datafields {
		if field.Tag == tag {
			result = append(result, [2]string{field.Ind1, field.Ind2})
		}
	}
	return result
}

// GetField returns all fields with corresponding tag.
func (r *Record) GetField(tag string) []Field {
	var res []Field
	for _, field := range r.Controlfields {
		if field.Tag == tag {
			res = append(res, &Controlfield{Tag: tag, Value: field.Value})
		}
	}
	for _, field := range r.Datafields {
		if field.Tag == tag {
			res = append(res, &Datafield{Tag: tag, Ind1: field.Ind1, Ind2: field.Ind2, Subfields: field.Subfields})
		}
	}
	return res
}

// GetValue returns a slice of the subfields of given code.
func (cf *Controlfield) GetValue(code string) []string {
	return []string{cf.Value}
}

func (df *Datafield) GetValue(code string) []string {
	var result []string
	for _, sub := range df.Subfields {
		if sub.Code == code {
			result = append(result, sub.Value)
		}
	}
	return result
}
