package feeder

import (
	"fmt"
	"strings"
)

func ValidateSKU(value string) error {
	if len(value) != 10 && len(value) != 11 {
		return fmt.Errorf("invalid sku length (%d)", len(value))
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
		if ch != '\n'{
			return fmt.Errorf("sku must be finished with new-line sequence")
		}
	default:
		return fmt.Errorf("unexpected character index")
	}

	return nil
}

func IsTerminateSequence(value string) bool {
	return value == "terminate\r\n" || value == "terminate\n"
}
