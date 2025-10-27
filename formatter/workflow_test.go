package formatter

import (
	"os"
	"path"
	"reflect"
	"testing"
)

func TestMainWorkflow_parse(t *testing.T) {
	tests := []struct {
		name        string
		w           *MainWorkflow
		wantNMAPRun NMAPRun
		wantErr     bool
		fileContent string
		fileName    string
	}{
		{
			name: "Wrong path (file does not exists)",
			w: &MainWorkflow{
				Config: &Config{
					InputFileConfig: InputFileConfig{
						Path: "",
					},
				},
			},
			wantNMAPRun: NMAPRun{},
			wantErr:     true,
		},
		{
			name: "Non-xml file",
			w: &MainWorkflow{
				Config: &Config{}, // will be set dynamically
			},
			wantNMAPRun: NMAPRun{},
			wantErr:     true,
			fileName:    "main_workflow_parse_2_test",
			fileContent: "[NOT XML file]",
		},
		{
			name: "XML file (empty content)",
			w: &MainWorkflow{
				Config: &Config{},
			},
			wantNMAPRun: NMAPRun{},
			wantErr:     false,
			fileContent: `<?xml version="1.0"?>
			<?xml-stylesheet href="file:///usr/local/bin/../share/nmap/nmap.xsl" type="text/xsl"?>
			<nmaprun></nmaprun>`,
			fileName: "main_workflow_parse_3_test",
		},
		{
			name: "Bad XML file",
			w: &MainWorkflow{
				Config: &Config{},
			},
			wantNMAPRun: NMAPRun{},
			wantErr:     true,
			fileContent: "<?x< version=",
			fileName:    "main_workflow_parse_4_test_wrong_xml",
		},
		{
			name: "XML file with some matching output",
			w: &MainWorkflow{
				Config: &Config{},
			},
			wantNMAPRun: NMAPRun{
				Scanner: "nmap",
				Version: "5.59BETA3",
				ScanInfo: ScanInfo{
					Services: "1-1000",
				},
			},
			wantErr: false,
			fileContent: `<?xml version="1.0"?>
			<nmaprun scanner="nmap" version="5.59BETA3">
				<scaninfo services="1-1000"/>
			</nmaprun>`,
			fileName: "main_workflow_parse_4_test",
		},
		{
			name: "XML types test for trace: hops",
			w: &MainWorkflow{
				Config: &Config{},
			},
			wantNMAPRun: NMAPRun{
				Scanner: "nmap",
				Version: "5.59BETA3",
				ScanInfo: ScanInfo{
					Services: "1-1000",
				},
				Host: []Host{
					{
						HostAddress: []HostAddress{
							{
								Address:     "10.10.10.20",
								AddressType: "ipv4",
							},
						},
						Trace: Trace{
							Port:     256,
							Protocol: "tcp",
							Hops: []Hop{
								{
									TTL:    1,
									IPAddr: "192.168.200.1",
									RTT:    0,
								},
								{
									TTL:    2,
									IPAddr: "192.168.1.1",
									RTT:    10.20,
								},
								{
									TTL:    3,
									IPAddr: "10.10.10.20",
									RTT:    23.30,
								},
							},
						},
					},
				},
			},
			wantErr: false,
			fileContent: `<?xml version="1.0"?>
			<nmaprun scanner="nmap" version="5.59BETA3">
				<scaninfo services="1-1000" />
				<host>
					<address addr="10.10.10.20" addrtype="ipv4" />
					<trace port="256" proto="tcp">
						<hop ttl="1" ipaddr="192.168.200.1" rtt="--" />
						<hop ttl="2" ipaddr="192.168.1.1" rtt="10.20" />
						<hop ttl="3" ipaddr="10.10.10.20" rtt="23.30" />
					</trace>
				</host>
			</nmaprun>`,
			fileName: "main_workflow_parse_5_test_hops",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fileName != "" {
				name := path.Join(os.TempDir(), tt.fileName)
				// Creating file with test-case content
				err := os.WriteFile(name, []byte(tt.fileContent), os.ModePerm)
				if err != nil {
					t.Errorf("Could not write file, error %v", err)
				}
				// deferring file removal after the test
				defer func() {
					_ = os.Remove(name)
				}()
				f, err := os.Open(name)
				if err != nil {
					t.Errorf("could not read source file: %v", err)
				}
				defer func() {
					_ = f.Close()
				}()
				tt.w.Config.InputFileConfig = InputFileConfig{
					Path:   name,
					Source: f,
				}
			}
			gotNMAPRun, err := tt.w.parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("MainWorkflow.parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNMAPRun, tt.wantNMAPRun) {
				t.Errorf("MainWorkflow.parse() = %+v, want %+v", gotNMAPRun, tt.wantNMAPRun)
			}
		})
	}
}

func TestMainWorkflow_Execute(t *testing.T) {
	tests := []struct {
		name        string
		w           *MainWorkflow
		wantErr     bool
		fileName    string
		fileContent string
		before      func(file string, t *testing.T)
	}{
		{
			name: "Parse of the file has failed",
			w: &MainWorkflow{
				Config: &Config{},
			},
			wantErr:     true,
			fileName:    "main_workflow_Execute_3_test",
			fileContent: "[NOT XML file]",
		},
		{
			name: "Empty CSV with header",
			w: &MainWorkflow{
				Config: &Config{
					OutputFormat: CSVOutput,
				},
			},
			wantErr:  false,
			fileName: "main_workflow_Execute_5_test",
			fileContent: `<?xml version="1.0"?>
			<nmaprun></nmaprun>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(tt.fileName, t)
			}
			if tt.fileName != "" {
				name := path.Join(os.TempDir(), tt.fileName)
				err := os.WriteFile(name, []byte(tt.fileContent), os.ModePerm)
				if err != nil {
					t.Errorf("Could not write file, error %v", err)
				}
				defer func() {
					_ = os.Remove(name)
				}()
				defer func() {
					_ = os.Remove(name + "_output")
				}()
				tt.w.Config.InputFileConfig = InputFileConfig{
					Path: name,
				}
				tt.w.Config.OutputFile = OutputFile(name + "_output")
			}
			tt.w.SetOutputFile()
			tt.w.SetInputFile()
			if err := tt.w.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("MainWorkflow.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
