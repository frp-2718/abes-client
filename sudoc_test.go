package sudoc

import (
	"net/http"
	"reflect"
	"testing"
)

// Ad-hoc experimentation.
func assertNotNil(object interface{}, t *testing.T, label string) {
	if object == nil || (reflect.ValueOf(object).Kind() == reflect.Ptr && reflect.ValueOf(object).IsNil()) {
		t.Errorf("%s is nil", label)
	}
}

func TestSudoc(t *testing.T) {
	var tests = []struct {
		sudoc *Sudoc
		label string
	}{
		{New(nil), "New(nil)"},
		{New(http.DefaultClient), "New(http.DefaultClient)"},
	}
	for _, test := range tests {
		assertNotNil(test.sudoc, t, test.label)
		if test.sudoc != nil {
			assertNotNil(test.sudoc.client, t, test.label+".client")
			assertNotNil(test.sudoc.Bibs, t, test.label+".Bibs")
		}
	}
}
