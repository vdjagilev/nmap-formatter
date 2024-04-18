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

type D2LangFormatter struct {
	config *Config
}

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

		hostId := hex.EncodeToString(fnv.Sum(nil))
		graph, _, _ = d2oracle.Create(graph, nil, hostId)
		graph, _ = d2oracle.Set(graph, nil, hostId+".label", nil, &hostLabel)
		graph, _ = d2oracle.Set(graph, nil, "nmap -> "+hostId, nil, nil)

		for j := range host.Port {
			port := &host.Port[j]
			portId := fmt.Sprintf("%s-port%d", hostId, port.PortID)
			graph, _, _ = d2oracle.Create(graph, nil, portId)
			portLabel := fmt.Sprintf("%d/%s\n%s\n%s", port.PortID, port.Protocol, port.State.State, port.Service.Name)
			graph, _ = d2oracle.Set(graph, nil, portId+".label", nil, &portLabel)
			shape := "circle"
			graph, _ = d2oracle.Set(graph, nil, portId+".shape", nil, &shape)
			width := "25"
			graph, _ = d2oracle.Set(graph, nil, portId+".width", nil, &width)
			graph, _ = d2oracle.Move(graph, nil, portId, hostId+"."+portId, true)
		}
	}
	_, err = f.config.Writer.Write([]byte(d2format.Format(graph.AST)))
	return
}

// defaultTemplateContent does not return anything in this case
func (f *D2LangFormatter) defaultTemplateContent() string {
	return ""
}
