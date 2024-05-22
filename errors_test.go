package sudoc

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestDecodeError(t *testing.T) {
	response := http.Response{
		StatusCode: http.StatusInternalServerError,
		Body: io.ReadCloser(io.NopCloser(bytes.NewBufferString(`
        <?xml version="1.0" encoding="UTF-8" ?>
        <sudoc service="multiwhere">
        <error>Found a null xml in result : values={ppn=a,b,c}, query=select autorites.MULTIWHERE(#ppn#) from dual</error>
        </sudoc>
        `))),
	}
	fmt.Println(decodeError(&response))
}
