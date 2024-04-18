package formatter

// OutputFormat is a resulting type of file that is converted/formatted from XML to HTML (for example)
type OutputFormat string

const (
	// HTMLOutput constant defines OutputFormat is HyperText Markup Language which can be viewed using browsers
	HTMLOutput OutputFormat = "html"
	// CSVOutput constant defines OutputFormat for Comma-Separated Values CSV file which is viewed most of the time in Excel
	CSVOutput OutputFormat = "csv"
	// MarkdownOutput constant defines OutputFormat for Markdown, which is handy and easy format to read-write
	MarkdownOutput OutputFormat = "md"
	// JSONOutput constant defines OutputFormat for JavaScript Object Notation, which is more useful for machine-related operations (parsing)
	JSONOutput OutputFormat = "json"
	// DotOutput constant defined OutputFormat for Dot (Graphviz), which can be used to generate various graphs
	DotOutput OutputFormat = "dot"
	// SqliteOutput constant defines OutputFormat for sqlite file, which can be used to generate sqlite embedded databases
	SqliteOutput OutputFormat = "sqlite"
	// ExcelOutput constant defines OutputFormat for Excel file, which can be used to generate Excel files
	ExcelOutput OutputFormat = "excel"
	// D2LangOutput constant defines OutputFormat for D2 language, which can be used to generate D2 language files
	D2LangOutput OutputFormat = "d2"
)

// IsValid checks whether requested output format is valid
func (of OutputFormat) IsValid() bool {
	// markdown & md is essentially the same thing
	switch of {
	case "markdown", "md", "html", "csv", "json", "dot", "sqlite", "excel", "d2":
		return true
	}
	return false
}
