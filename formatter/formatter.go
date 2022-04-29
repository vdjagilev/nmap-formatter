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
	}
	return nil
}

// Formatter interface describes only one function `Format()` that is responsible for data "formatting"
type Formatter interface {
	// Format the data and output it to appropriate io.Writer
	Format(td *TemplateData, templateContent string) error
	// TemplateContent() (string, error)
	defaultTemplateContent() string
}

func TemplateContent(f Formatter, c *Config) (string, error) {
	if c.TemplatePath != "" {
		f, err := os.Open(c.TemplatePath)
		if err != nil {
			return "", err
		}
		defer f.Close()
		content, err := io.ReadAll(f)
		if err != nil {
			return "", err
		}
		return string(content), nil
	}
	return f.defaultTemplateContent(), nil
}
