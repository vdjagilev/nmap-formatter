package formatter

// OutputOptions describes various output options for nmap formatter
type OutputOptions struct {
	// The hosts that are down wont be displayed in the TOC
	SkipDownHosts bool
	// JSONPrettyPrint defines if JSON output would be pretty-printed (human-readable) or not (machine readable)
	JSONPrettyPrint bool
	// SkipSummary skips general summary for HTML & Markdown
	SkipSummary bool
	// SkipTraceroute skips traceroute information for HTML & Markdown
	SkipTraceroute bool
	// SkipMetrics skips metrics related data for HTML
	SkipMetrics bool
	// SkipPortScripts skips port scripts information for HTML
	SkipPortScripts bool
}
