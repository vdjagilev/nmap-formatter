package formatter

import (
	"testing"
)

func TestOutputFormat_IsValid(t *testing.T) {
	tests := []struct {
		name string
		of   OutputFormat
		want bool
	}{
		{
			name: "Empty",
			of:   "",
			want: false,
		},
		{
			name: "Wrong",
			of:   "123",
			want: false,
		},
		{
			name: "Markdown 1",
			of:   "md",
			want: true,
		},
		{
			name: "Markdown 2",
			of:   "markdown",
			want: true,
		},
		{
			name: "HTML",
			of:   "html",
			want: true,
		},
		{
			name: "CSV",
			of:   "csv",
			want: true,
		},
		{
			name: "dot",
			of:   "dot",
			want: true,
		},
		{
			name: "sqlite",
			of:   "sqlite",
			want: true,
		},
		{
			name: "excel",
			of:   "sqlite",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.of.IsValid(); got != tt.want {
				t.Errorf("OutputFormat.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
