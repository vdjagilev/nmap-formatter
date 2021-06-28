package types

type Ports struct {
	Port []Port `xml:"port"`
}

type Port struct {
	Protocol string      `xml:"protocol,attr"`
	PortID   string      `xml:"portid,attr"`
	State    PortState   `xml:"state"`
	Service  PortService `xml:"service"`
	Script   []Script    `xml:"script"`
}

type PortState struct {
	State     string `xml:"state,attr"`
	Reason    string `xml:"reason,attr"`
	ReasonTTL string `xml:"reason_ttl,attr"`
}

type PortService struct {
	Name      string   `xml:"name,attr"`
	Product   string   `xml:"product,attr"`
	Version   string   `xml:"version,attr"`
	ExtraInfo string   `xml:"extrainfo,attr"`
	Method    string   `xml:"method,attr"`
	Conf      string   `xml:"conf,attr"`
	CPE       []string `xml:"cpe"`
}

type ExtraPorts struct {
	State string `xml:"state,attr"`
	Count string `xml:"count,attr"`
}

type Script struct {
	ID     string `xml:"id,attr"`
	Output string `xml:"output,attr"`
}
