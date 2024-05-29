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

type InvalidRequestError struct{}

func (e *InvalidRequestError) Error() string {
	return fmt.Sprint("empty request")
}

type NotFoundError struct {
	PPN string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.PPN)
}

type UnexpectedError struct {
	PPN string
}

func (e *UnexpectedError) Error() string {
	return fmt.Sprint("unexpected error", e.PPN)
}

func decodeError(r *http.Response) error {
	if r.Body == nil {
		return &UnexpectedError{}
	}
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
		re := regexp.MustCompile(`ppn=([^}]+)`)
		match := re.FindStringSubmatch(e.ErrorText)
		ppn := ""
		if match != nil {
			ppn = match[1]
		}
		return &NotFoundError{ppn}
	} else if strings.Contains(e.ErrorText, "Invalid char") {
		return &InvalidRequestError{}
	} else {
		return &UnexpectedError{}
	}
}

// multiwhere returns 500 for any error
func checkForErrors(r *http.Response) error {
	if r.StatusCode == http.StatusInternalServerError {
		return decodeError(r)
	} else if r.StatusCode >= 200 && r.StatusCode <= 299 {
		return nil
	} else {
		return errors.New(fmt.Sprintf("unknown error: HTTP %d", r.StatusCode))
	}
}
