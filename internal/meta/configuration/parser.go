package configuration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"unicode/utf8"
)

// Parse разбирает JSON-конфигурацию DITA формата v1.
func Parse(data []byte) (Definition, error) {
	var definition Definition

	if !utf8.Valid(data) {
		return definition, fmt.Errorf("configuration: invalid UTF-8")
	}

	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || trimmed[0] != '{' {
		return definition, fmt.Errorf("configuration: root must be a JSON object")
	}

	// Проверяем синтаксис и отсутствие повторяющихся ключей.
	check := json.NewDecoder(bytes.NewReader(data))
	check.UseNumber()

	if err := scanJSONValue(check); err != nil {
		return definition, fmt.Errorf("configuration: invalid JSON: %w", err)
	}

	if _, err := check.Token(); err != io.EOF {
		if err == nil {
			return definition, fmt.Errorf("configuration: multiple JSON values")
		}
		return definition, fmt.Errorf("configuration: trailing data: %w", err)
	}

	// Разбираем структуру с запретом неизвестных полей.
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&definition); err != nil {
		return definition, fmt.Errorf("configuration: decode definition: %w", err)
	}

	return definition, nil
}

// scanJSONValue рекурсивно проверяет JSON и запрещает повторяющиеся ключи.
func scanJSONValue(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}

	switch value := token.(type) {
	case json.Delim:
		switch value {
		case '{':
			keys := make(map[string]struct{})

			for decoder.More() {
				keyToken, err := decoder.Token()
				if err != nil {
					return err
				}

				key, ok := keyToken.(string)
				if !ok {
					return fmt.Errorf("object key must be a string")
				}

				if _, exists := keys[key]; exists {
					return fmt.Errorf("duplicate object key %q", key)
				}
				keys[key] = struct{}{}

				if err := scanJSONValue(decoder); err != nil {
					return err
				}
			}

			end, err := decoder.Token()
			if err != nil {
				return err
			}
			if end != json.Delim('}') {
				return fmt.Errorf("expected closing object delimiter")
			}

		case '[':
			for decoder.More() {
				if err := scanJSONValue(decoder); err != nil {
					return err
				}
			}

			end, err := decoder.Token()
			if err != nil {
				return err
			}
			if end != json.Delim(']') {
				return fmt.Errorf("expected closing array delimiter")
			}

		default:
			return fmt.Errorf("unexpected delimiter %q", value)
		}
	}

	return nil
}
