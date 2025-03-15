package sduration

import "time"

// SDuration is a wrapper around time.Duration that allows for JSON/Text/SQL marshalling and unmarshalling.
type SDuration time.Duration

// Duration returns the time.Duration value of the SDuration.
func (d SDuration) Duration() time.Duration {
	return time.Duration(d)
}
