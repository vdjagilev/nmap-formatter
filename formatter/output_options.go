package formatter

// OutputOptions describes various output options for nmap formatter
type OutputOptions struct {
	// The hosts that are down wont be displayed in the TOC
	SkipDownHosts bool
}
