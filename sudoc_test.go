package sudoc

import (
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	http_client := &http.Client{}
	sudoc := New(nil)
	if sudoc == nil {
		t.Fatal("sudoc was not initialized")
	}
	if sudoc.client == nil {
		t.Fatal("sudoc HTTP client was not set")
	}
	sudoc = New(http_client)
	if sudoc.client != http_client {
		t.Fatalf("sudoc HTTP client was not set: have %v, want %v", sudoc.client, http_client)
	}
}
