package formatter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple CamelCase",
			input:    "CamelCase",
			expected: "camel_case",
		},
		{
			name:     "Already snake_case",
			input:    "snake_case",
			expected: "snake_case",
		},
		{
			name:     "Single uppercase",
			input:    "X",
			expected: "x",
		},
		{
			name:     "Single lowercase",
			input:    "x",
			expected: "x",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Multiple consecutive uppercase",
			input:    "HTTPServer",
			expected: "httpserver",
		},
		{
			name:     "Starting with lowercase",
			input:    "startLowerCase",
			expected: "start_lower_case",
		},
		{
			name:     "Mixed case with numbers",
			input:    "Port80Open",
			expected: "port80_open",
		},
		{
			name:     "All uppercase",
			input:    "XML",
			expected: "xml",
		},
		{
			name:     "NMAP field examples",
			input:    "StartStr",
			expected: "start_str",
		},
		{
			name:     "NMAP field examples 2",
			input:    "NumServices",
			expected: "num_services",
		},
		{
			name:     "NMAP field examples 3",
			input:    "TimeStr",
			expected: "time_str",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toSnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("toSnakeCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSnakeCaseEncoder(t *testing.T) {
	tests := []struct {
		name        string
		data        interface{}
		prettyPrint bool
		checkKeys   []string
	}{
		{
			name: "Simple struct",
			data: struct {
				FirstName string
				LastName  string
				Age       int
			}{
				FirstName: "John",
				LastName:  "Doe",
				Age:       30,
			},
			prettyPrint: false,
			checkKeys:   []string{"first_name", "last_name", "age"},
		},
		{
			name: "Nested struct",
			data: struct {
				HostName string
				Status   struct {
					IsUp   bool
					Reason string
				}
			}{
				HostName: "example.com",
				Status: struct {
					IsUp   bool
					Reason string
				}{
					IsUp:   true,
					Reason: "echo-reply",
				},
			},
			prettyPrint: true,
			checkKeys:   []string{"host_name", "status", "is_up", "reason"},
		},
		{
			name: "Array of structs",
			data: struct {
				Hosts []struct {
					IPAddr string
					PortID int
				}
			}{
				Hosts: []struct {
					IPAddr string
					PortID int
				}{
					{IPAddr: "192.168.1.1", PortID: 80},
					{IPAddr: "192.168.1.2", PortID: 443},
				},
			},
			prettyPrint: false,
			checkKeys:   []string{"hosts", "ipaddr", "port_id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			encoder := newSnakeCaseEncoder(buf, tt.prettyPrint)

			err := encoder.Encode(tt.data)
			if err != nil {
				t.Fatalf("Encode() error = %v", err)
			}

			result := buf.String()

			// Verify all expected keys are present in snake_case
			for _, key := range tt.checkKeys {
				expectedKey := `"` + key + `"`
				if !strings.Contains(result, expectedKey) {
					t.Errorf("Expected key %s not found in output:\n%s", expectedKey, result)
				}
			}

			// Verify the JSON is still valid
			var decoded interface{}
			if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
				t.Errorf("Output is not valid JSON: %v\nOutput: %s", err, result)
			}
		})
	}
}

