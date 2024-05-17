package sudoc

import (
	"net/http"
	"strings"
)

func (s *Sudoc) do(url string) (*http.Response, error) {
	return http.Get(url)
}

func (s *Sudoc) buildURL(base, path string) string {
	return base + path
}

func (s *Sudoc) concatPPNs(ppns []string, max_ppns int) []string {
	if max_ppns < 1 {
		max_ppns = 1
	} else if max_ppns > MAX_MULTIWHERE_PPNS {
		max_ppns = MAX_MULTIWHERE_PPNS
	}
	res := []string{}
	for len(ppns) > max_ppns {
		res = append(res, strings.Join(ppns[:max_ppns], ","))
		ppns = ppns[max_ppns:]
	}
	if len(ppns) > 0 {
		res = append(res, strings.Join(ppns, ","))
	}
	return res
}
