package internal

import (
	"testing"
)

func TestFormatToLocal(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "standard afternoon time",
			input:   "2025-06-15T18:30Z",
			want:    "15. Jun, 18:30",
			wantErr: false,
		},
		{
			name:    "midnight edge case",
			input:   "2025-01-01T00:00Z",
			want:    "01. Jan, 00:00",
			wantErr: false,
		},
		{
			name:    "end of year",
			input:   "2025-12-31T23:59Z",
			want:    "31. Dec, 23:59",
			wantErr: false,
		},
		{
			name:    "leap day",
			input:   "2024-02-29T12:00Z",
			want:    "29. Feb, 12:00",
			wantErr: false,
		},
		{
			name:    "empty string returns error",
			input:   "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "wrong format returns error",
			input:   "June 15, 2025",
			want:    "",
			wantErr: true,
		},
		{
			name:    "missing Z suffix returns error",
			input:   "2025-06-15T18:30",
			want:    "",
			wantErr: true,
		},
		{
			name:    "full ISO 8601 with seconds returns error",
			input:   "2025-06-15T18:30:00Z",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatToLocal(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("FormatToLocal(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("FormatToLocal(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