func TestSnakeCaseEncoder_NMAPRun(t *testing.T) {
	// Test with actual NMAP structures
	nmapRun := NMAPRun{
		Scanner:  "nmap",
		Args:     "-sV",
		Start:    1234567890,
		StartStr: "2009-02-13 23:31:30 UTC",
		Version:  "7.80",
		ScanInfo: ScanInfo{
			Type:        "syn",
			Protocol:    "tcp",
			NumServices: 100,
			Services:    "1-100",
		},
		Host: []Host{
			{
				StartTime: 1234567890,
				EndTime:   1234567900,
				HostAddress: []HostAddress{
					{
						Address:     "192.168.1.1",
						AddressType: "ipv4",
					},
				},
				Status: HostStatus{
					State:  "up",
					Reason: "echo-reply",
				},
			},
		},
		Verbose: Verbose{
			Level: 1,
		},
		Debugging: Debugging{
			Level: 0,
		},
		RunStats: RunStats{
			Finished: Finished{
				Time:    1234567900,
				TimeStr: "2009-02-13 23:31:40 UTC",
				Elapsed: 10.5,
				Summary: "Scan completed",
				Exit:    "success",
			},
			Hosts: StatHosts{
				Up:    1,
				Down:  0,
				Total: 1,
			},
		},
	}

	tests := []struct {
		name        string
		prettyPrint bool
	}{
		{
			name:        "Compact output",
			prettyPrint: false,
		},
		{
			name:        "Pretty printed output",
			prettyPrint: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			encoder := newSnakeCaseEncoder(buf, tt.prettyPrint)

			err := encoder.Encode(&nmapRun)
			if err != nil {
				t.Fatalf("Encode() error = %v", err)
			}

			result := buf.String()

			// Verify key conversions
			expectedKeys := []string{
				`"scanner"`,
				`"args"`,
				`"start"`,
				`"start_str"`,
				`"version"`,
				`"scan_info"`,
				`"type"`,
				`"protocol"`,
				`"num_services"`,
				`"services"`,
				`"host"`,
				`"start_time"`,
				`"end_time"`,
				`"host_address"`,
				`"address"`,
				`"address_type"`,
				`"status"`,
				`"state"`,
				`"reason"`,
				`"verbose"`,
				`"level"`,
				`"debugging"`,
				`"run_stats"`,
				`"finished"`,
				`"time"`,
				`"time_str"`,
				`"elapsed"`,
				`"summary"`,
				`"exit"`,
				`"hosts"`,
				`"up"`,
				`"down"`,
				`"total"`,
			}

			for _, key := range expectedKeys {
				if !strings.Contains(result, key) {
					t.Errorf("Expected key %s not found in output", key)
				}
			}

			// Verify original CamelCase keys are NOT present
			unexpectedKeys := []string{
				`"Scanner"`,
				`"StartStr"`,
				`"ScanInfo"`,
				`"NumServices"`,
				`"HostAddress"`,
				`"AddressType"`,
				`"StartTime"`,
				`"EndTime"`,
				`"RunStats"`,
				`"TimeStr"`,
			}

			for _, key := range unexpectedKeys {
				if strings.Contains(result, key) {
					t.Errorf("Unexpected CamelCase key %s found in output", key)
				}
			}

			// Verify the JSON is still valid
			var decoded NMAPRun
			if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
				t.Errorf("Output is not valid JSON: %v", err)
			}
		})
	}
}

func TestJSONFormatter_FormatWithSnakeCase(t *testing.T) {
	tests := []struct {
		name        string
		snakeCase   bool
		prettyPrint bool
		checkKey    string
		shouldExist bool
	}{
		{
			name:        "Default (CamelCase)",
			snakeCase:   false,
			prettyPrint: false,
			checkKey:    `"StartStr"`,
			shouldExist: true,
		},
		{
			name:        "Snake case enabled",
			snakeCase:   true,
			prettyPrint: false,
			checkKey:    `"start_str"`,
			shouldExist: true,
		},
		{
			name:        "Snake case disabled should not have snake_case keys",
			snakeCase:   false,
			prettyPrint: false,
			checkKey:    `"start_str"`,
			shouldExist: false,
		},
		{
			name:        "Snake case enabled should not have CamelCase keys",
			snakeCase:   true,
			prettyPrint: false,
			checkKey:    `"StartStr"`,
			shouldExist: false,
		},
		{
			name:        "Snake case with pretty print",
			snakeCase:   true,
			prettyPrint: true,
			checkKey:    `"start_str"`,
			shouldExist: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &jsonMockedWriter{}
			formatter := &JSONFormatter{
				&Config{
					Writer: writer,
				},
			}

			nmapRun := NMAPRun{
				Scanner:  "nmap",
				StartStr: "2009-02-13 23:31:30 UTC",
			}

			td := &TemplateData{
				NMAPRun: nmapRun,
				OutputOptions: OutputOptions{
					JSONOptions: JSONOutputOptions{
						PrettyPrint: tt.prettyPrint,
						SnakeCase:   tt.snakeCase,
					},
				},
			}

			err := formatter.Format(td, "")
			if err != nil {
				t.Fatalf("Format() error = %v", err)
			}

			result := string(writer.data)
			keyExists := strings.Contains(result, tt.checkKey)

			if keyExists != tt.shouldExist {
				t.Errorf("Key %s existence = %v, want %v\nOutput: %s",
					tt.checkKey, keyExists, tt.shouldExist, result)
			}

			// Verify the output is valid JSON
			var decoded interface{}
			if err := json.Unmarshal(writer.data, &decoded); err != nil {
				t.Errorf("Output is not valid JSON: %v\nOutput: %s", err, result)
			}
		})
	}
}
