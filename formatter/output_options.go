package formatter

// OutputOptions describes various output options for nmap formatter
type OutputOptions struct {
	// The hosts that are down wont be displayed in the TOC
	SkipDownHosts bool
	// JSONPrettyPrint defines if JSON output would be pretty-printed (human-readable) or not (machine readable)
	JSONPrettyPrint bool
	// SkipSummary skips general summary for HTML & Markdown formats
	SkipSummary bool
	// SkipTraceroute skips traceroute information for HTML & Markdown formats
	SkipTraceroute bool
	// SkipMetrics skips metrics related data for HTML & Markdown formats
	SkipMetrics bool
	// SkipPortScripts skips port scripts information for HTML & Markdown formats
	SkipPortScripts bool
}
