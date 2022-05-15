package formatter

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

// Workflow interface that describes the main functions that are used in nmap-formatter
type Workflow interface {
	Execute() (err error)
	SetConfig(c *Config)
	SetInputFile()
	SetOutputFile()
}

// MainWorkflow is main workflow implementation struct
type MainWorkflow struct {
	Config *Config
}

// SetConfig is a simple setter-function that sets the configuration
func (w *MainWorkflow) SetConfig(c *Config) {
	w.Config = c
}

// SetOutputFile sets output file (file descriptor) depending on config and returns error
// if there is output file reading issue
func (w *MainWorkflow) SetOutputFile() {
	if w.Config.OutputFile == "" {
		w.Config.Writer = os.Stdout
	} else {
		// Error has been checked before executing this function
		f, _ := os.OpenFile(string(w.Config.OutputFile), os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.ModePerm)
		w.Config.Writer = f
	}
}

// SetInputFile sets an input file (file descriptor) in the config
func (w *MainWorkflow) SetInputFile() {
	var inputFile *os.File
	if w.Config.InputFileConfig.IsStdin {
		inputFile = os.Stdin
	} else {
		// Error has been checked before executing this function
		inputFile, _ = os.Open(w.Config.InputFileConfig.Path)
	}
	w.Config.InputFileConfig.Source = inputFile
}

// Execute is the core of the application which executes required steps
// one-by-one to achieve formatting from input -> output.
func (w *MainWorkflow) Execute() (err error) {
	// Reading & parsing the input file
	NMAPRun, err := w.parse()
	if err != nil {
		return
	}

	// Build template data with NMAPRun entry & various output options
	templateData := TemplateData{
		NMAPRun:       NMAPRun,
		OutputOptions: w.Config.OutputOptions,
	}

	// Setting custom options for template if they exist
	if len(w.Config.CustomOptions) > 0 {
		templateData.CustomOptions = w.Config.CustomOptionsMap()
	}

	// Getting new instance of formatter based on provided config
	formatter := New(w.Config)

	// This part usually should not happen
	if formatter == nil {
		return fmt.Errorf("no formatter is defined")
	}

	// Trying to read template content (read a file, or get default in case where no option was used)
	templateContent, err := TemplateContent(formatter, w.Config)
	if err != nil {
		return fmt.Errorf("error getting template content: %v", err)
	}
	return formatter.Format(&templateData, templateContent)
}

// parse reads & unmarshalles the input file into NMAPRun struct
func (w *MainWorkflow) parse() (run NMAPRun, err error) {
	input, err := w.Config.InputFileConfig.ReadContents()
	if err != nil {
		return
	}
	if err = xml.Unmarshal(input, &run); err != nil {
		return
	}

	// A temporary solution to parse `<address>` node separately and choose only ipv4 or ipv6 for
	// host addresses (avoid using mac-addresses), issue: #105
	if len(run.Host) > 0 {
		overrideHostAddresses(&run, bytes.NewReader(input))
	}
	return run, nil
}

// overrideHostAddresses fixes output by overriding HostAddress.Address
// struct for the reason if IPv4/IPv6 `<address>` type is getting overwritten
// by MAC-type address. This fix is temporary in order to provide a bugfix, but to avoid any BC-breaks
// Will be removed in new major release
func overrideHostAddresses(run *NMAPRun, reader io.Reader) {
	var tag string
	var hostID int = 0
	var unmarshalledHost *Host = &run.Host[hostID]
	var decoder *xml.Decoder = xml.NewDecoder(reader)
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch element := token.(type) {
		case xml.StartElement:
			tag = element.Name.Local
			if tag == "host" {
				unmarshalledHost = &run.Host[hostID]
				hostID++
			} else if tag == "address" && unmarshalledHost != nil {
				setHostAddress(unmarshalledHost, &element)
			}
		default:
		}
	}
}

// setHostAddress setting the Address & AddressType
func setHostAddress(h *Host, e *xml.StartElement) {
	if hasAddressIPAttr(e) {
		for _, attr := range e.Attr {
			switch attr.Name.Local {
			case "addr":
				h.HostAddress.Address = attr.Value
			case "addrtype":
				h.HostAddress.AddressType = attr.Value
			}
		}
	}
}

// hasAddressIPAttr determines whether XML element attributes have IPv4 or IPv6 address type
func hasAddressIPAttr(e *xml.StartElement) bool {
	for _, attr := range e.Attr {
		if attr.Name.Local == "addrtype" &&
			(attr.Value == "ipv4" || attr.Value == "ipv6") {
			return true
		}
	}
	return false
}
