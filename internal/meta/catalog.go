package meta

import "errors"

// Catalog описывает минимальное мета-описание справочника DITA.
type Catalog struct {
	ID     int64
	Name   string
	Prefix string
}

// Validate проверяет обязательные свойства мета-описания справочника.
func (c Catalog) Validate() error {
	if c.Name == "" {
		return errors.New("catalog name is required")
	}
	if c.Prefix == "" {
		return errors.New("catalog prefix is required")
	}
	return nil
}
