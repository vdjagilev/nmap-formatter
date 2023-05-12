package formatter

import (
	"testing"

	// including go-sqlite3 to for unit tests in memory
	_ "github.com/mattn/go-sqlite3"
)

func TestSqliteFormatter_Format(t *testing.T) {
	const DBDSN = "file::memory:?cache=shared"
	var config = &Config{
		OutputOptions: OutputOptions{
			SqliteOutputOptions: SqliteOutputOptions{
				DSN: DBDSN,
			},
		},
		CurrentVersion: "1",
	}
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
			name: "No data",
			fields: fields{
				config,
			},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{},
				},
				templateContent: "",
			},
			wantErr: false,
		},
		{
			name: "1 scan, 2 hosts, 2 ports",
			fields: fields{
				config,
			},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Scanner: "nmap",
						Host: []Host{
							{
								StartTime: 0,
								EndTime:   0,
								Port: []Port{
									{
										Protocol: "tcp",
										PortID:   80,
										State:    PortState{},
										Service:  PortService{},
										Script: []Script{
											{
												ID:     "http",
												Output: "abc",
											},
										},
									}, {
										Protocol: "tcp",
										PortID:   443,
										State:    PortState{},
										Service:  PortService{},
										Script: []Script{
											{
												ID:     "https",
												Output: "abc",
											},
										},
									},
								},
								HostAddress: []HostAddress{
									{
										Address:     "192.168.1.1",
										AddressType: "ipv4",
									},
									{
										Address:     "192.168.2.1",
										AddressType: "ipv4",
									},
								},
								HostNames: HostNames{
									HostName: []HostName{
										{
											Name: "example.com",
											Type: "A",
										},
										{
											Name: "www.example.com",
											Type: "AA",
										},
									},
								},
								Status: HostStatus{},
								OS: OS{
									OSPortUsed: []OSPortUsed{
										{
											State:    "up",
											Protocol: "tcp",
											PortID:   80,
										},
										{
											State:    "filtered",
											Protocol: "tcp",
											PortID:   443,
										},
									},
									OSClass: OSClass{
										Type:     "a",
										Vendor:   "b",
										OSFamily: "c",
										OSGen:    "d",
										Accuracy: "e",
										CPE:      []string{"a", "b"},
									},
									OSMatch: []OSMatch{
										{
											Name:     "a",
											Accuracy: "1",
											Line:     "1",
										},
									},
								},
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
										PortID:   80,
										State:    PortState{},
										Service:  PortService{},
										Script: []Script{
											{
												ID:     "http",
												Output: "abc",
											},
										},
									}, {
										Protocol: "tcp",
										PortID:   443,
										State:    PortState{},
										Service:  PortService{},
										Script: []Script{
											{
												ID:     "https",
												Output: "abc",
											},
										},
									},
								},
								HostAddress: []HostAddress{
									{
										Address:     "192.168.1.1",
										AddressType: "ipv4",
									},
									{
										Address:     "192.168.2.1",
										AddressType: "ipv4",
									},
								},
								HostNames: HostNames{
									HostName: []HostName{
										{
											Name: "example.com",
											Type: "A",
										},
										{
											Name: "www.example.com",
											Type: "AA",
										},
									},
								},
								Status: HostStatus{},
								OS:     OS{},
								Trace: Trace{
									Port:     80,
									Protocol: "tcp",
									Hops: []Hop{
										{
											TTL:    2,
											IPAddr: "10.0.0.1",
											RTT:    0,
											Host:   "",
										},
										{
											TTL:    1,
											IPAddr: "192.168.1.1",
											RTT:    0,
											Host:   "",
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
				},
				templateContent: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &SqliteFormatter{
				config: tt.fields.config,
			}
			if err := f.Format(tt.args.td, tt.args.templateContent); (err != nil) != tt.wantErr {
				t.Errorf("SqliteFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
