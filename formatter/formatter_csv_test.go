package formatter

import (
	"reflect"
	"testing"
)

func TestCSVFormatter_convert(t *testing.T) {
	header := []string{"IP", "Port", "Protocol", "State", "Service", "Reason", "Product", "Version", "Extra info"}
	type args struct {
		td *TemplateData
	}
	tests := []struct {
		name     string
		f        *CSVFormatter
		args     args
		wantData [][]string
	}{
		{
			name: "Empty CSV",
			f:    &CSVFormatter{},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Host: []Host{},
					},
				},
			},
			wantData: [][]string{
				header,
			},
		},
		{
			name: "1 Host is down",
			f:    &CSVFormatter{},
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
										Address:     "127.0.0.1",
										AddressType: "ipv4",
										Vendor:      "",
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
						},
					},
					OutputOptions: OutputOptions{
						CSVOptions: CSVOutputOptions{},
					},
				},
			},
			wantData: [][]string{
				header,
				{"127.0.0.1 (down)", "", "", "", "", "", "", "", ""},
			},
		},
		{
			name: "1 host is down",
			f:    &CSVFormatter{},
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
										Address: "127.0.0.1",
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
						},
					},
					OutputOptions: OutputOptions{
						CSVOptions: CSVOutputOptions{},
					},
				},
			},
			wantData: [][]string{
				header,
				{"127.0.0.1 (down)", "", "", "", "", "", "", "", ""},
			},
		},
		{
			name: "1 host 1 port (up)",
			f:    &CSVFormatter{},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Host: []Host{
							{
								StartTime: 0,
								EndTime:   0,
								Port: []Port{
									{
										Protocol: "tcp",
										PortID:   80,
										State: PortState{
											State:  "open",
											Reason: "syn-ack",
										},
										Service: PortService{
											Name:    "http",
											Product: "nginx",
											Version: "1.21.1",
										},
										Script: []Script{},
									},
								},
								HostAddress: []HostAddress{
									{
										Address: "127.0.0.1",
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
					OutputOptions: OutputOptions{
						CSVOptions: CSVOutputOptions{},
					},
				},
			},
			wantData: [][]string{
				header,
				{"127.0.0.1 (up)", "", "", "", "", "", "", "", ""},
				{"", "80", "tcp", "open", "http", "syn-ack", "nginx", "1.21.1", ""},
			},
		},
		{
			name: "1 host 2 ports (up)",
			f:    &CSVFormatter{},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Host: []Host{
							{
								StartTime: 0,
								EndTime:   0,
								Port: []Port{
									{
										Protocol: "tcp",
										PortID:   80,
										State: PortState{
											State:  "open",
											Reason: "syn-ack",
										},
										Service: PortService{
											Name:    "http",
											Product: "nginx",
											Version: "1.21.1",
										},
										Script: []Script{},
									},
									{
										Protocol: "tcp",
										PortID:   443,
										State: PortState{
											State:  "open",
											Reason: "syn-ack",
										},
										Service: PortService{
											Name:    "http",
											Product: "nginx",
											Version: "1.21.1",
										},
										Script: []Script{},
									},
								},
								HostAddress: []HostAddress{
									{
										Address: "127.0.0.1",
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
					OutputOptions: OutputOptions{
						CSVOptions: CSVOutputOptions{},
					},
				},
			},
			wantData: [][]string{
				header,
				{"127.0.0.1 (up)", "", "", "", "", "", "", "", ""},
				{"", "80", "tcp", "open", "http", "syn-ack", "nginx", "1.21.1", ""},
				{"", "443", "tcp", "open", "http", "syn-ack", "nginx", "1.21.1", ""},
			},
		},
		{
			name: "1 host up 2 ports, 1 host down",
			f:    &CSVFormatter{},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Host: []Host{
							{
								StartTime: 0,
								EndTime:   0,
								Port: []Port{
									{
										Protocol: "tcp",
										PortID:   80,
										State: PortState{
											State:  "open",
											Reason: "syn-ack",
										},
										Service: PortService{
											Name:    "http",
											Product: "nginx",
											Version: "1.21.1",
										},
										Script: []Script{},
									},
									{
										Protocol: "tcp",
										PortID:   443,
										State: PortState{
											State:  "open",
											Reason: "syn-ack",
										},
										Service: PortService{
											Name:    "http",
											Product: "nginx",
											Version: "1.21.1",
										},
										Script: []Script{},
									},
								},
								HostAddress: []HostAddress{
									{
										Address: "127.0.0.1",
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
										Address: "192.168.1.1",
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
						},
					},
					OutputOptions: OutputOptions{
						CSVOptions: CSVOutputOptions{},
					},
				},
			},
			wantData: [][]string{
				header,
				{"127.0.0.1 (up)", "", "", "", "", "", "", "", ""},
				{"", "80", "tcp", "open", "http", "syn-ack", "nginx", "1.21.1", ""},
				{"", "443", "tcp", "open", "http", "syn-ack", "nginx", "1.21.1", ""},
				{"192.168.1.1 (down)", "", "", "", "", "", "", "", ""},
			},
		},
		{
			name: "1 host up 2 ports, 1 host down",
			f: &CSVFormatter{
				config: &Config{},
			},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Host: []Host{
							{
								StartTime: 0,
								EndTime:   0,
								Port: []Port{
									{
										Protocol: "tcp",
										PortID:   80,
										State: PortState{
											State:  "open",
											Reason: "syn-ack",
										},
										Service: PortService{
											Name:    "http",
											Product: "nginx",
											Version: "1.21.1",
										},
										Script: []Script{},
									},
									{
										Protocol: "tcp",
										PortID:   443,
										State: PortState{
											State:  "open",
											Reason: "syn-ack",
										},
										Service: PortService{
											Name:    "http",
											Product: "nginx",
											Version: "1.21.1",
										},
										Script: []Script{},
									},
								},
								HostAddress: []HostAddress{
									{
										Address: "127.0.0.1",
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
										Address: "192.168.1.1",
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
						},
					},
					OutputOptions: OutputOptions{
						CSVOptions: CSVOutputOptions{},
					},
				},
			},
			wantData: [][]string{
				header,
				{"127.0.0.1 (up)", "", "", "", "", "", "", "", ""},
				{"", "80", "tcp", "open", "http", "syn-ack", "nginx", "1.21.1", ""},
				{"", "443", "tcp", "open", "http", "syn-ack", "nginx", "1.21.1", ""},
				{"192.168.1.1 (down)", "", "", "", "", "", "", "", ""},
			},
		},
		{
			name: "2 hosts (2+1 ports)",
			f: &CSVFormatter{
				config: &Config{
					SkipDownHosts: true,
				},
			},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Host: []Host{
							{
								StartTime: 0,
								EndTime:   0,
								Port: []Port{
									{
										Protocol: "tcp",
										PortID:   80,
										State: PortState{
											State:  "open",
											Reason: "syn-ack",
										},
										Service: PortService{
											Name:    "http",
											Product: "nginx",
											Version: "1.21.1",
										},
										Script: []Script{},
									},
									{
										Protocol: "tcp",
										PortID:   443,
										State: PortState{
											State:  "open",
											Reason: "syn-ack",
										},
										Service: PortService{
											Name:    "http",
											Product: "nginx",
											Version: "1.21.1",
										},
										Script: []Script{},
									},
								},
								HostAddress: []HostAddress{
									{
										Address: "127.0.0.1",
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
								Port: []Port{
									{
										Protocol: "tcp",
										PortID:   22,
										State: PortState{
											State:  "open",
											Reason: "syn-ack",
										},
										Service: PortService{
											Name:    "ssh",
											Product: "OpenSSH",
											Version: "5.3p1 Debian 3ubuntu7",
										},
										Script: []Script{},
									},
								},
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
						},
					},
					OutputOptions: OutputOptions{
						CSVOptions: CSVOutputOptions{},
					},
				},
			},
			wantData: [][]string{
				header,
				{"127.0.0.1 (up)", "", "", "", "", "", "", "", ""},
				{"", "80", "tcp", "open", "http", "syn-ack", "nginx", "1.21.1", ""},
				{"", "443", "tcp", "open", "http", "syn-ack", "nginx", "1.21.1", ""},
				{"192.168.1.1 (up)", "", "", "", "", "", "", "", ""},
				{"", "22", "tcp", "open", "ssh", "syn-ack", "OpenSSH", "5.3p1 Debian 3ubuntu7", ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotData := tt.f.convert(tt.args.td); !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("CSVFormatter.convert() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func TestCSVFormatter_Format(t *testing.T) {
	writer := &csvMockedWriter{}
	type args struct {
		td *TemplateData
	}
	tests := []struct {
		name       string
		f          *CSVFormatter
		args       args
		wantErr    bool
		wantOutput string
	}{
		{
			name: "Successful header write",
			f: &CSVFormatter{
				config: &Config{
					Writer:        writer,
					SkipDownHosts: true,
				},
			},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Host: []Host{},
					},
					OutputOptions: OutputOptions{
						CSVOptions: CSVOutputOptions{},
					},
				},
			},
			wantErr:    false,
			wantOutput: "IP,Port,Protocol,State,Service,Reason,Product,Version,Extra info\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.Format(tt.args.td, ""); (err != nil) != tt.wantErr {
				t.Errorf("CSVFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantOutput != string(writer.data) {
				t.Errorf("CSVFormatter.Format() written data = %v, wantOutput = %v", string(writer.data), tt.wantOutput)
			}
		})
	}
}

type csvMockedWriter struct {
	data []byte
}

func (w *csvMockedWriter) Write(p []byte) (n int, err error) {
	w.data = p
	return len(p), nil
}

func (w *csvMockedWriter) Close() error {
	return nil
}
