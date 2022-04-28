package formatter

// OutputOptions describes various output options for nmap formatter
type OutputOptions struct {
	HTMLOptions     HTMLOutputOptions
	MarkdownOptions MarkdownOutputOptions
	JSONOptions     JSONOutputOptions
	CSVOptions      CSVOutputOptions
}

type HTMLOutputOptions struct {
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
}

type MarkdownOutputOptions struct {
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

type JSONOutputOptions struct {
	// PrettyPrint defines if JSON output would be pretty-printed (human-readable) or not (machine readable)
	PrettyPrint bool
}

type CSVOutputOptions struct {
	// The hosts that are down won't be displayed
	SkipDownHosts bool
}
