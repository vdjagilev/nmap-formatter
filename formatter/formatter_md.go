package formatter

import (
	_ "embed"
	"html/template"
	"strings"

	"github.com/vdjagilev/nmap-formatter/types"
)

type MarkdownFormatter struct {
	Config *Config
}

//go:embed resources/templates/markdown.tmpl
var MarkdownTemplate string

// Format the data and output it to appropriate io.Writer
func (f *MarkdownFormatter) Format(td *types.TemplateData) (err error) {
	tmpl := template.New("markdown")
	f.defineFunctions(tmpl)
	tmpl, err = tmpl.Parse(MarkdownTemplate)
	if err != nil {
		return
	}
	return tmpl.Execute(f.Config.Writer, td)
}

func (f *MarkdownFormatter) defineFunctions(tmpl *template.Template) {
	tmpl.Funcs(
		template.FuncMap{
			"md_toc": func(v string) string {
				// Removing dots, replacing spaces with hyphens,
				// then convert it to lower-case
				return strings.ToLower(
					strings.ReplaceAll(
						strings.ReplaceAll(
							v,
							".",
							"",
						),
						" ",
						"-",
					),
				)
			},
			"md": func(v string) string {
				// Removing all the tick symbols
				return strings.ReplaceAll(v, "`", "")
			},
			"noesc": func(v string) template.HTML {
				// Removing all tick symbols and displaying raw data
				return template.HTML(strings.ReplaceAll(v, "`", ""))
			},
		},
	)
}
