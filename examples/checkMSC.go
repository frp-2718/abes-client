package main

import (
	"fmt"

	"github.com/frp-2718/abes-client/abes"
)

func main() {
	ppns := []string{"154923206", "111111111", "144089661", "103835091"}

	ac := abes.NewAbesClient(nil)

	for _, ppn := range ppns {
		marc, err := ac.UnimarcXML.GetRecord(ppn)
		fmt.Printf("PPN %s ", ppn)
		if err != nil {
			fmt.Println("non trouvé")
		} else {
			if hasMSC(marc) {
				fmt.Println("OK")
			} else {
				fmt.Println("686$2msc à créer ou corriger")
			}
		}
	}
}

// Quick and dirty test function to avoid overloading the example.
func hasMSC(record *abes.MarcRecord) bool {
	fields := record.GetField("686")
	if len(fields) == 0 {
		return false
	}
	for _, f := range fields {
		code := f.GetValue("2")
		if len(code) > 0 && code[0] == "msc" {
			return true
		}
	}
	return false
}
