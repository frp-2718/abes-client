package marc

import (
	"testing"
)

var errorXML = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<error>Les données bibliographiques sont indéfinies
  <ppn>PPN_IND1</ppn>
</error>`)

var correctXML = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<record>
  <leader>     cam0 22        450 </leader>
  <controlfield tag="003">http://www.sudoc.fr/155075381</controlfield>
  <controlfield tag="008">Aax3</controlfield>
  <controlfield tag="008">Oay3</controlfield>
  <datafield tag="010" ind1=" " ind2=" ">
    <subfield code="a">978-2-253-02983-0</subfield>
  </datafield>
  <datafield tag="200" ind1="1" ind2=" ">
    <subfield code="a">Orlando</subfield>
    <subfield code="f">Virginia Woolf</subfield>
    <subfield code="f">trad de Catherine Pappo</subfield>
  </datafield>
  <datafield tag="410" ind1=" " ind2="|">
    <subfield code="0">00102714X</subfield>
  </datafield>
  <datafield tag="702" ind1=" " ind2="1">
    <subfield code="3">02704324X</subfield>
    <subfield code="a">Pappo-Musard</subfield>
    <subfield code="b">Catherine</subfield>
    <subfield code="4">730</subfield>
  </datafield>
  <datafield tag="702" ind1=" " ind2="1">
    <subfield code="3">028265866</subfield>
    <subfield code="a">Nordon</subfield>
    <subfield code="b">Pierre</subfield>
  </datafield>
  <datafield tag="830" ind1=" " ind2=" ">
    <subfield code="a">pas un doublon</subfield>
  </datafield>
  <datafield tag="930" ind1="2" ind2=" ">
    <subfield code="5">ETAB1</subfield>
    <subfield code="c">Libre-accès</subfield>
    <subfield code="a">8 WOO</subfield>
  </datafield>
  <datafield tag="940" ind1=" " ind2=" ">
    <subfield code="5">ETAB2</subfield>
    <subfield code="a">20130626</subfield>
  </datafield>
  <datafield tag="930" ind1=" " ind2=" ">
    <subfield code="5">ETAB3</subfield>
    <subfield code="j">g</subfield>
  </datafield>
  <datafield tag="940" ind1=" " ind2=" ">
    <subfield code="5">ETAB4</subfield>
    <subfield code="a">20110923</subfield>
  </datafield>
</record>`)

func TestNewRecord(t *testing.T) {
	var tests = [][]byte{nil, errorXML}
	for _, test := range tests {
		_, err := NewRecord(test)
		if err == nil {
			t.Errorf(`NewRecord(%s) does not return error`, test)
		}
	}
	goodRecord, err := NewRecord(correctXML)
	if err != nil {
		t.Errorf(`NewRecord(correctXML) returns error : %s\n.`, err)
	}
	if len(goodRecord.Datafields) != 10 {
		t.Error(`len(goodRecord.Datafields) != 10`)
	}
}

func TestLeader(t *testing.T) {
	goodRecord, err := NewRecord([]byte(correctXML))
	if err != nil {
		t.Error(`Unable to create a new record from correctXML.`)
	}
	if result := goodRecord.Leader; result != "     cam0 22        450 " {
		t.Errorf(`goodRecord.Leader : got %s want "     cam0 22        450 "`, result)
	}
}

func TestIndicators(t *testing.T) {
	var tests = []struct {
		input string
		want  [][2]string
	}{
		{"008", nil},
		{"003", nil},
		{"000", nil},
		{"200", [][2]string{{"1", " "}}},
		{"410", [][2]string{{" ", "|"}}},
		{"702", [][2]string{{" ", "1"}, {" ", "1"}}},
		{"830", [][2]string{{" ", " "}}},
		{"invalid", nil},
	}
	goodRecord, err := NewRecord([]byte(correctXML))
	if err != nil {
		t.Error(`Unable to create a new record from correctXML.`)
	}
	for _, test := range tests {
		if got := goodRecord.Indicators(test.input); !equalInd(got, test.want) {
			t.Errorf("goodRecord.Indicators(%q) = %v", test.input, got)
		}
	}
}

