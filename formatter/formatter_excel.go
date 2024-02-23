package formatter

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

var EXCEL_COL_WIDTH float64 = 50

// ExcelFormatter is struct defined for Excel Output use-case
type ExcelFormatter struct {
	config *Config
}

type CellData struct {
	sheetName string
	style     int
	file      *excelize.File
}

func (cd *CellData) writeCell(cell string, value string) error {
	err := cd.file.SetCellValue(cd.sheetName, cell, value)
	if err != nil {
		return err
	}

	return cd.file.SetCellStyle(cd.sheetName, cell, cell, cd.style)
}

// Format the data to Excel and output it to an Excel file
func (f *ExcelFormatter) Format(td *TemplateData, templateContent string) (err error) {
	file := excelize.NewFile()
	sheetName := "Sheet1"

	// Create a style for center alignment
	style, err := file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
	})
	if err != nil {
		return err
	}

	// Reusable cell data struct for writing cells
	cd := &CellData{
		sheetName: sheetName,
		style:     style,
		file:      file,
	}

	// Set column headers with titles
	err = f.writeHeaders(cd)
	if err != nil {
		return err
	}

	err = f.writeHostRows(td.NMAPRun.Host, cd)
	if err != nil {
		return err
	}

	return file.Write(f.config.Writer, excelize.Options{})
}

func (f *ExcelFormatter) writeHostRows(h []Host, cd *CellData) error {
	var err error
	row := 2 // Start from row 2 for data

	for i := range h {
		host := h[i]

		// Skipping hosts that are down
		if f.config.OutputOptions.ExcelOptions.SkipDownHosts && host.Status.State != "up" {
			continue
		}

		joinedAddresses := host.JoinedAddresses("/")
		joinedHostnames := host.JoinedHostNames("/")
		address := ""

		if joinedHostnames == "" {
			address = joinedAddresses
		} else {
			address = fmt.Sprintf(
				"%s (%s)",
				joinedAddresses,
				joinedHostnames,
			)
		}

		// Set the IP/Host value
		cell := fmt.Sprintf("A%d", row)
		err = cd.writeCell(cell, address)
		if err != nil {
			return err
		}

		startRow := row // Remember the start row for this host
		err = f.writePorts(host.Port, cd, &row)
		if err != nil {
			return err
		}

		// Merge cells in the IP/Host column for this host
		if row > startRow+1 {
			err = cd.file.MergeCell(
				cd.sheetName,
				fmt.Sprintf("A%d", startRow),
				fmt.Sprintf("A%d", row-1),
			)

			if err != nil {
				return err
			}
		}
	}
	return err
}

func (f *ExcelFormatter) writePorts(p []Port, cd *CellData, row *int) error {
	var err error

	for i := range p {
		port := p[i]

		// Set the Service value for column B for services
		err = cd.writeCell(
			fmt.Sprintf("%c%d", 'B', *row),
			fmt.Sprintf("%d/%s %s", port.PortID, port.Protocol, port.Service.Name),
		)
		if err != nil {
			return err
		}
		*row++
	}
	return err
}

func (f *ExcelFormatter) writeHeaders(cd *CellData) error {
	err := cd.writeCell("A1", "IP/Host")
	if err != nil {
		return err
	}

	// Setting the width of the columns in order not to cut the text
	err = cd.file.SetColWidth(cd.sheetName, "A", "B", EXCEL_COL_WIDTH)
	if err != nil {
		return err
	}

	return cd.writeCell("B1", "Services")
}

func (f *ExcelFormatter) defaultTemplateContent() string {
	return ""
}
