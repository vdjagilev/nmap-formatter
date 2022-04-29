package formatter

import (
	_ "embed"
	"regexp"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

type testHTMLMockedFormatterWriter struct {
	data []byte
	err  error
}

func (f *testHTMLMockedFormatterWriter) Write(p []byte) (n int, err error) {
	if f.err != nil {
		return 0, f.err
	}
	f.data = append(f.data, p...)
	return len(p), nil
}

func TestHTMLFormatter_Format(t *testing.T) {
	type args struct {
		td *TemplateData
	}
	tests := []struct {
		name     string
		f        *HTMLFormatter
		args     args
		wantErr  bool
		validate func(f *HTMLFormatter, output string, t *testing.T)
	}{
		{
			name: "html successfully converted",
			f:    &HTMLFormatter{},
			args: args{
				td: &TemplateData{
					NMAPRun:       NMAPRun{},
					OutputOptions: OutputOptions{},
				},
			},
			wantErr: false,
		},
		{
			name: "if basic html is parseable",
			f:    &HTMLFormatter{},
			args: args{
				td: &TemplateData{
					NMAPRun:       NMAPRun{},
					OutputOptions: OutputOptions{},
				},
			},
			wantErr: false,
			validate: func(f *HTMLFormatter, output string, t *testing.T) {
				_, err := html.Parse(strings.NewReader(output))
				if err != nil {
					t.Fatalf("Failed to parse HTML: %v", err)
				}
			},
		},
		{
			name: "title check",
			f:    &HTMLFormatter{},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						StartStr: "start time",
					},
					OutputOptions: OutputOptions{},
				},
			},
			wantErr: false,
			validate: func(f *HTMLFormatter, output string, t *testing.T) {
				expect := "start time"
				re := regexp.MustCompile(`<title>NMAP Scan result: (?P<time>.+)</title>`)
				time := re.FindStringSubmatch(output)[re.SubexpIndex("time")]
				if time != expect {
					t.Fatalf("Expected <title> to have start-time = %v, got = %v", expect, time)
				}
			},
		},
		{
			name: "check TOC 2 hosts up",
			f:    &HTMLFormatter{},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Host: []Host{
							{
								StartTime: 0,
								EndTime:   0,
								Port:      []Port{},
								HostAddress: HostAddress{
									Address: "192.168.1.1",
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
								HostAddress: HostAddress{
									Address: "192.168.1.2",
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
						HTMLOptions: HTMLOutputOptions{
							SkipDownHosts: true,
						},
					},
				},
			},
			wantErr: false,
			validate: func(f *HTMLFormatter, output string, t *testing.T) {
				expect := 2
				re := regexp.MustCompile(`<li><a href="#192\.168\.1\.([12])">`)
				actual := len(re.FindAllString(output, -1))
				if expect != actual {
					t.Fatalf("Expected %d addresses in TOC, got %d", expect, actual)
				}
			},
		},
		{
			name: "Check TOC 2 hosts, 1 down",
			f:    &HTMLFormatter{},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Host: []Host{
							{
								StartTime: 0,
								EndTime:   0,
								Port:      []Port{},
								HostAddress: HostAddress{
									Address: "192.168.1.1",
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
								HostAddress: HostAddress{
									Address: "192.168.1.2",
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
						HTMLOptions: HTMLOutputOptions{
							SkipDownHosts: true,
						},
					},
				},
			},
			wantErr: false,
			validate: func(f *HTMLFormatter, output string, t *testing.T) {
				expect := 1
				re := regexp.MustCompile(`<li><a href="#192\.168\.1\.([12])">`)
				actual := len(re.FindAllString(output, -1))
				if expect != actual {
					t.Fatalf("Expected %d addresses in TOC, got %d", expect, actual)
				}
			},
		},
		{
			name: "Check TOC 2 hosts, 1 down, not skipping down hosts",
			f:    &HTMLFormatter{},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Host: []Host{
							{
								StartTime: 0,
								EndTime:   0,
								Port:      []Port{},
								HostAddress: HostAddress{
									Address: "192.168.1.1",
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
								HostAddress: HostAddress{
									Address: "192.168.1.2",
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
						HTMLOptions: HTMLOutputOptions{
							SkipDownHosts: false,
						},
					},
				},
			},
			wantErr: false,
			validate: func(f *HTMLFormatter, output string, t *testing.T) {
				expect := 2
				re := regexp.MustCompile(`<li><a href="#192\.168\.1\.([12])">`)
				actual := len(re.FindAllString(output, -1))
				if expect != actual {
					t.Fatalf("Expected %d addresses in TOC, got %d", expect, actual)
				}
			},
		},
		{
			name: "Check for ports, 2 hosts",
			f:    &HTMLFormatter{},
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
											State: "open",
										},
										Service: PortService{},
										Script:  []Script{},
									},
									{
										Protocol: "tcp",
										PortID:   443,
										State: PortState{
											State: "up",
										},
										Service: PortService{},
										Script:  []Script{},
									},
								},
								HostAddress: HostAddress{
									Address: "192.168.1.1",
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
										PortID:   8080,
										State: PortState{
											State: "open",
										},
										Service: PortService{},
										Script:  []Script{},
									},
								},
								HostAddress: HostAddress{},
								HostNames:   HostNames{},
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
			validate: func(f *HTMLFormatter, output string, t *testing.T) {
				ports := []string{"80", "443", "8080"}
				for _, p := range ports {
					re := regexp.MustCompile(`<td class="port-\w+">` + p)
					found := re.FindString(output)
					if found == "" {
						t.Fatalf("The port %s was not found in port listing tables", p)
					}
				}
			},
		},
		{
			name: "3 hosts 1 down (skip down: false)",
			f:    &HTMLFormatter{},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Host: []Host{
							{
								StartTime: 0,
								EndTime:   0,
								Port:      []Port{},
								HostAddress: HostAddress{
									Address: "192.168.1.1",
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
								HostAddress: HostAddress{
									Address: "192.168.1.2",
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
								HostAddress: HostAddress{
									Address: "192.168.1.3",
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
						HTMLOptions: HTMLOutputOptions{
							SkipDownHosts: false,
						},
					},
				},
			},
			wantErr: false,
			validate: func(f *HTMLFormatter, output string, t *testing.T) {
				expect := 3
				re := regexp.MustCompile(`<h2 class="host-address-header (host-up|host-down)">(\s*)192\.168\.1\.[0-9]+`)
				actual := len(re.FindAllString(output, -1))
				if expect != actual {
					t.Fatalf("Expected %d host headers, got %d", expect, actual)
				}
			},
		},
		{
			name: "3 hosts 1 down (skip down: true)",
			f:    &HTMLFormatter{},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Host: []Host{
							{
								StartTime: 0,
								EndTime:   0,
								Port:      []Port{},
								HostAddress: HostAddress{
									Address: "192.168.1.1",
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
								HostAddress: HostAddress{
									Address: "192.168.1.2",
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
								HostAddress: HostAddress{
									Address: "192.168.1.3",
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
						HTMLOptions: HTMLOutputOptions{
							SkipDownHosts: true,
						},
					},
				},
			},
			wantErr: false,
			validate: func(f *HTMLFormatter, output string, t *testing.T) {
				expect := 2
				re := regexp.MustCompile(`<h2 class="host-address-header (host-up|host-down)">(\s*)192\.168\.1\.[0-9]+`)
				actual := len(re.FindAllString(output, -1))
				if expect != actual {
					t.Fatalf("Expected %d host headers, got %d", expect, actual)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := testHTMLMockedFormatterWriter{}
			tt.f.config = &Config{
				Writer: &writer,
			}
			if err := tt.f.Format(tt.args.td, HTMLSimpleTemplate); (err != nil) != tt.wantErr {
				t.Errorf("HTMLFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
			} else if tt.validate != nil {
				tt.validate(tt.f, string(writer.data), t)
			}
		})
	}
}
