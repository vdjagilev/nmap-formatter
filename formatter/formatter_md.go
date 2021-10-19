package formatter

import (
	// Used in this place to have all required functionality within one binary file. No need for separate folders/files, just embed template
	_ "embed"
	"fmt"
	"html/template"
	"strings"
)

// MarkdownFormatter is a formatter struct used to deliver markdown file format
type MarkdownFormatter struct {
	config *Config
}

//go:embed resources/templates/markdown.tmpl
// MarkdownTemplate variable is used to store markdown.tmpl embed file contents
var MarkdownTemplate string

// Format the data and output it to appropriate io.Writer
func (f *MarkdownFormatter) Format(td *TemplateData) (err error) {
	tmpl := template.New("markdown")
	f.defineTemplateFunctions(tmpl)
	tmpl, err = tmpl.Parse(MarkdownTemplate)
	if err != nil {
		return
	}
	return tmpl.Execute(f.config.Writer, td)
}

func (f *MarkdownFormatter) defineTemplateFunctions(tmpl *template.Template) {
	tmpl.Funcs(
		template.FuncMap{
			"md_toc":   markdownTOCEntry,
			"md":       markdownEntry,
			"noesc":    markdownNoEscape,
			"md_title": markdownHostAnchorTitle,
			"md_link":  markdownAnchorLink,
		},
	)
}

// markdownHostAnchorTitle helps to generate a title for specific hostname
func markdownHostAnchorTitle(h *Host) string {
	title := h.HostAddress.Address
	for i := range h.HostNames.HostName {
		title += fmt.Sprintf(" / %s", h.HostNames.HostName[i].Name)
	}
	title += fmt.Sprintf(" (%s)", h.Status.State)
	return title
}

// markdownAnchorLink is converting generated anchor title to table-of-contents entry
func markdownAnchorLink(h *Host) string {
	return markdownTOCEntry(markdownHostAnchorTitle(h))
}

// markdownTOCEntry returns lower-cased Table-of-Contents
// anchor entry that should work as an internal link
func markdownTOCEntry(v string) string {
	r := strings.NewReplacer(
		".", "",
		" ", "-",
		"/", "",
		"(", "",
		")", "",
	)
	// replace special characters and lower-case
	return strings.ToLower(r.Replace(v))
}

func markdownEntry(v string) string {
	// Removing all the tick symbols
	return strings.ReplaceAll(v, "`", "")
}

func markdownNoEscape(v string) template.HTML {
	// Removing all tick symbols and displaying raw data
	return template.HTML(markdownEntry(v))
}