func TestGetField(t *testing.T) {
	cf1 := Controlfield{Tag: "008", Value: "Aax3"}
	cf2 := Controlfield{Tag: "008", Value: "Oay3"}
	f1 := []Field{&cf1, &cf2}

	sf1 := Subfield{Code: "a", Value: "Orlando"}
	sf2 := Subfield{Code: "f", Value: "Virginia Woolf"}
	sf3 := Subfield{Code: "f", Value: "trad de Catherine Pappo"}
	df1 := Datafield{Tag: "200", Ind1: "1", Ind2: " ",
		Subfields: []Subfield{sf1, sf2, sf3}}
	f2 := []Field{&df1}

	sf4 := Subfield{Code: "5", Value: "ETAB1"}
	sf5 := Subfield{Code: "c", Value: "Libre-accès"}
	sf6 := Subfield{Code: "a", Value: "8 WOO"}
	df2 := Datafield{Tag: "930", Ind1: "2", Ind2: " ",
		Subfields: []Subfield{sf4, sf5, sf6}}

	sf7 := Subfield{Code: "5", Value: "ETAB3"}
	sf8 := Subfield{Code: "j", Value: "g"}
	df3 := Datafield{Tag: "930", Ind1: " ", Ind2: " ",
		Subfields: []Subfield{sf7, sf8}}
	f3 := []Field{&df2, &df3}

	goodRecord, err := NewRecord([]byte(correctXML))
	if err != nil {
		t.Error(`Unable to create a new record from correctXML.`)
	}

	var tests = []struct {
		input string
		want  []Field
	}{
		{"008", f1},
		{"200", f2},
		{"930", f3},
	}
	for _, test := range tests {
		if !equalFields(goodRecord.GetField(test.input), test.want) {
			t.Errorf(`GetFields(%q) failed`, test.input)
		}
	}
	if goodRecord.GetField("invalid") != nil {
		t.Error(`GetField does not return nil on invalid tag`)
	}
	df := goodRecord.GetField("702")[0].(*Datafield)
	if len(df.Subfields) != 4 {
		t.Error(`"702" subfields != 4`)
	}
}

func TestValue(t *testing.T) {
	goodRecord, err := NewRecord([]byte(correctXML))
	if err != nil {
		t.Error(`Unable to create a new record from correctXML.`)
	}

	var tests = []struct {
		inputTag  string
		inputCode string
		want      []string
	}{
		{"003", "", []string{"http://www.sudoc.fr/155075381"}},
		{"410", "0", []string{"00102714X"}},
		{"200", "f", []string{"Virginia Woolf", "trad de Catherine Pappo"}},
		{"940", "invalid", nil},
	}
	for _, test := range tests {
		f := goodRecord.GetField(test.inputTag)[0]
		if !equalStrings(f.GetValue(test.inputCode), test.want) {
			t.Errorf(`%s $%s != %q`, test.inputTag, test.inputCode, test.want)
		}
	}
}

// equalStrings compares two slices of strings.
func equalStrings(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, s := range s1 {
		if s != s2[i] {
			return false
		}
	}
	return true
}

// equalFields compares slices of Field.
func equalFields(s1, s2 []Field) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, f := range s1 {
		if !equalField(f, s2[i]) {
			return false
		}
	}
	return true
}

// equalField deeply compares Field structs.
func equalField(f1, f2 Field) bool {
	if cf1, ok := f1.(*Controlfield); ok {
		if cf2, ok := f2.(*Controlfield); ok {
			if cf1.Tag == cf2.Tag && cf1.Value == cf2.Value {
				return true
			}
		}
		return false
	}
	// f1 is a Datafield
	cf1 := f1.(*Datafield)
	if cf2, ok := f2.(*Datafield); ok {
		if cf1.Tag == cf2.Tag && cf1.Ind1 == cf2.Ind1 && cf1.Ind2 == cf2.Ind2 {
			for i, sub := range cf1.Subfields {
				if !equalSub(sub, cf2.Subfields[i]) {
					return false
				}
			}
			return true
		}
	}
	return false
}

// equalSub compares Subfield structs.
func equalSub(s1, s2 Subfield) bool {
	return s1.Code == s2.Code && s1.Value == s2.Value
}

// equalInd is a special purpose comparison function for indicators.
func equalInd(i1, i2 [][2]string) bool {
	if len(i1) != len(i2) {
		return false
	}
	for i, tab := range i1 {
		if tab[0] != i2[i][0] || tab[1] != i2[i][1] {
			return false
		}
	}
	return true
}
