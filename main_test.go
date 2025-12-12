package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCleanInput(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected []string
	}{
		"valid":                      {"a b c", []string{"a", "b", "c"}},
		"validWithLUpperCaseStrings": {"ABC ABC", []string{"abc", "abc"}},
		"noSep":                      {"abc", []string{"abc"}},
		"sepAtStart":                 {" a b c", []string{"a", "b", "c"}},
		"sepAtEnd":                   {"a b c ", []string{"a", "b", "c"}},
		"sepAtStartAndEnd":           {" a b c ", []string{"a", "b", "c"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := cleanInput(tc.input)
			diff := cmp.Diff(tc.expected, got)
			if diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
