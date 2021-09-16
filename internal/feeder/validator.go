package feeder

import (
	"fmt"
	"strings"
)

// ValidateSKU return an error if value it's not a valid sku
func ValidateSKU(value string) error {
	if len(value) != 10 && len(value) != 11 {
		return fmt.Errorf("invalid sku length (%d)", len(value))
	}

	// Validate ends with new line (\n)
	if value[len(value)-1] != '\n' {
		return fmt.Errorf("sku must be finished with new-line sequence")
	}

	// convert to lowercase to compare
	value = strings.ToLower(value)
	for i, char := range value {
		err := validateChar(i, char)
		if err != nil {
			return err
		}
	}

	return nil
}

// validateChar validates if a character is valid for the given position (i)
func validateChar(i int, ch int32) error {
	switch i {
	case 0, 1, 2, 3:
		if ch < 'a' || ch > 'z' {
			return fmt.Errorf("first 4 characters must be a letter")
		}
	case 4:
		if ch != '-' {
			return fmt.Errorf("separator character must be a dash")
		}
	case 5, 6, 7, 8:
		if ch < '0' || ch > '9' {
			return fmt.Errorf("last 4 characters must be a number")
		}
	case 9:
		if ch != '\r' && ch != '\n' {
			return fmt.Errorf("sku must be finished with new-line sequence")
		}
	case 10:
		if ch != '\n' {
			return fmt.Errorf("sku must be finished with new-line sequence")
		}
	default:
		return fmt.Errorf("unexpected character index")
	}

	return nil
}

// IsTerminateSequence return if the given value matches a terminate sequence
func IsTerminateSequence(value string) bool {
	return value == "terminate\r\n" || value == "terminate\n"
}
