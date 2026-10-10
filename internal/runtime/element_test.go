package runtime

import (
	"testing"

	"github.com/Dmitriy-495/dita/internal/meta"
)

func TestNewElement(t *testing.T) {
	catalog := meta.Catalog{
		ID:     1,
		Name:   "Объекты эксплуатации",
		Prefix: "ОО",
	}

	got, err := NewElement(catalog, 1, "Котельная №1")
	if err != nil {
		t.Fatalf("NewElement() error = %v", err)
	}

	want := Element{
		ID:   1,
		Code: "ОО-000001",
		Name: "Котельная №1",
	}

	if got != want {
		t.Fatalf("NewElement() = %+v, want %+v", got, want)
	}
}

func TestNewElementRejectsInvalidData(t *testing.T) {
	catalog := meta.Catalog{
		ID:     1,
		Name:   "Объекты эксплуатации",
		Prefix: "ОО",
	}

	tests := []struct {
		name string
		id   int64
		item string
	}{
		{name: "zero id", id: 0, item: "Котельная"},
		{name: "negative id", id: -1, item: "Котельная"},
		{name: "empty name", id: 1, item: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewElement(catalog, tt.id, tt.item); err == nil {
				t.Fatal("NewElement() expected error")
			}
		})
	}
}
