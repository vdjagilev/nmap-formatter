package formatter

// OutputOptions describes various output options for nmap formatter
type OutputOptions struct {
	// The hosts that are down wont be displayed in the TOC
	SkipDownHosts bool
	// JSONPrettyPrint defines if JSON output would be pretty-printed (human-readable) or not (machine readable)
	JSONPrettyPrint bool
}
