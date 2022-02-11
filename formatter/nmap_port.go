package formatter

// Ports has a list of Port structs. Ports itself is defined in Host struct
type Ports struct {
	Port []Port `xml:"port"`
}

// Port record contains main information about port that was scanned
type Port struct {
	Protocol string      `xml:"protocol,attr"`
	PortID   int         `xml:"portid,attr"`
	State    PortState   `xml:"state"`
	Service  PortService `xml:"service"`
	Script   []Script    `xml:"script"`
}

// PortState describes information about the port state and why it's state was defined that way
type PortState struct {
	State     string `xml:"state,attr"`
	Reason    string `xml:"reason,attr"`
	ReasonTTL string `xml:"reason_ttl,attr"`
}

// PortService struct contains information about the service that is located on certain port
type PortService struct {
	Name      string   `xml:"name,attr"`
	Product   string   `xml:"product,attr"`
	Version   string   `xml:"version,attr"`
	ExtraInfo string   `xml:"extrainfo,attr"`
	Method    string   `xml:"method,attr"`
	Conf      string   `xml:"conf,attr"`
	CPE       []string `xml:"cpe"`
}

// ExtraPorts contains information about certain amount of ports that were (for example) filtered
type ExtraPorts struct {
	State string `xml:"state,attr"`
	Count int    `xml:"count,attr"`
}

// Script defines a script ID and script output (result)
type Script struct {
	ID     string `xml:"id,attr"`
	Output string `xml:"output,attr"`
}
