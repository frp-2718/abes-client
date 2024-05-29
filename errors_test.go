package sudoc

import (
	"bytes"
	"errors"
	"fmt"
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
			want: &UnexpectedError{},
		},
		{
			name: "empty body",
			input: &http.Response{
				StatusCode: http.StatusInternalServerError},
			want: &UnexpectedError{},
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
				case *UnexpectedError:
					var unknown *UnexpectedError
					if !errors.As(got, &unknown) {
						t.Errorf("want UnexpectedError, got %v", got)
					}
				default:
					t.Errorf("unexpected error type: %T", got)
				}
			}
		})
	}
}

func TestCheckForError(t *testing.T) {
	tests := []struct {
		name  string
		input *http.Response
		want  error
	}{
		{
			name: "500",
			input: &http.Response{
				StatusCode: http.StatusInternalServerError},
			want: &UnexpectedError{},
		},
		{
			name: "404",
			input: &http.Response{
				StatusCode: http.StatusNotFound},
			want: errors.New(fmt.Sprintf("unknown error: HTTP %d", http.StatusNotFound)),
		},
		{
			name: "200",
			input: &http.Response{
				StatusCode: http.StatusOK},
			want: nil,
		},
		{
			name: "206",
			input: &http.Response{
				StatusCode: http.StatusPartialContent},
			want: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := checkForErrors(test.input)
			if test.want != nil {
				if got == nil {
					t.Errorf("want nil, got %v", got)
				} else if test.want.Error() != got.Error() {
					t.Errorf("want %v, got %v", test.want.Error(), got.Error())
				}
			}
		})
	}
}
