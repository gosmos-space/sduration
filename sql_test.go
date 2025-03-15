package sduration

import (
	"database/sql/driver"
	"testing"
	"time"
)

func TestValue(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := tt.duration.Value()
				if err != nil {
					t.Errorf("Value() error = %v", err)
					return
				}

				if got != tt.want {
					t.Errorf("Value() got = %v, want %v", got, tt.want)
				}

				// Verify the returned type implements driver.Value
				var _ driver.Value = got
			},
		)
	}
}

func TestScan(t *testing.T) {
	tests := []struct {
		name    string
		src     interface{}
		want    SDuration
		wantErr bool
	}{
		{
			name:    "string",
			src:     "5s",
			want:    SDuration(5 * time.Second),
			wantErr: false,
		},
		{
			name:    "string_complex",
			src:     "1h30m45s",
			want:    SDuration(1*time.Hour + 30*time.Minute + 45*time.Second),
			wantErr: false,
		},
		{
			name:    "bytes",
			src:     []byte("5s"),
			want:    SDuration(5 * time.Second),
			wantErr: false,
		},
		{
			name:    "int64",
			src:     int64(5 * time.Second),
			want:    SDuration(5 * time.Second),
			wantErr: false,
		},
		{
			name:    "float64",
			src:     float64(5 * time.Second),
			want:    SDuration(5 * time.Second),
			wantErr: false,
		},
		{
			name:    "nil",
			src:     nil,
			want:    SDuration(0),
			wantErr: false,
		},
		{
			name:    "invalid_string",
			src:     "not_a_duration",
			want:    SDuration(0),
			wantErr: true,
		},
		{
			name:    "invalid_type",
			src:     struct{}{},
			want:    SDuration(0),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				var d SDuration
				err := d.Scan(tt.src)

				if (err != nil) != tt.wantErr {
					t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if !tt.wantErr && d != tt.want {
					t.Errorf("Scan() got = %v, want %v", d.Duration(), tt.want.Duration())
				}
			},
		)
	}
}

// TestSQLRoundTrip tests the full roundtrip of Value -> Scan
func TestSQLRoundTrip(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// Value
				val, err := tt.duration.Value()
				if err != nil {
					t.Errorf("Value() error = %v", err)
					return
				}

				// Scan
				var scanned SDuration
				err = scanned.Scan(val)
				if err != nil {
					t.Errorf("Scan() error = %v", err)
					return
				}

				// Compare
				if scanned != tt.duration {
					t.Errorf(
						"Round-trip: got = %v, want %v",
						scanned.Duration(), tt.duration.Duration(),
					)
				}
			},
		)
	}
}
