package client

import (
	"log"
	"strings"
	"sudoc/marc"
	"sudoc/requests"
)

type Client interface {
	GetLocations(url string) []string
}

type SudocClient struct {
	fetcher requests.Fetcher
}

func NewSudocClient() SudocClient {
	var client SudocClient
	client.fetcher = requests.NewHttpFetch(nil)
	return client
}

func (c SudocClient) GetLocations(ppn string) []string {
	data, err := c.fetcher.FetchMarc(ppn)
	if err != nil {
		return []string{}
	}
	record, err := marc.NewRecord(data)
	if err != nil {
		log.Fatal(err)
	}
	items := record.GetField("930")
	var all_loc []string
	for _, item := range items {
		localizations := item.GetValue("5")
		for _, l := range localizations {
			all_loc = append(all_loc, strings.Split(l, ":")[0])
		}
	}
	return all_loc
}
