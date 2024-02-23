package formatter

import "testing"

type excelMockedWriter struct {
	data []byte
}

func (w *excelMockedWriter) Write(p []byte) (n int, err error) {
	w.data = p
	return len(p), nil
}

func (w *excelMockedWriter) Close() error {
	return nil
}

func TestExcelFormatter_Format(t *testing.T) {
	writer := &excelMockedWriter{}
	type args struct {
		td *TemplateData
	}
	tests := []struct {
		name    string
		f       *ExcelFormatter
		args    args
		wantErr bool
	}{
		{
			name: "Empty Test ExcelFormatter.Format",
			f: &ExcelFormatter{
				config: &Config{
					Writer: writer,
				},
			},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Host: []Host{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test ExcelFormatter.Format Basic 3 hosts",
			f: &ExcelFormatter{
				config: &Config{
					Writer: writer,
				},
			},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Host: []Host{
							{
								StartTime: 0,
								EndTime:   0,
								Port:      []Port{},
								HostAddress: []HostAddress{
									{
										Address: "192.168.1.1",
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.Format(tt.args.td, ""); (err != nil) != tt.wantErr {
				t.Errorf("ExcelFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
