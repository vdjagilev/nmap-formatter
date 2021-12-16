package formatter

import (
	// Used in this place to have all required functionality within one binary file. No need for separate folders/files, just embed template
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"strings"
)

// MarkdownFormatter is a formatter struct used to deliver markdown file format
type MarkdownFormatter struct {
	config *Config
}

//go:embed resources/templates/markdown.tmpl
// MarkdownTemplate variable is used to store markdown.tmpl embed file contents
var MarkdownTemplate string

type markdownOutputFilter struct {
	writer  io.Writer
	content []byte
}

func (m *markdownOutputFilter) Write(p []byte) (n int, err error) {
	m.content = append(m.content, p...)
	return len(p), nil
}

// split is used to split markdown content with new line delimiters
func (m *markdownOutputFilter) split() [][]byte {
	lines := [][]byte{}
	line := []byte{}
	for _, b := range m.content {
		if b == '\n' {
			// No need to add new line if previous one was defined as new line
			if len(line) != 0 {
				lines = append(lines, line)
			}
			lines = append(lines, []byte{})
			line = []byte{}
			continue
		}
		line = append(line, b)
	}
	if len(line) != 0 {
		lines = append(lines, line)
	}
	return lines
}

func (m *markdownOutputFilter) filter() []byte {
	content := []byte{}
	contentLines := m.split()
	newLines := 0
	isCodeBlock := false
	for _, l := range contentLines {
		if len(l) >= 3 &&
			l[0] == '`' &&
			l[1] == '`' &&
			l[2] == '`' {
			isCodeBlock = !isCodeBlock
		}

		if !isCodeBlock {
			if len(l) == 0 {
				newLines++
			} else {
				newLines = 0
			}

			// Skip other new lines
			if newLines > 2 {
				continue
			}

			if newLines > 0 {
				l = append(l, '\n')
			}
		} else {
			// Simply add newline if it's empty slice
			if len(l) == 0 {
				l = append(l, '\n')
			}
		}
		content = append(content, l...)
	}
	return content
}

// Format the data and output it to appropriate io.Writer
func (f *MarkdownFormatter) Format(td *TemplateData) (err error) {
	tmpl := template.New("markdown")
	f.defineTemplateFunctions(tmpl)
	tmpl, err = tmpl.Parse(MarkdownTemplate)
	if err != nil {
		return
	}
	markdownOutput := &markdownOutputFilter{writer: f.config.Writer, content: []byte{}}
	err = tmpl.Execute(markdownOutput, td)
	if err != nil {
		return err
	}
	_, err = f.config.Writer.Write(markdownOutput.filter())
	return err
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
