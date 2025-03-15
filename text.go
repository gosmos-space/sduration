package sduration

import "time"

func (d *SDuration) UnmarshalText(text []byte) error {
	duration, err := time.ParseDuration(string(text))

	if err != nil {
		return err
	}

	*d = SDuration(duration)

	return nil
}

func (d SDuration) MarshalText() (text []byte, err error) {
	return []byte(d.Duration().String()), nil
}
