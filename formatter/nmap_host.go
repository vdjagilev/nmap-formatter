package formatter

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

// Host describes host related entry (`host` node)
type Host struct {
	StartTime     int           `xml:"starttime,attr"`
	EndTime       int           `xml:"endtime,attr"`
	Port          []Port        `xml:"ports>port"`
	HostAddress   []HostAddress `xml:"address"`
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

// JoinedAddresses joins all possible host addresses with a delimiter string
func (h *Host) JoinedAddresses(delimiter string) string {
	var addr string = ""
	for i := range h.HostAddress {
		// First element does not require prepended delimiter
		if i == 0 {
			addr += h.HostAddress[i].Address
		} else {
			addr += fmt.Sprintf(" %s %s", delimiter, h.HostAddress[i].Address)
		}
	}
	return addr
}

// TCPTSSequence describes all information related to `<tcptssequence>` node
type TCPTSSequence struct {
	Class  string `xml:"class,attr"`
	Values string `xml:"values,attr"`
}

// IPIDSequence describes all information related to `<ipidsequence>` node
type IPIDSequence struct {
	Class  string `xml:"class,attr"`
	Values string `xml:"values,attr"`
}

// TCPSequence describes all information related to `<tcpsequence>`
type TCPSequence struct {
	Index      string `xml:"index,attr"`
	Difficulty string `xml:"difficulty,attr"`
	Values     string `xml:"values,attr"`
}

// Uptime shows the information about host uptime
type Uptime struct {
	Seconds  int    `xml:"seconds,attr"`
	LastBoot string `xml:"lastboot,attr"`
}

// Distance describes amount of hops to the target
type Distance struct {
	Value int `xml:"value,attr"`
}

// HostStatus describes the state (up or down) of the host and the reason
type HostStatus struct {
	State  string `xml:"state,attr"`
	Reason string `xml:"reason,attr"`
}

// HostAddress struct contains the host address (IP) and type of it.
type HostAddress struct {
	Address     string `xml:"addr,attr"`
	AddressType string `xml:"addrtype,attr"`
}

// HostNames struct contains list of hostnames (domains) that this host has
type HostNames struct {
	HostName []HostName `xml:"hostname"`
}

// HostName defines the name of the host and type of DNS record (like PTR for example)
type HostName struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

// OS describes all information about underlying operating system that this host operates
type OS struct {
	OSPortUsed []OSPortUsed `xml:"portused"`
	OSClass    []OSClass    `xml:"osclass"`
	OSMatch    []OSMatch    `xml:"osmatch"`
}

// OSPortUsed defines which ports were used for OS detection
type OSPortUsed struct {
	State    string `xml:"state,attr"`
	Protocol string `xml:"proto,attr"`
	PortID   int    `xml:"portid,attr"`
}

// OSClass contains all information about operating system family
type OSClass struct {
	Type     string   `xml:"type,attr"`
	Vendor   string   `xml:"vendor,attr"`
	OSFamily string   `xml:"osfamily,attr"`
	OSGen    string   `xml:"osgen,attr"`
	Accuracy string   `xml:"accuracy,attr"`
	CPE      []string `xml:"cpe"`
}

// OSMatch is a record of OS that matched with certain accuracy
type OSMatch struct {
	Name     string `xml:"name,attr"`
	Accuracy string `xml:"accuracy,attr"`
	Line     string `xml:"line,attr"`
}

// Trace struct contains trace information with hops
type Trace struct {
	Port     int    `xml:"port,attr"`
	Protocol string `xml:"proto,attr"`
	Hops     []Hop  `xml:"hop"`
}

// Hop struct contains information about HOP record with time to live, host name, IP
type Hop struct {
	TTL    int    `xml:"ttl,attr"`
	IPAddr string `xml:"ipaddr,attr"`
	RTT    RTT    `xml:"rtt,attr"`
	Host   string `xml:"host,attr"`
}

// RTT is a separate type that is located in Hop struct
type RTT float64

// UnmarshalXMLAttr is a separate function that attempts to parse RTT float value
// if it fails to do so, it sets the value to 0.0
func (r *RTT) UnmarshalXMLAttr(attr xml.Attr) error {
	value, err := strconv.ParseFloat(attr.Value, 64)
	if err != nil {
		value = 0.0
	}
	*(*float64)(r) = value
	return nil
}
