package sduration

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		duration SDuration
		want     string
	}{
		{
			name:     "zero",
			duration: SDuration(0),
			want:     `"0s"`,
		},
		{
			name:     "seconds",
			duration: SDuration(5 * time.Second),
			want:     `"5s"`,
		},
		{
			name:     "minutes",
			duration: SDuration(10 * time.Minute),
			want:     `"10m0s"`,
		},
		{
			name:     "hours",
			duration: SDuration(2 * time.Hour),
			want:     `"2h0m0s"`,
		},
		{
			name:     "complex",
			duration: SDuration(1*time.Hour + 30*time.Minute + 45*time.Second),
			want:     `"1h30m45s"`,
		},
		{
			name:     "negative",
			duration: SDuration(-30 * time.Second),
			want:     `"-30s"`,
		},
		{
			name:     "milliseconds",
			duration: SDuration(250 * time.Millisecond),
			want:     `"250ms"`,
		},
		{
			name:     "microseconds",
			duration: SDuration(500 * time.Microsecond),
			want:     `"500µs"`,
		},
		{
			name:     "nanoseconds",
			duration: SDuration(50 * time.Nanosecond),
			want:     `"50ns"`,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := tt.duration.MarshalJSON()
				if err != nil {
					t.Errorf("MarshalJSON() error = %v", err)
					return
				}

				if string(got) != tt.want {
					t.Errorf("MarshalJSON() got = %v, want %v", string(got), tt.want)
				}
			},
		)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    SDuration
		wantErr bool
	}{
		{
			name:    "seconds",
			json:    `"5s"`,
			want:    SDuration(5 * time.Second),
			wantErr: false,
		},
		{
			name:    "minutes",
			json:    `"10m"`,
			want:    SDuration(10 * time.Minute),
			wantErr: false,
		},
		{
			name:    "hours",
			json:    `"2h"`,
			want:    SDuration(2 * time.Hour),
			wantErr: false,
		},
		{
			name:    "complex",
			json:    `"1h30m45s"`,
			want:    SDuration(1*time.Hour + 30*time.Minute + 45*time.Second),
			wantErr: false,
		},
		{
			name:    "negative",
			json:    `"-30s"`,
			want:    SDuration(-30 * time.Second),
			wantErr: false,
		},
		{
			name:    "milliseconds",
			json:    `"250ms"`,
			want:    SDuration(250 * time.Millisecond),
			wantErr: false,
		},
		{
			name:    "microseconds",
			json:    `"500µs"`,
			want:    SDuration(500 * time.Microsecond),
			wantErr: false,
		},
		{
			name:    "nanoseconds",
			json:    `"50ns"`,
			want:    SDuration(50 * time.Nanosecond),
			wantErr: false,
		},
		{
			name:    "zero",
			json:    `"0s"`,
			want:    SDuration(0),
			wantErr: false,
		},
		{
			name:    "empty string",
			json:    `""`,
			want:    SDuration(0),
			wantErr: true,
		},
		{
			name:    "invalid format",
			json:    `"not a duration"`,
			want:    SDuration(0),
			wantErr: true,
		},
		{
			name:    "missing unit",
			json:    `"42"`,
			want:    SDuration(0),
			wantErr: true,
		},
		{
			name:    "missing quotes",
			json:    `5s`,
			want:    SDuration(0),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				var d SDuration
				err := d.UnmarshalJSON([]byte(tt.json))

				if (err != nil) != tt.wantErr {
					t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if !tt.wantErr && d != tt.want {
					t.Errorf("UnmarshalJSON() got = %v, want %v", d, tt.want)
				}
			},
		)
	}
}

// TestJSONStructMarshalUnmarshal tests the complete cycle of marshaling and
// unmarshaling SDuration within a struct in a JSON context
func TestJSONStructMarshalUnmarshal(t *testing.T) {
	type Config struct {
		Timeout SDuration `json:"timeout"`
		Retry   SDuration `json:"retry"`
	}

	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "basic_durations",
			config: Config{
				Timeout: SDuration(30 * time.Second),
				Retry:   SDuration(5 * time.Minute),
			},
		},
		{
			name: "complex_durations",
			config: Config{
				Timeout: SDuration(1*time.Hour + 30*time.Minute + 45*time.Second),
				Retry:   SDuration(250 * time.Millisecond),
			},
		},
		{
			name: "zero_values",
			config: Config{
				Timeout: SDuration(0),
				Retry:   SDuration(0),
			},
		},
		{
			name: "negative_values",
			config: Config{
				Timeout: SDuration(-10 * time.Second),
				Retry:   SDuration(-500 * time.Millisecond),
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// Marshal to JSON
				jsonData, err := json.Marshal(tt.config)
				if err != nil {
					t.Errorf("json.Marshal() error = %v", err)
					return
				}

				// Unmarshal back
				var decoded Config
				err = json.Unmarshal(jsonData, &decoded)
				if err != nil {
					t.Errorf("json.Unmarshal() error = %v", err)
					return
				}

				// Compare
				if decoded.Timeout != tt.config.Timeout {
					t.Errorf(
						"Timeout field mismatch: got = %v, want %v",
						decoded.Timeout.Duration(), tt.config.Timeout.Duration(),
					)
				}
				if decoded.Retry != tt.config.Retry {
					t.Errorf(
						"Retry field mismatch: got = %v, want %v",
						decoded.Retry.Duration(), tt.config.Retry.Duration(),
					)
				}
			},
		)
	}
}

// TestJSONRoundTrip tests marshaling and unmarshaling a single SDuration value directly
func TestJSONRoundTrip(t *testing.T) {
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
				// Marshal the duration to JSON
				jsonData, err := json.Marshal(tt.duration)
				if err != nil {
					t.Errorf("json.Marshal() error = %v", err)
					return
				}

				// Unmarshal the JSON back to a duration
				var unmarshaled SDuration
				err = json.Unmarshal(jsonData, &unmarshaled)
				if err != nil {
					t.Errorf("json.Unmarshal() error = %v", err)
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
