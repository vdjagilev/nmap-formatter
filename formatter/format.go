package formatter

type OutputFormat string

const (
	HTMLOutput     OutputFormat = "html"
	CSVOutput      OutputFormat = "csv"
	MarkdownOutput OutputFormat = "md"
	JSONOutput     OutputFormat = "json"
)

// IsValid checks whether requested output format is valid
func (of OutputFormat) IsValid() bool {
	switch of {
	case "markdown", "md", "html", "csv", "json":
		return true
	}
	return false
}

// FileOutputFormat returns appropriate file format, users can provide short
func (of OutputFormat) FileOutputFormat() OutputFormat {
	switch of {
	case "markdown", "md":
		return MarkdownOutput
	case "html":
		return HTMLOutput
	case "csv":
		return CSVOutput
	case "json":
		return JSONOutput
	}
	return HTMLOutput
}
