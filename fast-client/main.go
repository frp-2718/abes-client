package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"sudoc/client"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("no ppn file provided")
	}

	client := client.NewSudocClient()

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		ppn := scanner.Text()
		locations := client.GetLocations(ppn)
		if slices.Contains(locations, "914712302") && len(locations) > 1 {
			fmt.Printf("%s;%v", ppn, locations)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
