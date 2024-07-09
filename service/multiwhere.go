package service

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/frp-2718/sudoc-client/models"
)

type MultiwhereService struct {
	httpClient *http.Client
	baseURL    string
}

func NewMultiwhereService(client *http.Client, baseURL string) *MultiwhereService {
	return &MultiwhereService{
		httpClient: client,
		baseURL:    baseURL,
	}
}

func (s *MultiwhereService) GetResponse(requestParams map[string]string) (*models.Response, error) {
	url := buildURL(s.baseURL, requestParams["ppns"])
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request error")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io error")
	}

	return &models.Response{
		Data: body,
	}, nil
}

func buildURL(base, ppns string) string {
	if !strings.HasSuffix(base, "/") && !strings.HasPrefix(ppns, "/") {
		return base + "/" + ppns
	} else if strings.HasSuffix(base, "/") && strings.HasPrefix(ppns, "/") {
		return base + ppns[1:]
	} else {
		return base + ppns
	}
}
