package configuration

import (
	"encoding/json"
	"testing"
)

func TestDefinitionUnmarshal(t *testing.T) {
	const source = `{
		"format_version": "1",
		"configuration": {
			"id": "8f4f7f8a-7d34-4e2f-9c6a-5a5c4c1a1001",
			"code": "example",
			"name": "Example configuration",
			"version": "1.0.0",
			"platform_compatibility": ">=0.2.0 <0.3.0"
		},
		"object_types": [{
			"id": "8f4f7f8a-7d34-4e2f-9c6a-5a5c4c1a2001",
			"kind": "catalog",
			"code": "locations",
			"name": "Locations",
			"attributes": [{
				"id": "8f4f7f8a-7d34-4e2f-9c6a-5a5c4c1a3001",
				"code": "name",
				"name": "Name",
				"type": "string",
				"length": 255,
				"required": true
			}]
		}]
	}`

	var got Definition
	if err := json.Unmarshal([]byte(source), &got); err != nil {
		t.Fatalf("unmarshal definition: %v", err)
	}
	if got.FormatVersion != FormatVersionV1 {
		t.Fatalf("format version = %q, want %q", got.FormatVersion, FormatVersionV1)
	}
	if len(got.ObjectTypes) != 1 || got.ObjectTypes[0].Kind != ObjectKindCatalog {
		t.Fatalf("unexpected object types: %+v", got.ObjectTypes)
	}
	if len(got.ObjectTypes[0].Attributes) != 1 {
		t.Fatalf("attribute count = %d, want 1", len(got.ObjectTypes[0].Attributes))
	}
	attribute := got.ObjectTypes[0].Attributes[0]
	if attribute.Type != AttributeTypeString || attribute.Length == nil || *attribute.Length != 255 || !attribute.Required {
		t.Fatalf("unexpected attribute: %+v", attribute)
	}
}

func TestStandardObjectKinds(t *testing.T) {
	got := []ObjectKind{ObjectKindCatalog, ObjectKindDocument, ObjectKindPlanning, ObjectKindTask, ObjectKindNote}
	want := []ObjectKind{"catalog", "document", "planning", "task", "note"}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("kind[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestStandardAttributeTypes(t *testing.T) {
	got := []AttributeType{AttributeTypeString, AttributeTypeText, AttributeTypeInteger, AttributeTypeDecimal, AttributeTypeBoolean, AttributeTypeDate, AttributeTypeDateTime, AttributeTypeReference}
	want := []AttributeType{"string", "text", "integer", "decimal", "boolean", "date", "datetime", "reference"}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("attribute type[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
