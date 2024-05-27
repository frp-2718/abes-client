package sudoc

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
)

func TestDecodeError(t *testing.T) {
	tests := []struct {
		name  string
		input *http.Response
		want  error
	}{
		{
			name: "one wrong ppn",
			input: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body: io.ReadCloser(io.NopCloser(bytes.NewBufferString(`
                <?xml version="1.0" encoding="UTF-8" ?>
                <sudoc service="multiwhere">
                <error>Found a null xml in result : values={ppn=a}, query=select autorites.MULTIWHERE(#ppn#) from dual</error>
                </sudoc>`)))},
			want: &NotFoundError{"a"},
		},
		{
			name: "several wrong ppn",
			input: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body: io.ReadCloser(io.NopCloser(bytes.NewBufferString(`
                <?xml version="1.0" encoding="UTF-8" ?>
                <sudoc service="multiwhere">
                <error>Found a null xml in result : values={ppn=a,b,c}, query=select autorites.MULTIWHERE(#ppn#) from dual</error>
                </sudoc>`)))},
			want: &NotFoundError{"a,b,c"},
		},
		{
			name: "empty request",
			input: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body: io.ReadCloser(io.NopCloser(bytes.NewBufferString(`
                    <?xml version="1.0" encoding="UTF-8"?>
                    <sudoc service="multiwhere">
                    <error>Invalid char in query string, values={}, query=select autorites.MULTIWHERE(#ppn#) from dual</error>
                    </sudoc>`)))},
			want: &InvalidRequestError{},
		},
		{
			name: "unknown error",
			input: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body: io.ReadCloser(io.NopCloser(bytes.NewBufferString(`
                    <?xml version="1.0" encoding="UTF-8"?>
                    <sudoc service="multiwhere">
                    <error>unknown</error>
                    </sudoc>`)))},
			want: &UnknownError{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := decodeError(test.input)

			if test.want != nil {
				switch want := test.want.(type) {
				case *NotFoundError:
					var notfound *NotFoundError
					if !errors.As(got, &notfound) {
						t.Errorf("want NotFoundError, got %v", got)
					} else if notfound.PPN != want.PPN {
						t.Errorf("want %v, got %v", want, got)
					}
				case *InvalidRequestError:
					var invalid *InvalidRequestError
					if !errors.As(got, &invalid) {
						t.Errorf("want InvalidRequestError, got %v", got)
					}
				case *UnknownError:
					var unknown *UnknownError
					if !errors.As(got, &unknown) {
						t.Errorf("want UnknownError, got %v", got)
					}
				default:
					t.Errorf("unexpected error type: %T", got)
				}
			}
		})
	}
}

// func decodeError(r *http.Response) error {
// 	content, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		log.Println("decodeError: unable to read HTTP response body")
// 		return err
// 	}
// 	r.Body.Close()

// 	var e multiwhere_error
// 	err = xml.Unmarshal(content, &e)
// 	if err != nil {
// 		log.Println("decodeError: unable to unmarshal XML content")
// 		return err
// 	}

// 	if strings.Contains(e.ErrorText, "null xml") {
// 		re := regexp.MustCompile(`ppn=([^}]+)`)
// 		match := re.FindStringSubmatch(e.ErrorText)
// 		ppn := ""
// 		if match != nil {
// 			ppn = match[1]
// 		}
// 		return &NotFoundError{ppn}
// 	} else if strings.Contains(e.ErrorText, "invalid character") {
// 		return &InvalidRequestError{}
// 	} else {
// 		return errors.New("unknown error")
// 	}
// }

// // multiwhere returns 500 for any error
// func checkForErrors(r *http.Response) error {
// 	if r.StatusCode == http.StatusInternalServerError {
// 		return decodeError(r)
// 	} else if r.StatusCode >= 200 && r.StatusCode <= 299 {
// 		return nil
// 	} else {
// 		return errors.New(fmt.Sprintf("Unknown error: HTTP %d", r.StatusCode))
// 	}
// }
