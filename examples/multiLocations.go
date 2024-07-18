package main

import (
	"fmt"

	"github.com/frp-2718/abes-client/abes"
)

func main() {
	ppns := []string{"144089661", "154923206"}

	ac := abes.NewAbesClient(nil)
	res, _ := ac.Multiwhere.GetMultiLocations(ppns, 20)

	for k, v := range res {
		fmt.Println(k)
		for _, l := range v {
			fmt.Println(l)
		}
		fmt.Println("===========================")
	}
}
