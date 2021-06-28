package formatter

import (
	"encoding/csv"

	"github.com/vdjagilev/nmap-formatter/types"
)

type CSVFormatter struct {
	Config *Config
}

// Format the data to CSV and output it to appropriate io.Writer
func (f *CSVFormatter) Format(td *types.TemplateData) (err error) {
	return csv.NewWriter(f.Config.Writer).WriteAll(f.convert(td))
}

// convert uses NMAPRun struct to convert all data to [][]string type
func (f *CSVFormatter) convert(td *types.TemplateData) (data [][]string) {
	data = append(data, []string{"IP", "Port", "Protocol", "State", "Service", "Reason", "Product", "Version", "Extra info"})
	for _, host := range td.NMAPRun.Host {
		for j, port := range host.Ports.Port {
			// Show host IP only once, to avoid repetitions
			hostIP := " "
			if j == 0 {
				hostIP = host.HostAddress.Address
			}
			data = append(
				data,
				[]string{
					hostIP,
					port.PortID,
					port.Protocol,
					port.State.State,
					port.Service.Name,
					port.State.Reason,
					port.Service.Product,
					port.Service.Version,
					port.Service.ExtraInfo,
				},
			)
		}
	}
	return
}
