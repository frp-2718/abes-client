package abes

import (
	"io"
	"net/http"
)

// UnimarcXML wraps www.sudoc.fr/<PPN>.xml
type UnimarcXMLService struct {
	service
	endpoint string
}

// newUnimarcXMLService returnes a configured UnimarcXMLService instance.
func newUnimarcXMLService(client *http.Client, endpoint string) *UnimarcXMLService {
	s := new(UnimarcXMLService)
	s.client = client
	s.endpoint = endpoint
	return s
}

// GetRecord returns the parsed MARC Record corresponding to the provided PPN.
func (s *UnimarcXMLService) GetRecord(ppn string) (*MarcRecord, error) {
	res, err := s.client.Get(s.buildURL(ppn))
	if err != nil {
		return nil, &NetworkError{"HTTP protocol error"}
	}

	if res.StatusCode != http.StatusOK {
		return nil, &NotFoundError{}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &SystemError{"unable to read body"}
	}
	res.Body.Close()

	record, err := NewRecord(body)
	if err != nil {
		return nil, &SystemError{"unable to parse MARC data"}
	}
	return record, nil
}

func (s *UnimarcXMLService) buildURL(ppn string) string {
	return s.endpoint + ppn + ".xml"
}
