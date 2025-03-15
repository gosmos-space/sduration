package sduration

import (
	"testing"
	"time"
)

func TestUnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		want    SDuration
		wantErr bool
	}{
		{
			name:    "seconds",
			text:    "5s",
			want:    SDuration(5 * time.Second),
			wantErr: false,
		},
		{
			name:    "minutes",
			text:    "10m",
			want:    SDuration(10 * time.Minute),
			wantErr: false,
		},
		{
			name:    "hours",
			text:    "2h",
			want:    SDuration(2 * time.Hour),
			wantErr: false,
		},
		{
			name:    "complex",
			text:    "1h30m45s",
			want:    SDuration(1*time.Hour + 30*time.Minute + 45*time.Second),
			wantErr: false,
		},
		{
			name:    "negative",
			text:    "-30s",
			want:    SDuration(-30 * time.Second),
			wantErr: false,
		},
		{
			name:    "milliseconds",
			text:    "250ms",
			want:    SDuration(250 * time.Millisecond),
			wantErr: false,
		},
		{
			name:    "microseconds",
			text:    "500µs",
			want:    SDuration(500 * time.Microsecond),
			wantErr: false,
		},
		{
			name:    "nanoseconds",
			text:    "50ns",
			want:    SDuration(50 * time.Nanosecond),
			wantErr: false,
		},
		{
			name:    "zero",
			text:    "0s",
			want:    SDuration(0),
			wantErr: false,
		},
		{
			name:    "empty string",
			text:    "",
			want:    SDuration(0),
			wantErr: true,
		},
		{
			name:    "invalid format",
			text:    "not a duration",
			want:    SDuration(0),
			wantErr: true,
		},
		{
			name:    "missing unit",
			text:    "42",
			want:    SDuration(0),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				var d SDuration
				err := d.UnmarshalText([]byte(tt.text))

				if (err != nil) != tt.wantErr {
					t.Errorf("UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if !tt.wantErr && d != tt.want {
					t.Errorf("UnmarshalText() got = %v, want %v", d, tt.want)
				}
			},
		)
	}
}

func TestMarshalText(t *testing.T) {
	tests := []struct {
		name     string
		duration SDuration
		want     string
	}{
		{
			name:     "zero",
			duration: SDuration(0),
			want:     "0s",
		},
		{
			name:     "seconds",
			duration: SDuration(5 * time.Second),
			want:     "5s",
		},
		{
			name:     "minutes",
			duration: SDuration(10 * time.Minute),
			want:     "10m0s",
		},
		{
			name:     "hours",
			duration: SDuration(2 * time.Hour),
			want:     "2h0m0s",
		},
		{
			name:     "complex",
			duration: SDuration(1*time.Hour + 30*time.Minute + 45*time.Second),
			want:     "1h30m45s",
		},
		{
			name:     "negative",
			duration: SDuration(-30 * time.Second),
			want:     "-30s",
		},
		{
			name:     "milliseconds",
			duration: SDuration(250 * time.Millisecond),
			want:     "250ms",
		},
		{
			name:     "microseconds",
			duration: SDuration(500 * time.Microsecond),
			want:     "500µs",
		},
		{
			name:     "nanoseconds",
			duration: SDuration(50 * time.Nanosecond),
			want:     "50ns",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := tt.duration.MarshalText()
				if err != nil {
					t.Errorf("MarshalText() error = %v", err)
					return
				}

				if string(got) != tt.want {
					t.Errorf("MarshalText() got = %v, want %v", string(got), tt.want)
				}
			},
		)
	}
}

// TestRoundTripMarshalUnmarshal tests the complete cycle of marshaling
// and then unmarshaling to ensure consistency
func TestRoundTripMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		duration SDuration
	}{
		{
			name:     "zero",
			duration: SDuration(0),
		},
		{
			name:     "seconds",
			duration: SDuration(5 * time.Second),
		},
		{
			name:     "complex",
			duration: SDuration(1*time.Hour + 30*time.Minute + 45*time.Second),
		},
		{
			name:     "negative",
			duration: SDuration(-30 * time.Second),
		},
		{
			name:     "milliseconds",
			duration: SDuration(250 * time.Millisecond),
		},
		{
			name:     "microseconds",
			duration: SDuration(500 * time.Microsecond),
		},
		{
			name:     "nanoseconds",
			duration: SDuration(50 * time.Nanosecond),
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// Marshal the duration to text
				text, err := tt.duration.MarshalText()
				if err != nil {
					t.Errorf("MarshalText() error = %v", err)
					return
				}

				// Unmarshal the text back to a duration
				var unmarshaled SDuration
				err = unmarshaled.UnmarshalText(text)
				if err != nil {
					t.Errorf("UnmarshalText() error = %v", err)
					return
				}

				// Compare the original and round-tripped durations
				if unmarshaled != tt.duration {
					t.Errorf(
						"Round-trip: got = %v, want %v",
						time.Duration(unmarshaled), time.Duration(tt.duration),
					)
				}
			},
		)
	}
}
