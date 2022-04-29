package formatter

import (
	// Used in this place to have all required functionality within one binary file. No need for separate folders/files, just embed template
	_ "embed"
	"html/template"
)

// HTMLFormatter is struct defined for HTML Output use-case
type HTMLFormatter struct {
	config *Config
}

//go:embed resources/templates/simple-html.gohtml
// HTMLSimpleTemplate variable is used to store embedded HTML template content
var HTMLSimpleTemplate string

// Format the data and output it to appropriate io.Writer
func (f *HTMLFormatter) Format(td *TemplateData, templateContent string) error {
	tmpl, err := template.New("html").Parse(templateContent)
	if err != nil {
		return err
	}
	return tmpl.Execute(f.config.Writer, td)
}

func (f *HTMLFormatter) defaultTemplateContent() string {
	return HTMLSimpleTemplate
}
