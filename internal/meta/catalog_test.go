package meta

import "testing"

func TestCatalogValidate(t *testing.T) {
	tests := []struct {
		name    string
		catalog Catalog
		wantErr bool
	}{
		{
			name: "valid catalog",
			catalog: Catalog{
				ID:     1,
				Name:   "Объекты эксплуатации",
				Prefix: "ОО",
			},
		},
		{
			name:    "name is required",
			catalog: Catalog{Prefix: "ОО"},
			wantErr: true,
		},
		{
			name:    "prefix is required",
			catalog: Catalog{Name: "Объекты эксплуатации"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.catalog.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
