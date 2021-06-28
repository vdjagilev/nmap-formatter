package types

type Host struct {
	StartTime     string        `xml:"starttime,attr"`
	EndTime       string        `xml:"endtime,attr"`
	Ports         Ports         `xml:"ports"`
	HostAddress   HostAddress   `xml:"address"`
	HostNames     HostNames     `xml:"hostnames"`
	Status        HostStatus    `xml:"status"`
	OS            OS            `xml:"os"`
	Trace         Trace         `xml:"trace"`
	Uptime        Uptime        `xml:"uptime"`
	Distance      Distance      `xml:"distance"`
	TCPSequence   TCPSequence   `xml:"tcpsequence"`
	IPIDSequence  IPIDSequence  `xml:"ipidsequence"`
	TCPTSSequence TCPTSSequence `xml:"tcptssequence"`
}

type TCPTSSequence struct {
	Class  string `xml:"class,attr"`
	Values string `xml:"values,attr"`
}

type IPIDSequence struct {
	Class  string `xml:"class,attr"`
	Values string `xml:"values,attr"`
}

type TCPSequence struct {
	Index      string `xml:"index,attr"`
	Difficulty string `xml:"difficulty,attr"`
	Values     string `xml:"values,attr"`
}

type Uptime struct {
	Seconds  string `xml:"seconds,attr"`
	LastBoot string `xml:"lastboot,attr"`
}

type Distance struct {
	Value string `xml:"value,attr"`
}

type HostStatus struct {
	State  string `xml:"state,attr"`
	Reason string `xml:"reason,attr"`
}

type HostAddress struct {
	Address     string `xml:"addr,attr"`
	AddressType string `xml:"addrtype,attr"`
}

type HostNames struct {
	HostName []HostName `xml:"hostname"`
}

type HostName struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type OS struct {
	OSPortUsed []OSPortUsed `xml:"portused"`
	OSClass    OSClass      `xml:"osclass"`
	OSMatch    OSMatch      `xml:"osmatch"`
}

type OSPortUsed struct {
	State    string `xml:"state,attr"`
	Protocol string `xml:"proto,attr"`
	PortID   string `xml:"portid,attr"`
}

type OSClass struct {
	Type     string   `xml:"type,attr"`
	Vendor   string   `xml:"vendor,attr"`
	OSFamily string   `xml:"osfamily,attr"`
	OSGen    string   `xml:"osgen,attr"`
	Accuracy string   `xml:"accuracy,attr"`
	CPE      []string `xml:"cpe"`
}

type OSMatch struct {
	Name     string `xml:"name,attr"`
	Accuracy string `xml:"accuracy,attr"`
	Line     string `xml:"line,attr"`
}

type Trace struct {
	Port     string `xml:"port,attr"`
	Protocol string `xml:"proto,attr"`
	Hops     []Hop  `xml:"hop"`
}

type Hop struct {
	TTL    string `xml:"ttl,attr"`
	IPAddr string `xml:"ipaddr,attr"`
	RTT    string `xml:"rtt,attr"`
	Host   string `xml:"host,attr"`
}
