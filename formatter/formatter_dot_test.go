package formatter

import (
	_ "embed"
	"reflect"
	"strings"
	"testing"
)

func Test_hopList(t *testing.T) {
	type args struct {
		hops       []Hop
		startHop   string
		endHopName string
		endHopKey  int
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "No hops",
			args: args{
				hops:       []Hop{},
				startHop:   "scanner",
				endHopName: "srv",
				endHopKey:  0,
			},
			want: map[string]string{
				"scanner": "srv0",
			},
		},
		{
			name: "3 Hops",
			args: args{
				hops: []Hop{
					{
						IPAddr: "192.168.250.1",
					},
					{
						IPAddr: "192.168.1.1",
					},
					{
						IPAddr: "10.10.10.1",
					},
				},
				startHop:   "scanner",
				endHopName: "srv",
				endHopKey:  0,
			},
			want: map[string]string{
				"scanner":          "hop192.168.250.1",
				"hop192.168.250.1": "hop192.168.1.1",
				"hop192.168.1.1":   "srv0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hopList(tt.args.hops, tt.args.startHop, tt.args.endHopName, tt.args.endHopKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hopList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_portStateColor(t *testing.T) {
	type args struct {
		port *Port
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Default color",
			args: args{
				port: &Port{
					State: PortState{
						State: "unknown",
					},
				},
			},
			want: "gray",
		},
		{
			name: "Open port color",
			args: args{
				port: &Port{
					State: PortState{
						State: "open",
					},
				},
			},
			want: "#228B22",
		},
		{
			name: "Filtered port color",
			args: args{
				port: &Port{
					State: PortState{
						State: "filtered",
					},
				},
			},
			want: "#FFAE00",
		},
		{
			name: "Closed port color",
			args: args{
				port: &Port{
					State: PortState{
						State: "closed",
					},
				},
			},
			want: "#DC143C",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := portStateColor(tt.args.port); got != tt.want {
				t.Errorf("portStateColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cleanIP(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Local network IP",
			args: args{
				ip: "192.168.1.1",
			},
			want: "19216811",
		},
		{
			name: "Loopback",
			args: args{
				ip: "127.0.0.1",
			},
			want: "127001",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanIP(tt.args.ip); got != tt.want {
				t.Errorf("cleanIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

type dotMockedWriter struct {
	data []byte
	err  error
}

func (w *dotMockedWriter) Write(p []byte) (n int, err error) {
	if w.err != nil {
		return 0, w.err
	}
	w.data = append(w.data, p...)
	return len(p), nil
}

func (w *dotMockedWriter) Close() error {
	return nil
}

func TestDotFormatter_Format(t *testing.T) {
	type fields struct {
		config *Config
	}
	type args struct {
		td              *TemplateData
		templateContent string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		validate func(f *DotFormatter, output string, t *testing.T)
		wantErr  bool
	}{
		{
			name: "Failed to parse template",
			fields: fields{
				&Config{
					Writer: &dotMockedWriter{},
				},
			},
			args: args{
				td:              &TemplateData{},
				templateContent: "wrong template content {{.nonexistent_variable}}",
			},
			validate: func(f *DotFormatter, output string, t *testing.T) {
			},
			wantErr: true,
		},
		{
			name: "1 target, 2 hops",
			fields: fields{
				&Config{
					Writer: &dotMockedWriter{},
				},
			},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Scanner: "nmap",
						Host: []Host{
							{
								StartTime: 0,
								EndTime:   0,
								Port:      []Port{},
								HostAddress: []HostAddress{
									{
										Address:     "10.10.10.12",
										AddressType: "ipv4",
									},
								},
								HostNames: HostNames{},
								Status: HostStatus{
									State: "up",
								},
								OS: OS{},
								Trace: Trace{
									Port:     20,
									Protocol: "tcp",
									Hops: []Hop{
										{
											IPAddr: "192.168.100.1",
										},
										{
											IPAddr: "192.168.1.1",
										},
										{
											IPAddr: "10.10.10.12",
										},
									},
								},
								Uptime:        Uptime{},
								Distance:      Distance{},
								TCPSequence:   TCPSequence{},
								IPIDSequence:  IPIDSequence{},
								TCPTSSequence: TCPTSSequence{},
							},
						},
					},
					CustomOptions: DotDefaultOptions,
				},
				templateContent: DotTemplate,
			},
			wantErr: false,
			validate: func(f *DotFormatter, output string, t *testing.T) {
				if !strings.Contains(output, "srv0 [label=\"10.10.10.12\", tooltip=\"10.10.10.12\", shape=hexagon, style=filled];") {
					t.Error("Does not contain correct sv0")
				}
				hops := []string{
					"hop1921681001 [label=\"\", tooltip=\"192.168.100.1\", shape=circle, height=.12, width=.12, style=filled];",
					"hop19216811 [label=\"\", tooltip=\"192.168.1.1\", shape=circle, height=.12, width=.12, style=filled];",
				}
				for i := range hops {
					if !strings.Contains(output, hops[i]) {
						t.Errorf("Could not find hop: %s", hops[i])
					}
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := &dotMockedWriter{}
			f := &DotFormatter{
				config: &Config{
					Writer: output,
				},
			}
			if err := f.Format(tt.args.td, tt.args.templateContent); (err != nil) != tt.wantErr {
				t.Errorf("DotFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
			} else if tt.validate != nil {
				tt.validate(f, string(output.data), t)
			}
		})
	}
}
