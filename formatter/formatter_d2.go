package formatter

import (
	"context"
	"encoding/hex"
	"fmt"
	"hash/fnv"

	"oss.terrastruct.com/d2/d2format"
	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2oracle"
	"oss.terrastruct.com/d2/lib/log"
	"oss.terrastruct.com/d2/lib/textmeasure"
)

// D2LangFormatter is struct defined for D2 Language Output use-case
type D2LangFormatter struct {
	config *Config
}

// Format the data to D2 Language and output it to a D2 Language file
func (f *D2LangFormatter) Format(td *TemplateData, templateContent string) (err error) {
	ruler, _ := textmeasure.NewRuler()
	layoutResolver := func(engine string) (d2graph.LayoutGraph, error) {
		return d2dagrelayout.DefaultLayout, nil
	}
	compileOpts := &d2lib.CompileOptions{
		LayoutResolver: layoutResolver,
		Ruler:          ruler,
	}

	_, graph, _ := d2lib.Compile(log.Stderr(context.Background()), "nmap", compileOpts, nil)

	for i := range td.NMAPRun.Host {
		host := &td.NMAPRun.Host[i]
		fnv := fnv.New128()

		if host.ShouldSkipHost(td.OutputOptions.D2LangOptions.SkipDownHosts) {
			continue
		}

		address := host.JoinedAddresses("/")
		hostnames := host.JoinedHostNames("/")
		hostLabel := address
		if hostnames != "" {
			hostLabel = fmt.Sprintf("%s\n(%s)", address, hostnames)
		}
		_, err := fnv.Write([]byte(address))
		if err != nil {
			return err
		}

		hostID := hex.EncodeToString(fnv.Sum(nil))
		graph, _, _ = d2oracle.Create(graph, nil, hostID)
		graph, _ = d2oracle.Set(graph, nil, hostID+".label", nil, &hostLabel)
		graph, _ = d2oracle.Set(graph, nil, "nmap -> "+hostID, nil, nil)

		for j := range host.Port {
			port := &host.Port[j]
			portID := fmt.Sprintf("%s-port%d", hostID, port.PortID)
			graph, _, _ = d2oracle.Create(graph, nil, portID)
			portLabel := fmt.Sprintf("%d/%s\n%s\n%s", port.PortID, port.Protocol, port.State.State, port.Service.Name)
			graph, _ = d2oracle.Set(graph, nil, portID+".label", nil, &portLabel)
			shape := "circle"
			graph, _ = d2oracle.Set(graph, nil, portID+".shape", nil, &shape)
			width := "25"
			graph, _ = d2oracle.Set(graph, nil, portID+".width", nil, &width)
			graph, _ = d2oracle.Move(graph, nil, portID, hostID+"."+portID, true)
		}
	}
	_, err = f.config.Writer.Write([]byte(d2format.Format(graph.AST)))
	return
}

// defaultTemplateContent does not return anything in this case
func (f *D2LangFormatter) defaultTemplateContent() string {
	return ""
}
