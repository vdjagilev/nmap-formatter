package formatter

import (
	// Used to embed graphviz template file
	_ "embed"
	"fmt"
	"strings"
	"text/template"
)

// DotFormatter is used to create Graphviz (dot) format
type DotFormatter struct {
	config *Config
}

// DotTemplate variable is used to store contents of graphviz template
//
//go:embed resources/templates/graphviz.tmpl
var DotTemplate string

const (
	// DotOpenPortColor defines default color of the opened port
	DotOpenPortColor = "#228B22"
	// DotFilteredPortColor defines default color of the filtered port
	DotFilteredPortColor = "#FFAE00"
	// DotClosedPortColor defies default color of the closed port
	DotClosedPortColor = "#DC143C"
	// DotDefaultColor defines default color of various elements (lines, boxes)
	DotDefaultColor = "gray"

	// DotFontStyle default font style
	DotFontStyle = "monospace"

	// DotLayout is a type of layout used in Graphviz (dot by default is the most fitting)
	DotLayout = "dot"
)

// DotDefaultOptions is a config map that is used in Graphviz template
var DotDefaultOptions = map[string]string{
	"default_font":  DotFontStyle,
	"layout":        DotLayout,
	"color_default": DotDefaultColor,
}

// DotTemplateData is a custom TemplateData struct that is used by DotFormatter
type DotTemplateData struct {
	NMAPRun   *NMAPRun
	Constants map[string]string
}

// Format the data and output it to appropriate io.Writer
func (f *DotFormatter) Format(td *TemplateData, templateContent string) (err error) {
	tmpl := template.New("dot")
	f.defineTemplateFunctions(tmpl)
	tmpl, err = tmpl.Parse(templateContent)
	if err != nil {
		return
	}
	dotTemplateData := DotTemplateData{
		NMAPRun:   &td.NMAPRun,
		Constants: DotDefaultOptions,
	}
	return tmpl.Execute(f.config.Writer, dotTemplateData)
}

// defaultTemplateContent returns default template content for any typical chosen formatter (HTML or Markdown)
func (f *DotFormatter) defaultTemplateContent() string {
	return DotTemplate
}

// defineTemplateFunctions defines all template functions that are used in dot templates
func (f *DotFormatter) defineTemplateFunctions(tmpl *template.Template) {
	tmpl.Funcs(
		template.FuncMap{
			"clean_ip":         cleanIP,
			"port_state_color": portStateColor,
			"hop_list":         hopList,
		},
	)
}

// cleanIP removes dots from IP address to make it possible to use in graphviz as an ID
func cleanIP(ip string) string {
	return strings.ReplaceAll(ip, ".", "")
}

// portStateColor returns hexademical color value for state port
func portStateColor(port *Port) string {
	switch port.State.State {
	case "open":
		return DotOpenPortColor
	case "filtered":
		return DotFilteredPortColor
	case "closed":
		return DotClosedPortColor
	}
	return DotDefaultColor
}

// hopList function returns a map with a list of hops where very first hop is `startHop` (scanner itself)
func hopList(hops []Hop, startHop string, endHopName string, endHopKey int) map[string]string {
	var hopList map[string]string = map[string]string{}
	var previous *Hop = nil
	for i := range hops {
		// Skip last hop, because it has the same IP as the target server
		if i == len(hops)-1 {
			break
		}
		if i == 0 {
			hopList[startHop] = fmt.Sprintf("hop%s", hops[i].IPAddr)
		} else {
			hopList[fmt.Sprintf("hop%s", previous.IPAddr)] = fmt.Sprintf("hop%s", hops[i].IPAddr)
		}
		previous = &hops[i]
	}
	if previous != nil {
		hopList[fmt.Sprintf("hop%s", previous.IPAddr)] = fmt.Sprintf("%s%d", endHopName, endHopKey)
	} else {
		hopList[startHop] = fmt.Sprintf("%s%d", endHopName, endHopKey)
	}
	return hopList
}
