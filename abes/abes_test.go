package abes

import (
	"fmt"
	"testing"
)

func TestAbesClient(t *testing.T) {
	abes := NewAbesClient(nil)
	ppns := []string{"144089661", "154923206"}
	res := abes.multiwhere.GetMultiLocations(ppns, 10)
	for k, v := range res {
		fmt.Println(k)
		for _, l := range v {
			fmt.Println(l)
		}
		fmt.Println("===================================")
	}
}
