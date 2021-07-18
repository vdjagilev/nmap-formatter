package formatter

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name string
		args args
		want Formatter
	}{
		{
			name: "Nil case",
			args: args{
				config: &Config{
					OutputFormat: OutputFormat("nil"),
				},
			},
			want: nil,
		},
		{
			name: "JSON output",
			args: args{
				config: &Config{
					OutputFormat: JSONOutput,
				},
			},
			want: &JSONFormatter{config: &Config{OutputFormat: JSONOutput}},
		},
		{
			name: "HTML output",
			args: args{
				config: &Config{
					OutputFormat: HTMLOutput,
				},
			},
			want: &HTMLFormatter{config: &Config{OutputFormat: HTMLOutput}},
		},
		{
			name: "Markdown output",
			args: args{
				config: &Config{
					OutputFormat: MarkdownOutput,
				},
			},
			want: &MarkdownFormatter{config: &Config{OutputFormat: MarkdownOutput}},
		},
		{
			name: "CSV output",
			args: args{
				config: &Config{
					OutputFormat: CSVOutput,
				},
			},
			want: &CSVFormatter{Config: &Config{OutputFormat: CSVOutput}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
