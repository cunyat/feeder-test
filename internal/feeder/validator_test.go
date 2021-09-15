package feeder

import "testing"

func TestValidateSKU(t *testing.T) {
	tests := []struct {
		name  string
		sku   string
		valid bool
	}{
		{sku: "abcd-1234\n", valid: true, name: "lowercase valid sku"},
		{sku: "ASDF-1234\n", valid: true, name: "uppercase valid sku"},
		{sku: "abcd-1234\r\n", valid: true, name: "windows new-line sequence"},
		{sku: "a1cd-1234\n", valid: false, name: "number in first chunk"},
		{sku: "abcd-1a34\n", valid: false, name: "character in first chunk"},
		{sku: "abcd_1234\n", valid: false, name: "bad separator"},
		{sku: "abcd-1234", valid: false, name: "no new-line sequence"},
		{sku: "abcd-1234\r", valid: false, name: "bad new-line sequence"},
		{sku: "acd-1234\n", valid: false, name: "less characters"},
		{sku: "aecd-123\n", valid: false, name: "less numbers"},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSKU(tt.sku)
			if tt.valid && err != nil {
				t.Errorf("sku %s should be valid, error returned: %s", tt.sku, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("sku %s should be invalid, no error returned", tt.sku)
			}
		})
	}
}
