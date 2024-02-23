package formatter

import (
	"os"
	"path"
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
			want: &CSVFormatter{config: &Config{OutputFormat: CSVOutput}},
		},
		{
			name: "DOT output",
			args: args{
				config: &Config{
					OutputFormat: DotOutput,
				},
			},
			want: &DotFormatter{config: &Config{OutputFormat: DotOutput}},
		},
		{
			name: "SQLite output",
			args: args{
				config: &Config{
					OutputFormat: SqliteOutput,
				},
			},
			want: &SqliteFormatter{config: &Config{OutputFormat: SqliteOutput}},
		},
		{
			name: "Excel output",
			args: args{
				config: &Config{
					OutputFormat: ExcelOutput,
				},
			},
			want: &ExcelFormatter{config: &Config{OutputFormat: ExcelOutput}},
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

func TestTemplateContent(t *testing.T) {
	type args struct {
		f Formatter
		c *Config
	}
	beforeFunc := func(path string, t *testing.T) {
		f, err := os.Create(path)
		if err != nil {
			t.Errorf("failed to create a file: %v", err)
		}
		defer f.Close()
	}
	afterFunc := func(path string) {
		os.Remove(path)
	}
	tests := []struct {
		name      string
		args      args
		want      string
		wantErr   bool
		file      string
		before    func(path string, t *testing.T)
		after     func(path string)
		runBefore bool
		runAfter  bool
	}{
		{
			name: "Template file does not exists",
			args: args{
				f: &MarkdownFormatter{},
				c: &Config{
					TemplatePath: "non-existing-file-123",
				},
			},
			want:      "",
			wantErr:   true,
			runBefore: false,
			runAfter:  false,
		},
		{
			name: "Simple template file",
			args: args{
				f: &MarkdownFormatter{},
				c: &Config{
					TemplatePath: path.Join(os.TempDir(), "simple_template_md_file.md"),
				},
			},
			want:      "",
			wantErr:   false,
			file:      path.Join(os.TempDir(), "simple_template_md_file.md"),
			before:    beforeFunc,
			after:     afterFunc,
			runBefore: true,
			runAfter:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.runBefore {
				tt.before(tt.file, t)
			}
			if tt.runAfter {
				defer tt.after(tt.file)
			}
			got, err := TemplateContent(tt.args.f, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("TemplateContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TemplateContent() = %v, want %v", got, tt.want)
			}
		})
	}
}
