package formatter

// NMAPRun represents main `<nmaprun>` node which contains meta-information about the scan
// For example: scanner, what arguments used during scan, nmap version, verbosity level, et cetera
// Main information about scanned hosts is in the `host` node
type NMAPRun struct {
	Scanner   string    `xml:"scanner,attr"`
	Args      string    `xml:"args,attr"`
	Start     int       `xml:"start,attr"`
	StartStr  string    `xml:"startstr,attr"`
	Version   string    `xml:"version,attr"`
	ScanInfo  ScanInfo  `xml:"scaninfo"`
	Host      []Host    `xml:"host"`
	Verbose   Verbose   `xml:"verbose"`
	Debugging Debugging `xml:"debugging"`
	RunStats  RunStats  `xml:"runstats"`
}

// ScanInfo shows what type of scan it was and number of services covered
type ScanInfo struct {
	Type        string `xml:"type,attr"`
	Protocol    string `xml:"protocol,attr"`
	NumServices int    `xml:"numservices,attr"`
	Services    string `xml:"services,attr"`
}

// Verbose defines verbosity level that was configured during NMAP execution
type Verbose struct {
	Level int `xml:"level,attr"`
}

// Debugging defines level of debug during NMAP execution
type Debugging struct {
	Level int `xml:"level,attr"`
}

// RunStats contains other nodes that refer to statistics of the scan
type RunStats struct {
	Finished Finished  `xml:"finished"`
	Hosts    StatHosts `xml:"hosts"`
}

// Finished is part of `RunStats` struct, it has all information related to the time (started, how much time it took) and summary incl. exit status code
type Finished struct {
	Time    int     `xml:"time,attr"`
	TimeStr string  `xml:"timestr,attr"`
	Elapsed float64 `xml:"elapsed,attr"`
	Summary string  `xml:"summary,attr"`
	Exit    string  `xml:"exit,attr"`
}

// StatHosts contains statistics about hosts that are up or down
type StatHosts struct {
	Up    int `xml:"up,attr"`
	Down  int `xml:"down,attr"`
	Total int `xml:"total,attr"`
}

// AllHops is getting all possible hops that occurred during the scan and
// merges them uniquely into one map
func (n *NMAPRun) AllHops() map[string]Hop {
	hops := map[string]Hop{}
	for i := range n.Host {
		for j := range n.Host[i].Trace.Hops {
			// Skip the last hop, because it has the same IP as the target server
			if j == len(n.Host[i].Trace.Hops)-1 {
				break
			}
			hop := n.Host[i].Trace.Hops[j]
			hops[hop.IPAddr] = hop
		}
	}
	return hops
}
