package main

import (
	"fmt"

	"github.com/frp-2718/sudoc-client"
)

func main() {
	ppns := []string{"27681536X", "144089661", "154923206"}
	sc := sudoc.New(nil)
	res := sc.Locations(ppns)
	for k, v := range res {
		fmt.Println(k)
		for _, l := range v {
			fmt.Println(l)
		}
		fmt.Println()
	}
}
