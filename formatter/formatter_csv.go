package formatter

import (
	"encoding/csv"
	"fmt"
)

// CSVFormatter is struct defined for CSV Output use-case
type CSVFormatter struct {
	Config *Config
}

// Format the data to CSV and output it to appropriate io.Writer
func (f *CSVFormatter) Format(td *TemplateData) (err error) {
	return csv.NewWriter(f.Config.Writer).WriteAll(f.convert(td))
}

// convert uses NMAPRun struct to convert all data to [][]string type
func (f *CSVFormatter) convert(td *TemplateData) (data [][]string) {
	data = append(data, []string{"IP", "Port", "Protocol", "State", "Service", "Reason", "Product", "Version", "Extra info"})
	for _, host := range td.NMAPRun.Host {
		// Skipping hosts that are down
		if td.OutputOptions.SkipDownHosts && host.Status.State != "up" {
			continue
		}
		address := fmt.Sprintf("%s (%s)", host.HostAddress.Address, host.Status.State)
		data = append(data, []string{address, "", "", "", "", "", "", "", ""})
		for _, port := range host.Ports.Port {
			data = append(
				data,
				[]string{
					"",
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
