package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/frp-2718/abes-client/abes"
)

func main() {
	isbn := "978-0134190440"
	isbns := []string{"978-2765408970", "154923206", "978273827883T"}

	// Custom HTTP client with timeout settings
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSHandshakeTimeout = 20 * time.Second
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
	ac := abes.NewAbesClient(httpClient)

	res_mono, err := ac.Isbn2ppn.GetPPNs(isbn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur lors de la recherche de l'ISBN %s : %v\n", isbn, err)
		return
	}
	fmt.Println("Interrogation d'un seul ISBN : " + isbn)
	fmt.Println()
	fmt.Println(res_mono)

	fmt.Println("\n=====================================")
	fmt.Println()

	fmt.Println("Interrogation de plusieurs ISBN : " + fmt.Sprint(isbns))
	fmt.Println()
	res_multi, err := ac.Isbn2ppn.GetMultiPPNs(isbns, 20)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur lors de la recherche des ISBN %v : %v\n", isbns, err)
		return
	}
	for _, v := range res_multi {
		fmt.Println(v)
	}
}
