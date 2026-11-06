// Package configuration содержит доменную модель декларативной конфигурации DITA.
package configuration

import "encoding/json"

// FormatVersionV1 — первая поддерживаемая версия формата JSON-конфигурации.
const FormatVersionV1 = "1"

// ObjectKind определяет стандартный вид объекта платформы.
type ObjectKind string

const (
	ObjectKindCatalog  ObjectKind = "catalog"
	ObjectKindDocument ObjectKind = "document"
	ObjectKindPlanning ObjectKind = "planning"
	ObjectKindTask     ObjectKind = "task"
	ObjectKindNote     ObjectKind = "note"
)

// AttributeType определяет тип значения реквизита.
type AttributeType string

const (
	AttributeTypeString    AttributeType = "string"
	AttributeTypeText      AttributeType = "text"
	AttributeTypeInteger   AttributeType = "integer"
	AttributeTypeDecimal   AttributeType = "decimal"
	AttributeTypeBoolean   AttributeType = "boolean"
	AttributeTypeDate      AttributeType = "date"
	AttributeTypeDateTime  AttributeType = "datetime"
	AttributeTypeReference AttributeType = "reference"
)

// ReferenceDeleteRule задаёт разрешённое поведение физического удаления цели ссылки.
type ReferenceDeleteRule string

const ReferenceDeleteRestrict ReferenceDeleteRule = "restrict"

// Definition — корневой документ JSON-конфигурации DITA.
type Definition struct {
	FormatVersion string        `json:"format_version"`
	Configuration Configuration `json:"configuration"`
	ObjectTypes   []ObjectType  `json:"object_types"`
}

// Configuration описывает идентичность и версию прикладной конфигурации.
type Configuration struct {
	ID                    string `json:"id"`
	Code                  string `json:"code"`
	Name                  string `json:"name"`
	Version               string `json:"version"`
	PlatformCompatibility string `json:"platform_compatibility"`
}

// ObjectType описывает конкретный тип объекта внутри конфигурации.
type ObjectType struct {
	ID         string      `json:"id"`
	Kind       ObjectKind  `json:"kind"`
	Code       string      `json:"code"`
	Name       string      `json:"name"`
	Attributes []Attribute `json:"attributes"`

	// Структура kind-specific настроек проверяется Validator.
	Hierarchy   json.RawMessage `json:"hierarchy,omitempty"`
	Folders     json.RawMessage `json:"folders,omitempty"`
	OwnerTypeID string          `json:"owner_type_id,omitempty"`
	Numbering   json.RawMessage `json:"numbering,omitempty"`
	Posting     json.RawMessage `json:"posting,omitempty"`
	Indexes     []Index         `json:"indexes,omitempty"`
}

// Attribute описывает типизированный реквизит объектного типа.
type Attribute struct {
	ID           string              `json:"id"`
	Code         string              `json:"code"`
	Name         string              `json:"name"`
	Type         AttributeType       `json:"type"`
	Required     bool                `json:"required"`
	Length       *int                `json:"length,omitempty"`
	Precision    *int                `json:"precision,omitempty"`
	Scale        *int                `json:"scale,omitempty"`
	Unique       bool                `json:"unique,omitempty"`
	Index        bool                `json:"index,omitempty"`
	Default      json.RawMessage     `json:"default,omitempty"`
	Min          json.RawMessage     `json:"min,omitempty"`
	Max          json.RawMessage     `json:"max,omitempty"`
	MinLength    *int                `json:"min_length,omitempty"`
	MaxLength    *int                `json:"max_length,omitempty"`
	TargetTypeID string              `json:"target_type_id,omitempty"`
	OnDelete     ReferenceDeleteRule `json:"on_delete,omitempty"`
}

// Index описывает составной индекс без SQL-фрагментов в конфигурации.
type Index struct {
	ID          string   `json:"id"`
	Attributes  []string `json:"attributes"`
	Unique      bool     `json:"unique"`
	WhereActive bool     `json:"where_active,omitempty"`
}
