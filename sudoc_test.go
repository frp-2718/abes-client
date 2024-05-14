package sudoc

import (
	"testing"
)

func TestSudoc(T *testing.T) {
	ppns := []string{"144089661", "154923206"}
	sc := New()
	sc.Locations(ppns)
}
