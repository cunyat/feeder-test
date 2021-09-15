package feeder

import "testing"

func TestValidateSKU(t *testing.T) {
	tests := []struct {
		name  string
		sku   string
		valid bool
	}{
		{valid: true, sku: "abcd-1234\n", name: "lowercase valid sku"},
		{valid: true, sku: "ASDF-4925\n", name: "uppercase valid sku"},
		{valid: true, sku: "abcd-2794\r\n", name: "windows new-line sequence"},
		{valid: false, sku: "a1cd-1234\n", name: "number in first chunk"},
		{valid: false, sku: "abcd-1a34\n", name: "character in first chunk"},
		{valid: false, sku: "abcd_1234\n", name: "bad separator"},
		{valid: false, sku: "abcd-1234", name: "no new-line sequence"},
		{valid: false, sku: "abcd-1234\r", name: "bad new-line sequence"},
		{valid: false, sku: "acd-1234\n", name: "less characters"},
		{valid: false, sku: "aecd-123\n", name: "less numbers"},
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
