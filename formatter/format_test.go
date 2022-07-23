package formatter

import (
	"reflect"
	"testing"
)

func TestOutputFormat_FileOutputFormat(t *testing.T) {
	tests := []struct {
		name string
		of   OutputFormat
		want OutputFormat
	}{
		{
			name: "Default variant",
			of:   "",
			want: "html",
		},
		{
			name: "Markdown 1",
			of:   "md",
			want: "md",
		},
		{
			name: "Markdown 2",
			of:   "markdown",
			want: "md",
		},
		{
			name: "HTML",
			of:   "html",
			want: "html",
		},
		{
			name: "CSV",
			of:   "csv",
			want: "csv",
		},
		{
			name: "JSON",
			of:   "json",
			want: "json",
		},
		{
			name: "dot",
			of:   "dot",
			want: "dot",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.of.FileOutputFormat(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OutputFormat.FileOutputFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.of.IsValid(); got != tt.want {
				t.Errorf("OutputFormat.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
