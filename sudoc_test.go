package sudoc

import (
	"net/http"
	"testing"
)

func TestSudoc(t *testing.T) {
	sudoc := New(nil)
	if sudoc == nil {
		t.Error("New(nil) returned nil")
	}
	if sudoc.client == nil {
		t.Error("New(nil) client initialization failed")
	}
	sudoc = New(http.DefaultClient)
	if sudoc == nil {
		t.Error("New(client) returned nil")
	}
	if sudoc.client == nil {
		t.Error("New(client) client initialization failed")
	}
}
