package formatter

import "testing"

type d2MockedWriter struct {
	data []byte
	err  error
}

func (d *d2MockedWriter) Write(p []byte) (n int, err error) {
	if d.err != nil {
		return 0, d.err
	}
	d.data = append(d.data, p...)
	return len(p), nil
}

func (d *d2MockedWriter) Close() error {
	return nil
}

func TestD2LangFormatter_Format(t *testing.T) {
	type fields struct {
		config *Config
	}
	type args struct {
		td              *TemplateData
		templateContent string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Basic",
			fields: fields{
				config: &Config{
					Writer: &d2MockedWriter{},
				},
			},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Scanner:  "",
						Args:     "",
						Start:    0,
						StartStr: "",
						Version:  "",
						ScanInfo: ScanInfo{},
						Host: []Host{
							{
								StartTime: 0,
								EndTime:   0,
								Port: []Port{
									{
										Protocol: "",
										PortID:   22,
										State: PortState{
											State:     "open",
											Reason:    "",
											ReasonTTL: "",
										},
										Service: PortService{
											Name: "ssh",
										},
										Script: []Script{},
									},
								},
								HostAddress: []HostAddress{
									{
										Address:     "10.10.10.1",
										AddressType: "ipv4",
										Vendor:      "",
									},
								},
								HostNames:     HostNames{},
								Status:        HostStatus{},
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
										Protocol: "",
										PortID:   80,
										State: PortState{
											State:     "open",
											Reason:    "",
											ReasonTTL: "",
										},
										Service: PortService{
											Name:      "",
											Product:   "",
											Version:   "",
											ExtraInfo: "",
											Method:    "",
											Conf:      "",
											CPE:       []string{},
										},
										Script: []Script{},
									},
								},
								HostAddress: []HostAddress{
									{
										Address:     "10.10.10.2",
										AddressType: "ipv4",
										Vendor:      "",
									},
								},
								HostNames:     HostNames{},
								Status:        HostStatus{},
								OS:            OS{},
								Trace:         Trace{},
								Uptime:        Uptime{},
								Distance:      Distance{},
								TCPSequence:   TCPSequence{},
								IPIDSequence:  IPIDSequence{},
								TCPTSSequence: TCPTSSequence{},
							},
						},
						Verbose:   Verbose{},
						Debugging: Debugging{},
						RunStats:  RunStats{},
					},
					OutputOptions: OutputOptions{},
					CustomOptions: map[string]string{},
				},
				templateContent: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &D2LangFormatter{
				config: tt.fields.config,
			}
			if err := f.Format(tt.args.td, tt.args.templateContent); (err != nil) != tt.wantErr {
				t.Errorf("D2LangFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
