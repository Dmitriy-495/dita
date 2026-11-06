package configuration

import (
	"strings"
	"testing"
)

const validDefinition = `{
	"format_version": "1",
	"configuration": {
		"id": "8f4f7f8a-7d34-4e2f-9c6a-5a5c4c1a1001",
		"code": "example",
		"name": "Example configuration",
		"version": "1.0.0",
		"platform_compatibility": ">=0.2.0 <0.3.0"
	},
	"object_types": []
}`

func TestParseValidDefinition(t *testing.T) {
	got, err := Parse([]byte(validDefinition))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if got.FormatVersion != FormatVersionV1 {
		t.Errorf("format_version = %q, want %q", got.FormatVersion, FormatVersionV1)
	}

	if got.Configuration.Code != "example" {
		t.Errorf("configuration.code = %q, want example", got.Configuration.Code)
	}
}

func TestParseRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty input", ""},
		{"invalid JSON", `{"format_version":`},
		{"root array", `[]`},
		{"root null", `null`},
		{"unknown root field", `{"unknown": true}`},
		{
			"unknown nested field",
			`{"format_version":"1","configuration":{"unknown":true},"object_types":[]}`,
		},
		{"multiple JSON values", validDefinition + ` {}`},
		{
			"duplicate root key",
			`{"format_version":"1","format_version":"1","configuration":{},"object_types":[]}`,
		},
		{
			"duplicate nested key",
			`{"format_version":"1","configuration":{"code":"a","code":"b"},"object_types":[]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse([]byte(tt.input))
			if err == nil {
				t.Fatal("Parse() expected an error")
			}
		})
	}
}

func TestParseRejectsInvalidUTF8(t *testing.T) {
	input := append([]byte(validDefinition), 0xff)

	_, err := Parse(input)
	if err == nil {
		t.Fatal("Parse() expected invalid UTF-8 error")
	}

	if !strings.Contains(err.Error(), "UTF-8") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseAllowsUnknownFieldsInsideRawMessage(t *testing.T) {
	input := `{
		"format_version": "1",
		"configuration": {
			"id": "8f4f7f8a-7d34-4e2f-9c6a-5a5c4c1a1001",
			"code": "example",
			"name": "Example",
			"version": "1.0.0",
			"platform_compatibility": ">=0.2.0 <0.3.0"
		},
		"object_types": [{
			"id": "8f4f7f8a-7d34-4e2f-9c6a-5a5c4c1a2001",
			"kind": "catalog",
			"code": "locations",
			"name": "Locations",
			"attributes": [],
			"hierarchy": {
				"parent_attribute": "parent",
				"future_option": true
			}
		}]
	}`

	if _, err := Parse([]byte(input)); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
}
