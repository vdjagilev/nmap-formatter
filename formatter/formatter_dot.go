package formatter

import (
	_ "embed"
	"fmt"
	"strings"
	"text/template"
)

type DotFormatter struct {
	config *Config
}

//go:embed resources/templates/graphviz.tmpl
// DotTemplate variable is used to store contents of graphviz template
var DotTemplate string

const (
	DotOpenPortColor     = "#228B22"
	DotFilteredPortColor = "#FFAE00"
	DotClosedPortColor   = "#DC143C"
	DotDefaultColor      = "gray"

	DotFontStyle = "monospace"

	DotLayout = "dot"
)

var DotDefaultOptions = map[string]string{
	"default_font":  DotFontStyle,
	"layout":        DotLayout,
	"color_default": DotDefaultColor,
}

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

func hopList(hops []Hop, startHop string, endHopName string, endHopKey int) map[string]string {
	var hopList map[string]string = map[string]string{}
	var previous *Hop = nil
	for i := range hops {
		if i == 0 {
			hopList[startHop] = fmt.Sprintf("hop%s", hops[i].IPAddr)
		} else {
			hopList[fmt.Sprintf("hop%s", previous.IPAddr)] = fmt.Sprintf("hop%s", hops[i].IPAddr)
		}
		previous = &hops[i]
	}
	hopList[fmt.Sprintf("hop%s", previous.IPAddr)] = fmt.Sprintf("%s%d", endHopName, endHopKey)
	return hopList
}
