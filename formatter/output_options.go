package formatter

// OutputOptions describes various output options for nmap formatter
type OutputOptions struct {
	HTMLOptions         HTMLOutputOptions
	MarkdownOptions     MarkdownOutputOptions
	JSONOptions         JSONOutputOptions
	CSVOptions          CSVOutputOptions
	SqliteOutputOptions SqliteOutputOptions
	ExcelOptions        ExcelOutputOptions
	D2LangOptions       D2LangOutputOptions
}

// HTMLOutputOptions stores options related only to HTML conversion/formatting
type HTMLOutputOptions struct {
	// SkipHeader skips the header of the HTML output
	SkipHeader bool
	// SkipTOC skips the table of contents in the HTML output
	SkipTOC bool
	// SkipDownHosts skips hosts that are down (including TOC)
	SkipDownHosts bool
	// SkipSummary skips general summary for HTML
	SkipSummary bool
	// SkipTraceroute skips traceroute information for HTML
	SkipTraceroute bool
	// SkipMetrics skips metrics related data for HTML
	SkipMetrics bool
	// SkipPortScripts skips port scripts information for HTML
	SkipPortScripts bool
	// DarkMode sets a style to be mostly in dark colours, if false, light colours would be used
	DarkMode bool
	// FloatingContentsTable is an option to make contents table float on the side of the page
	FloatingContentsTable bool
}

// MarkdownOutputOptions stores options related only to Markdown conversion/formatting
type MarkdownOutputOptions struct {
	// SkipHeader skips the header
	SkipHeader bool
	// SkipTOC skips the table of contents
	SkipTOC bool
	// SkipDownHosts skips hosts that are down (including TOC)
	SkipDownHosts bool
	// SkipSummary skips general summary for Markdown
	SkipSummary bool
	// SkipPortScripts skips port scripts information for Markdown
	SkipPortScripts bool
	// SkipTraceroute skips traceroute information for Markdown
	SkipTraceroute bool
	// SkipMetrics skips metrics related data for Markdown
	SkipMetrics bool
}

// JSONOutputOptions store option related only to JSON conversion/formatting
type JSONOutputOptions struct {
	// PrettyPrint defines if JSON output would be pretty-printed (human-readable) or not (machine readable)
	PrettyPrint bool
}

// CSVOutputOptions store option related only to CSV conversion/formatting
type CSVOutputOptions struct {
	// The hosts that are down won't be displayed
	SkipDownHosts bool
}

// SqliteOutputOptions store options related to SQLite database formatting
type SqliteOutputOptions struct {
	// DSN is a Data Source Name to sqlite embedded database, by default it's empty which results in stdout output,
	// however, if both DSN and OutputFile provided, OutputFile option takes preference over SqliteOutputOptions.DSN
	DSN string
	// ScanIdentifier is a unique string passed by the user to identify unique scans. If it's empty, it's generated automatically
	ScanIdentifier string
}

// ExcelOutputOptions store options related to Excel file formatting
type ExcelOutputOptions struct {
	// The hosts that are down won't be displayed
	SkipDownHosts bool
}

// D2LangOutputOptions store options related to D2 language file formatting
type D2LangOutputOptions struct {
	// The hosts that are down won't be displayed
	SkipDownHosts bool
}
