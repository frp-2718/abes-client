package sudoc

import (
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
	sudoc := New()
	assertNotNil(sudoc, t, "New()")
	if sudoc != nil {
		assertNotNil(sudoc.client, t, ".client")
		assertNotNil(sudoc.Bibs, t, ".Bibs")
	}
}
