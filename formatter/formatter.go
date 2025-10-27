package formatter

import (
	"io"
	"os"
)

// New returns new instance of formatter the exact struct
// of formatter would depend on provided config
func New(config *Config) Formatter {
	switch config.OutputFormat {
	case JSONOutput:
		return &JSONFormatter{
			config,
		}
	case HTMLOutput:
		return &HTMLFormatter{
			config,
		}
	case MarkdownOutput:
		return &MarkdownFormatter{
			config,
		}
	case CSVOutput:
		return &CSVFormatter{
			config,
		}
	case ExcelOutput:
		return &ExcelFormatter{
			config,
		}
	case DotOutput:
		return &DotFormatter{
			config,
		}
	case SqliteOutput:
		return &SqliteFormatter{
			config,
		}
	case D2LangOutput:
		return &D2LangFormatter{
			config,
		}
	}
	return nil
}

// Formatter interface describes only one function `Format()` that is responsible for data "formatting"
type Formatter interface {
	// Format the data and output it to appropriate io.Writer
	Format(td *TemplateData, templateContent string) error
	// defaultTemplateContent returns default template content for any typical chosen formatter (HTML or Markdown)
	defaultTemplateContent() string
}

// TemplateContent reads customly provided template content or fails with error
func TemplateContent(f Formatter, c *Config) (string, error) {
	if c.TemplatePath != "" {
		file, err := os.Open(c.TemplatePath)
		if err != nil {
			return "", err
		}
		defer func() {
			_ = file.Close()
		}()
		content, err := io.ReadAll(file)
		if err != nil {
			return "", err
		}
		return string(content), nil
	}
	return f.defaultTemplateContent(), nil
}
