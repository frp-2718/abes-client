package abes

import (
	"reflect"
	"testing"
)

func TestConcatValues(t *testing.T) {
	tests := []struct {
		name       string
		values     []string
		maxValues  int
		serviceMax int
		want       []string
	}{
		{"empty input", nil, 2, 3, []string{}},
		{"single value", []string{"a"}, 2, 3, []string{"a"}},
		{"split into batches", []string{"a", "b", "c", "d", "e"}, 2, 3, []string{"a,b", "c,d", "e"}},
		{"default to service limit", []string{"a", "b", "c", "d"}, 0, 3, []string{"a,b,c", "d"}},
		{"cap at service limit", []string{"a", "b", "c", "d"}, 5, 3, []string{"a,b,c", "d"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := concatValues(tt.values, tt.maxValues, tt.serviceMax)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("concatValues(%v, %d, %d) = %v, want %v", tt.values, tt.maxValues, tt.serviceMax, got, tt.want)
			}
		})
	}
}

func TestJoinURLPath(t *testing.T) {
	tests := []struct {
		base, path string
		want       string
	}{
		{"", "", "/"},
		{"base", "", "base/"},
		{"", "path", "/path"},
		{"base", "path", "base/path"},
		{"base/", "/path", "base/path"},
		{"base/", "path", "base/path"},
		{"base", "/path", "base/path"},
	}

	for _, tt := range tests {
		t.Run(tt.base+" + "+tt.path, func(t *testing.T) {
			if got := joinURLPath(tt.base, tt.path); got != tt.want {
				t.Errorf("joinURLPath(%q, %q) = %q, want %q", tt.base, tt.path, got, tt.want)
			}
		})
	}
}
