package qdrant_test

import (
	"testing"

	"github.com/qdrant/go-client/qdrant"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected *qdrant.Version
		hasError bool
	}{
		{"v1.2.3", &qdrant.Version{Major: 1, Minor: 2, Rest: "3"}, false},
		{"1.2.3", &qdrant.Version{Major: 1, Minor: 2, Rest: "3"}, false},
		{"v1.2", &qdrant.Version{Major: 1, Minor: 2, Rest: ""}, false},
		{"1.2", &qdrant.Version{Major: 1, Minor: 2, Rest: ""}, false},
		{"v1.2.3.4", &qdrant.Version{Major: 1, Minor: 2, Rest: "3.4"}, false},
		{"", nil, true},
		{"1", nil, true},
		{"1.", nil, true},
		{".1.", nil, true},
		{"1.something.1", nil, true},
	}

	for _, test := range tests {
		result, err := qdrant.ParseVersion(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("ParseVersion(%q) error = %v, wantErr %v", test.input, err, test.hasError)
			continue
		}
		if !test.hasError && result != nil && (result.Major != test.expected.Major ||
			result.Minor != test.expected.Minor ||
			result.Rest != test.expected.Rest) {
			t.Errorf("ParseVersion(%q) = %v, want %v", test.input, result, test.expected)
		}
	}
}

func TestIsCompatible(t *testing.T) {
	tests := []struct {
		clientVersion string
		serverVersion string
		expected      bool
	}{
		{"1.9.3.dev0", "2.8.1.dev12-something", false},
		{"1.9", "2.8", false},
		{"1", "2", false},
		{"1", "1", true},
		{"1.9.0", "2.9.0", false},
		{"1.1.0", "1.2.9", true},
		{"1.2.7", "1.1.8.dev0", true},
		{"1.2.1", "1.2.29", true},
		{"1.2.0", "1.2.0", true},
		{"1.2.0", "1.4.0", false},
		{"1.4.0", "1.2.0", false},
		{"1.9.0", "3.7.0", false},
		{"3.0.0", "1.0.0", false},
	}

	for _, test := range tests {
		clientVersion := test.clientVersion
		serverVersion := test.serverVersion
		result := qdrant.IsCompatible(&clientVersion, &serverVersion)
		if result != test.expected {
			t.Errorf("IsCompatible(%q, %q) = %v, want %v", clientVersion, serverVersion, result, test.expected)
		}
	}
}
