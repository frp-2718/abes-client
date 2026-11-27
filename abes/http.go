package abes

import (
	"io"
	"net/http"
	"strings"
)

const maxErrorBodyBytes = 512

// do performs a GET request and returns the response body. Transport
// failures are wrapped in *NetworkError, 404 answers in *NotFoundError and
// any other non-2xx status in *APIError. For multiwhere, a 500 carrying an
// <error> element also means "not found" (notFoundOn500).
func (s *service) do(op, url string, notFoundOn500 bool) ([]byte, error) {
	res, err := s.client.Get(url)
	if err != nil {
		return nil, &NetworkError{Op: op, URL: url, Err: err}
	}
	defer res.Body.Close()

	// On error statuses the body may be truncated mid-stream (observed on
	// multiwhere 500 responses), so a read error is only fatal on 200.
	body, readErr := io.ReadAll(res.Body)
	if res.StatusCode == http.StatusOK {
		if readErr != nil {
			return nil, &SystemError{Op: op, URL: url, Err: readErr}
		}
		return body, nil
	}

	msg := errorMessage(body)
	switch {
	case res.StatusCode == http.StatusNotFound:
		return nil, &NotFoundError{Op: op, URL: url, Message: msg}
	case notFoundOn500 && res.StatusCode == http.StatusInternalServerError && msg != "":
		return nil, &NotFoundError{Op: op, URL: url, Message: msg}
	default:
		return nil, &APIError{Op: op, URL: url, StatusCode: res.StatusCode, Body: truncate(body, maxErrorBodyBytes)}
	}
}

// errorMessage extracts the text of an <error> element. ABES error bodies are
// not always well-formed XML (the <ppn>.xml 404 body embeds an unclosed <ppn>
// tag inside <error>), so parse defensively.
func errorMessage(body []byte) string {
	s := string(body)
	i := strings.Index(s, "<error>")
	if i < 0 {
		return ""
	}
	s = s[i+len("<error>"):]
	if j := strings.IndexByte(s, '<'); j >= 0 {
		s = s[:j]
	}
	return strings.TrimSpace(s)
}

func truncate(b []byte, max int) string {
	if len(b) <= max {
		return string(b)
	}
	return string(b[:max]) + "..."
}
