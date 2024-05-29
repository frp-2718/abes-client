package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/frp-2718/sudoc-client"
)

func main() {
	ppns := []string{"273075969", "273091824", "273094378", "219527149",
		"224236024", "273095757", "273117165", "273562053", "273799304",
		"273961152", "274048655", "27515064X", "275435601", "275538842",
		"275862984"}

	client := sudoc.New(nil)

	for _, ppn := range ppns {
		locations := client.Locations([]string{ppn})
		rcrs := []string{}
		for _, library := range locations[ppn] {
			rcrs = append(rcrs, library.RCR)
		}
		if len(rcrs) > 1 && slices.Contains(rcrs, "914712302") {
			fmt.Printf("%s;%s\n", ppn, strings.Join(rcrs, ","))
		}
	}
}
