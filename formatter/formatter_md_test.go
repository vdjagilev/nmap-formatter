package formatter

import (
	_ "embed"
	"html/template"
	"os"
	"reflect"
	"regexp"
	"testing"
)

func Test_markdownEntry(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Content without ticks should remain the same",
			args: args{
				v: `## Some content that is already defined
				Should remain the same`,
			},
			want: `## Some content that is already defined
				Should remain the same`,
		},
		{
			name: "Remove tick from 'ticked' sentence",
			args: args{
				v: "Let's remove this `part` of code",
			},
			want: "Let's remove this part of code",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := markdownEntry(tt.args.v); got != tt.want {
				t.Errorf("markdownEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_markdownNoEscape(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want template.HTML
	}{
		{
			name: "Basic tick removal",
			args: args{
				v: "Let's remove this `part` of code",
			},
			want: "Let's remove this part of code",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := markdownNoEscape(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("markdownNoEscape() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_markdownTOCEntry(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "IP address",
			args: args{
				v: "192.168.1.1",
			},
			want: "19216811",
		},
		{
			name: "Sentence",
			args: args{
				v: "Lorem Ipsum",
			},
			want: "lorem-ipsum",
		},
		{
			name: "IP addr & sentece",
			args: args{
				v: "192.168.2.2 Test Host",
			},
			want: "19216822-test-host",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := markdownTOCEntry(tt.args.v); got != tt.want {
				t.Errorf("markdownTOCEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

type markdownMockedWriter struct {
	data []byte
	err  error
}

func (w *markdownMockedWriter) Write(p []byte) (n int, err error) {
	if w.err != nil {
		return 0, w.err
	}
	w.data = append(w.data, p...)
	return len(p), nil
}

func TestMarkdownFormatter_Format(t *testing.T) {
	type args struct {
		td *TemplateData
	}
	tests := []struct {
		name     string
		f        *MarkdownFormatter
		args     args
		wantErr  bool
		validate func(f *MarkdownFormatter, output string, t *testing.T)
	}{
		{
			name: "Basic execution",
			f: &MarkdownFormatter{
				config: &Config{
					Writer: os.Stdout,
				},
			},
			args: args{
				td: &TemplateData{
					NMAPRun:       NMAPRun{},
					OutputOptions: OutputOptions{},
				},
			},
			wantErr: false,
		},
		{
			name: "Have 3 hosts (1 is down, skip down: true)",
			f:    &MarkdownFormatter{},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Host: []Host{
							{
								StartTime: "",
								EndTime:   "",
								Ports:     Ports{},
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
								StartTime: "",
								EndTime:   "",
								Ports:     Ports{},
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
								StartTime: "",
								EndTime:   "",
								Ports:     Ports{},
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
						SkipDownHosts: true,
					},
				},
			},
			wantErr: false,
			validate: func(f *MarkdownFormatter, output string, t *testing.T) {
				expect := 2
				re := regexp.MustCompile(`## 192\.168\.1\.[0-9]+`)
				actual := len(re.FindAllString(output, -1))
				if expect != actual {
					t.Fatalf("Expected %d host headers, got %d", expect, actual)
				}
			},
		},
		{
			name: "Have 3 hosts (1 is down, skip down: false)",
			f:    &MarkdownFormatter{},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Host: []Host{
							{
								StartTime: "",
								EndTime:   "",
								Ports:     Ports{},
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
								StartTime: "",
								EndTime:   "",
								Ports:     Ports{},
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
								StartTime: "",
								EndTime:   "",
								Ports:     Ports{},
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
						SkipDownHosts: false,
					},
				},
			},
			wantErr: false,
			validate: func(f *MarkdownFormatter, output string, t *testing.T) {
				expect := 3
				re := regexp.MustCompile(`## 192\.168\.1\.[0-9]+`)
				actual := len(re.FindAllString(output, -1))
				if expect != actual {
					t.Fatalf("Expected %d host headers, got %d", expect, actual)
				}
			},
		},
		{
			name: "Have 3 ports (1 host is down, skip down: true)",
			f:    &MarkdownFormatter{},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Host: []Host{
							{
								StartTime: "",
								EndTime:   "",
								Ports: Ports{
									Port: []Port{
										{
											Protocol: "tcp",
											PortID:   "80",
											State: PortState{
												State: "open",
											},
											Service: PortService{},
											Script:  []Script{},
										},
										{
											Protocol: "tcp",
											PortID:   "443",
											State: PortState{
												State: "up",
											},
											Service: PortService{},
											Script:  []Script{},
										},
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
								StartTime: "",
								EndTime:   "",
								Ports:     Ports{},
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
								StartTime: "",
								EndTime:   "",
								Ports: Ports{
									Port: []Port{
										{
											Protocol: "tcp",
											PortID:   "8080",
											State: PortState{
												State: "open",
											},
											Service: PortService{},
											Script:  []Script{},
										},
									},
								},
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
						SkipDownHosts: true,
					},
				},
			},
			wantErr: false,
			validate: func(f *MarkdownFormatter, output string, t *testing.T) {
				ports := []string{"80", "443", "8080"}
				for _, p := range ports {
					re := regexp.MustCompile(`\| ` + p + ` \| tcp`)
					found := re.FindString(output)
					if found == "" {
						t.Fatalf("The port %s was not found in port listing tables", p)
					}
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := markdownMockedWriter{}
			tt.f.config = &Config{
				Writer: &writer,
			}
			if err := tt.f.Format(tt.args.td); (err != nil) != tt.wantErr {
				t.Errorf("MarkdownFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
			} else if tt.validate != nil {
				tt.validate(tt.f, string(writer.data), t)
			}
		})
	}
}
