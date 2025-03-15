package sduration

import (
	"encoding/json"
	"time"
)

// MarshalJSON implements the json.Marshaler interface.
// It returns the duration as a string in JSON format.
func (d SDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Duration().String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It expects a string in the format that time.ParseDuration can handle.
func (d *SDuration) UnmarshalJSON(bytes []byte) error {
	var s string

	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}

	duration, err := time.ParseDuration(s)

	if err != nil {
		return err
	}

	*d = SDuration(duration)

	return nil
}
