package main

import (
	"fmt"
	"log"

	"github.com/frp-2718/abes-client/abes"
)

func main() {
	ppn := "111111111"

	ac := abes.NewAbesClient(nil)
	res, err := ac.Multiwhere.GetLocations(ppn)
	if err != nil {
		log.Fatal(err)
	}

	if len(res) == 0 {
		fmt.Printf("PPN %s ne correspond pas à une notice bibliographqiue vivante.\n", ppn)
	} else {
		for _, l := range res {
			fmt.Println(l)
		}
	}
}
