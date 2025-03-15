package sduration

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Value implements the driver.Valuer interface.
// This method returns a driver.Value from the SDuration.
// It will be stored in the database as a string.
func (d SDuration) Value() (driver.Value, error) {
	return d.Duration().String(), nil
}

// Scan implements the sql.Scanner interface.
// This method scans a value from the database driver and
// attempts to convert it to an SDuration.
func (d *SDuration) Scan(src interface{}) error {
	if src == nil {
		*d = SDuration(0)
		return nil
	}

	// Handle different types that might come from the database
	switch v := src.(type) {
	case string:
		// Parse the duration string
		duration, err := time.ParseDuration(v)
		if err != nil {
			return fmt.Errorf("failed to parse duration from string '%s': %w", v, err)
		}

		*d = SDuration(duration)

		return nil

	case []byte:
		// Parse the duration from bytes (common in SQL drivers)
		duration, err := time.ParseDuration(string(v))

		if err != nil {
			return fmt.Errorf("failed to parse duration from bytes '%s': %w", string(v), err)
		}

		*d = SDuration(duration)

		return nil

	case int64:
		// Handle the case where the duration is stored as nanoseconds
		*d = SDuration(v)

		return nil

	case float64:
		// Handle the case where the duration is stored as a float
		*d = SDuration(time.Duration(v))

		return nil

	default:
		return fmt.Errorf("cannot scan type %T into SDuration", src)
	}
}
