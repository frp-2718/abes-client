package abes

import "net/http"

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
// A 404 answer (or a 200 body carrying an <error> element) is reported as a
// *NotFoundError carrying the server's message.
func (s *UnimarcXMLService) GetRecord(ppn string) (*MarcRecord, error) {
	url := s.buildURL(ppn)
	body, err := s.do("UnimarcXML.GetRecord", url, false)
	if err != nil {
		return nil, err
	}

	if msg := errorMessage(body); msg != "" {
		return nil, &NotFoundError{Op: "UnimarcXML.GetRecord", URL: url, Message: msg}
	}

	record, err := NewRecord(body)
	if err != nil {
		return nil, &SystemError{Op: "UnimarcXML.GetRecord", URL: url, Err: err}
	}
	return record, nil
}

func (s *UnimarcXMLService) buildURL(ppn string) string {
	return s.endpoint + ppn + ".xml"
}
