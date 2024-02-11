package formatter

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

// ExcelFormatter is struct defined for Excel Output use-case
type ExcelFormatter struct {
	config *Config
}

// Format the data to Excel and output it to an Excel file
func (f *ExcelFormatter) Format(td *TemplateData, templateContent string) (err error) {
	file := excelize.NewFile()
	sheetName := "Sheet1"

	// Create a style for center alignment
	style, err := file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	// Set the column headers
	file.SetCellValue(sheetName, "A1", "IP/Host")
	file.SetCellValue(sheetName, "B1", "Servicios")
	file.SetCellStyle(sheetName, "A1", "A1", style)
	file.SetCellStyle(sheetName, "B1", "B1", style)

	row := 2 // Start from row 2 for data

	for i := range td.NMAPRun.Host {
		var host *Host = &td.NMAPRun.Host[i]
		// Skipping hosts that are down
		if td.OutputOptions.ExcelOptions.SkipDownHosts && host.Status.State != "up" {
			continue
		}
		address := fmt.Sprintf("%s (%s)", host.JoinedAddresses("/"), host.JoinedHostNames("/"))

		// Set the IP/Host value
		cell := fmt.Sprintf("A%d", row)
		file.SetCellValue(sheetName, cell, address)
		file.SetCellStyle(sheetName, cell, cell, style)

		startRow := row // Remember the start row for this host

		for j := range host.Port {
			var port *Port = &host.Port[j]
			col := 'B' // Start from column B for Services

			// Set the Service value
			cell = fmt.Sprintf("%c%d", col, row)
			file.SetCellValue(sheetName, cell, fmt.Sprintf("%d/%s %s", port.PortID, port.Protocol, port.Service.Name))
			file.SetCellStyle(sheetName, cell, cell, style)
			row++ // Increment row for next service
		}

		// Merge cells in the IP/Host column for this host
		if row > startRow+1 { // Only merge if there's more than one service
			file.MergeCell(sheetName, fmt.Sprintf("A%d", startRow), fmt.Sprintf("A%d", row-1))
		}
	}

	// Save the Excel file
	err = file.SaveAs("nmap-output.xlsx")
	return err
}

func (f *ExcelFormatter) defaultTemplateContent() string {
	return HTMLSimpleTemplate
}
