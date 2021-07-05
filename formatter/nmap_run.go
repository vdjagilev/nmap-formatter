package formatter

type NMAPRun struct {
	Scanner   string    `xml:"scanner,attr"`
	Args      string    `xml:"args,attr"`
	Start     string    `xml:"start,attr"`
	StartStr  string    `xml:"startstr,attr"`
	Version   string    `xml:"version,attr"`
	ScanInfo  ScanInfo  `xml:"scaninfo"`
	Host      []Host    `xml:"host"`
	Verbose   Verbose   `xml:"verbose"`
	Debugging Debugging `xml:"debugging"`
	RunStats  RunStats  `xml:"runstats"`
}

type ScanInfo struct {
	Type        string `xml:"type,attr"`
	Protocol    string `xml:"protocol,attr"`
	NumServices string `xml:"numservices,attr"`
	Services    string `xml:"services,attr"`
}

type Verbose struct {
	Level string `xml:"level,attr"`
}

type Debugging struct {
	Level string `xml:"level,attr"`
}

type RunStats struct {
	Finished Finished  `xml:"finished"`
	Hosts    StatHosts `xml:"hosts"`
}

type Finished struct {
	Time    string `xml:"time,attr"`
	TimeStr string `xml:"timestr,attr"`
	Elapsed string `xml:"elapsed,attr"`
	Summary string `xml:"summary,attr"`
	Exit    string `xml:"exit,attr"`
}

type StatHosts struct {
	Up    string `xml:"up,attr"`
	Down  string `xml:"down,attr"`
	Total string `xml:"total,attr"`
}

// TODO: Other info like: scaninfo, verbose, runstats, et cetera
