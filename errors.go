package sudoc

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type multiwhere_error struct {
	ErrorText string `xml:"error"`
}

type NotFoundError struct {
	PPN string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.PPN)
}

type InvalidRequestError struct{}

func (e *InvalidRequestError) Error() string {
	return fmt.Sprint("Empty request")
}

func decodeError(r *http.Response) error {
	content, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("decodeError: unable to read HTTP response body")
		return err
	}
	r.Body.Close()

	var e multiwhere_error
	err = xml.Unmarshal(content, &e)
	if err != nil {
		log.Println("decodeError: unable to unmarshal XML content")
		return err
	}

	if strings.Contains(e.ErrorText, "null xml") {
		re := regexp.MustCompile(`ppn=(?P<ppnValue>[^}]+)`)
		ppn := re.FindString(e.ErrorText)
		return &NotFoundError{ppn}
	} else if strings.Contains(e.ErrorText, "invalid character") {
		return &InvalidRequestError{}
	} else {
		return errors.New("unknown error")
	}
}

// multiwhere returns 500 for any error
func checkForErrors(r *http.Response) error {
	if r.StatusCode == http.StatusInternalServerError {
		return decodeError(r)
	} else if r.StatusCode >= 200 && r.StatusCode <= 299 {
		return nil
	} else {
		return errors.New(fmt.Sprintf("Unknown error: HTTP %d", r.StatusCode))
	}
}
