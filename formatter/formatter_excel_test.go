package formatter

import (
	"bytes"
	"testing"

	"github.com/xuri/excelize/v2"
)

type excelMockedWriter struct {
	bytes.Buffer
}

func (w *excelMockedWriter) Close() error {
	return nil
}

func TestExcelFormatter_Format(t *testing.T) {
	type args struct {
		td *TemplateData
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantCells map[string]string
		wantMerge []string
	}{
		{
			name: "Empty Test ExcelFormatter.Format",
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Host: []Host{},
					},
				},
			},
			wantErr: false,
			wantCells: map[string]string{
				"A1": "IP/Host",
				"B1": "Services",
			},
		},
		{
			name: "Test ExcelFormatter.Format Basic 3 hosts",
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Host: []Host{
							{
								StartTime: 0,
								EndTime:   0,
								Port: []Port{
									{
										PortID:   80,
										Protocol: "tcp",
										Service: PortService{
											Name: "http",
										},
									},
									{
										PortID:   443,
										Protocol: "tcp",
										Service: PortService{
											Name: "https",
										},
									},
								},
								HostAddress: []HostAddress{
									{
										Address: "192.168.1.1",
									},
								},
								HostNames: HostNames{
									HostName: []HostName{
										{
											Name: "example.com",
											Type: "",
										},
										{
											Name: "www.example.com",
											Type: "",
										},
									},
								},
								Status: HostStatus{
									State: "up",
								},
								OS:            OS{},
								Trace:         Trace{},
								Uptime:        Uptime{},
								Distance:      Distance{},
								TCPSequence:   TCPSequence{},
								IPIDSequence:  IPIDSequence{},
								TCPTSSequence: TCPTSSequence{},
							},
							{
								StartTime: 0,
								EndTime:   0,
								Port:      []Port{},
								HostAddress: []HostAddress{
									{
										Address: "192.168.1.2",
									},
								},
								HostNames: HostNames{},
								Status: HostStatus{
									State: "down",
								},
								OS:            OS{},
								Trace:         Trace{},
								Uptime:        Uptime{},
								Distance:      Distance{},
								TCPSequence:   TCPSequence{},
								IPIDSequence:  IPIDSequence{},
								TCPTSSequence: TCPTSSequence{},
							},
							{
								StartTime: 0,
								EndTime:   0,
								Port:      []Port{},
								HostAddress: []HostAddress{
									{
										Address: "192.168.1.3",
									},
								},
								HostNames: HostNames{},
								Status: HostStatus{
									State: "up",
								},
								OS:            OS{},
								Trace:         Trace{},
								Uptime:        Uptime{},
								Distance:      Distance{},
								TCPSequence:   TCPSequence{},
								IPIDSequence:  IPIDSequence{},
								TCPTSSequence: TCPTSSequence{},
							},
						},
					},
				},
			},
			wantErr: false,
			wantCells: map[string]string{
				"A1": "IP/Host",
				"B1": "Services",
				"A2": "192.168.1.1 (example.com/www.example.com) [up]",
				"B2": "80/tcp http",
				"B3": "443/tcp https",
				"A4": "192.168.1.2 [down]",
				"B4": "-",
				"A5": "192.168.1.3 [up]",
				"B5": "-",
			},
			wantMerge: []string{"A2:A3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &excelMockedWriter{}
			f := &ExcelFormatter{
				config: &Config{
					Writer: writer,
				},
			}
			if err := f.Format(tt.args.td, ""); (err != nil) != tt.wantErr {
				t.Fatalf("ExcelFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			wb, err := excelize.OpenReader(bytes.NewReader(writer.Bytes()))
			if err != nil {
				t.Fatalf("excelize.OpenReader() error = %v", err)
			}
			defer func() {
				if closeErr := wb.Close(); closeErr != nil {
					t.Errorf("wb.Close() error = %v", closeErr)
				}
			}()

			const sheet = "Sheet1"
			for cell, wantVal := range tt.wantCells {
				gotVal, err := wb.GetCellValue(sheet, cell)
				if err != nil {
					t.Fatalf("GetCellValue(%s) error = %v", cell, err)
				}
				if gotVal != wantVal {
					t.Errorf("cell %s = %q, want %q", cell, gotVal, wantVal)
				}
			}

			for _, col := range []string{"A", "B"} {
				width, err := wb.GetColWidth(sheet, col)
				if err != nil {
					t.Fatalf("GetColWidth(%s) error = %v", col, err)
				}
				if width != ExcelColWidth {
					t.Errorf("col %s width = %v, want %v", col, width, ExcelColWidth)
				}
			}

			styleID, err := wb.GetCellStyle(sheet, "A1")
			if err != nil {
				t.Fatalf("GetCellStyle(A1) error = %v", err)
			}
			style, err := wb.GetStyle(styleID)
			if err != nil {
				t.Fatalf("GetStyle(%d) error = %v", styleID, err)
			}
			if style.Alignment == nil ||
				style.Alignment.Horizontal != "center" ||
				style.Alignment.Vertical != "center" ||
				!style.Alignment.WrapText {
				t.Errorf("unexpected A1 alignment style: %+v", style.Alignment)
			}

			merges, err := wb.GetMergeCells(sheet)
			if err != nil {
				t.Fatalf("GetMergeCells() error = %v", err)
			}
			if len(merges) != len(tt.wantMerge) {
				t.Fatalf("len(GetMergeCells()) = %d, want %d", len(merges), len(tt.wantMerge))
			}
			for i, wantRange := range tt.wantMerge {
				gotRange := merges[i].GetStartAxis() + ":" + merges[i].GetEndAxis()
				if gotRange != wantRange {
					t.Errorf("merge[%d] = %q, want %q", i, gotRange, wantRange)
				}
			}
		})
	}
}
