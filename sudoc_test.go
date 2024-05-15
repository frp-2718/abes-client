package sudoc

import (
	"net/http"
	"testing"
)

func TestNewSudoc(t *testing.T) {
	http_client := &http.Client{}
	sudoc := NewSudoc(nil)
	if sudoc == nil {
		t.Fatal("sudoc was not initialized")
	}
	if sudoc.client == nil {
		t.Fatal("sudoc HTTP client was not set")
	}
	sudoc = NewSudoc(http_client)
	if sudoc.client != http_client {
		t.Fatalf("sudoc HTTP client was not set: have %v, want %v", sudoc.client, http_client)
	}
}

// func TestSudoc(T *testing.T) {
// 	ppns := []string{"144089661", "154923206"}
// 	sc := NewSudoc(nil)
// 	res := sc.Locations(ppns)
// 	for k, v := range res {
// 		fmt.Println(k)
// 		for _, l := range v {
// 			fmt.Println(l)
// 		}
// 		fmt.Println()
// 	}
// }
