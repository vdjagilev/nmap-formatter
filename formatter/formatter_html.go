package formatter

import (
	_ "embed"
	"html/template"
)

type HTMLFormatter struct {
	config *Config
}

//go:embed resources/templates/simple-html.gohtml
var HTMLSimpleTemplate string

// Format the data and output it to appropriate io.Writer
func (f *HTMLFormatter) Format(td *TemplateData) error {
	tmpl, err := template.New("html").Parse(HTMLSimpleTemplate)
	if err != nil {
		return err
	}
	return tmpl.Execute(f.config.Writer, td)
}
