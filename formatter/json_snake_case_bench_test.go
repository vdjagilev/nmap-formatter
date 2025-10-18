package formatter

import (
	"bytes"
	"testing"
)

// closeableBuffer wraps bytes.Buffer to implement io.WriteCloser
type closeableBuffer struct {
	*bytes.Buffer
}

func (cb *closeableBuffer) Close() error {
	return nil
}

// Create a realistic NMAP run for benchmarking
func createBenchmarkNMAPRun() *NMAPRun {
	return &NMAPRun{
		Scanner:  "nmap",
		Args:     "-sV -sC -p- -T4",
		Start:    1234567890,
		StartStr: "2009-02-13 23:31:30 UTC",
		Version:  "7.80",
		ScanInfo: ScanInfo{
			Type:        "syn",
			Protocol:    "tcp",
			NumServices: 1000,
			Services:    "1-1000",
		},
		Host: []Host{
			{
				StartTime: 1234567890,
				EndTime:   1234567900,
				HostAddress: []HostAddress{
					{Address: "192.168.1.1", AddressType: "ipv4"},
					{Address: "00:11:22:33:44:55", AddressType: "mac", Vendor: "Vendor Inc"},
				},
				HostNames: HostNames{
					HostName: []HostName{
						{Name: "example.com", Type: "PTR"},
					},
				},
				Status: HostStatus{State: "up", Reason: "echo-reply"},
				Port: []Port{
					{PortID: 22, Protocol: "tcp", State: PortState{State: "open", Reason: "syn-ack"}, Service: PortService{Name: "ssh", Product: "OpenSSH", Version: "8.0"}},
					{PortID: 80, Protocol: "tcp", State: PortState{State: "open", Reason: "syn-ack"}, Service: PortService{Name: "http", Product: "nginx", Version: "1.18"}},
					{PortID: 443, Protocol: "tcp", State: PortState{State: "open", Reason: "syn-ack"}, Service: PortService{Name: "https", Product: "nginx", Version: "1.18"}},
				},
				OS: OS{
					OSMatch: []OSMatch{
						{Name: "Linux 4.15 - 5.6", Accuracy: "95"},
					},
				},
			},
			{
				StartTime: 1234567890,
				EndTime:   1234567905,
				HostAddress: []HostAddress{
					{Address: "192.168.1.2", AddressType: "ipv4"},
				},
				Status: HostStatus{State: "up", Reason: "echo-reply"},
				Port: []Port{
					{PortID: 3306, Protocol: "tcp", State: PortState{State: "open", Reason: "syn-ack"}, Service: PortService{Name: "mysql", Product: "MySQL", Version: "5.7.30"}},
					{PortID: 8080, Protocol: "tcp", State: PortState{State: "open", Reason: "syn-ack"}, Service: PortService{Name: "http-proxy", Product: "Squid", Version: "4.10"}},
				},
			},
		},
		Verbose:   Verbose{Level: 1},
		Debugging: Debugging{Level: 0},
		RunStats: RunStats{
			Finished: Finished{
				Time:    1234567900,
				TimeStr: "2009-02-13 23:31:40 UTC",
				Elapsed: 10.5,
				Summary: "Nmap done: 2 IP addresses (2 hosts up) scanned in 10.50 seconds",
				Exit:    "success",
			},
			Hosts: StatHosts{Up: 2, Down: 0, Total: 2},
		},
	}
}

func BenchmarkJSONFormatter_Default(b *testing.B) {
	nmapRun := createBenchmarkNMAPRun()
	td := &TemplateData{
		NMAPRun: *nmapRun,
		OutputOptions: OutputOptions{
			JSONOptions: JSONOutputOptions{
				PrettyPrint: false,
				SnakeCase:   false,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := &closeableBuffer{Buffer: new(bytes.Buffer)}
		formatter := &JSONFormatter{
			&Config{Writer: buf},
		}
		if err := formatter.Format(td, ""); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONFormatter_SnakeCase(b *testing.B) {
	nmapRun := createBenchmarkNMAPRun()
	td := &TemplateData{
		NMAPRun: *nmapRun,
		OutputOptions: OutputOptions{
			JSONOptions: JSONOutputOptions{
				PrettyPrint: false,
				SnakeCase:   true,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := &closeableBuffer{Buffer: new(bytes.Buffer)}
		formatter := &JSONFormatter{
			&Config{Writer: buf},
		}
		if err := formatter.Format(td, ""); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONFormatter_PrettyPrint(b *testing.B) {
	nmapRun := createBenchmarkNMAPRun()
	td := &TemplateData{
		NMAPRun: *nmapRun,
		OutputOptions: OutputOptions{
			JSONOptions: JSONOutputOptions{
				PrettyPrint: true,
				SnakeCase:   false,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := &closeableBuffer{Buffer: new(bytes.Buffer)}
		formatter := &JSONFormatter{
			&Config{Writer: buf},
		}
		if err := formatter.Format(td, ""); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONFormatter_SnakeCasePretty(b *testing.B) {
	nmapRun := createBenchmarkNMAPRun()
	td := &TemplateData{
		NMAPRun: *nmapRun,
		OutputOptions: OutputOptions{
			JSONOptions: JSONOutputOptions{
				PrettyPrint: true,
				SnakeCase:   true,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := &closeableBuffer{Buffer: new(bytes.Buffer)}
		formatter := &JSONFormatter{
			&Config{Writer: buf},
		}
		if err := formatter.Format(td, ""); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkToSnakeCase(b *testing.B) {
	testCases := []string{
		"CamelCase",
		"StartStr",
		"NumServices",
		"HTTPServer",
		"HostAddress",
		"AddressType",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			_ = toSnakeCase(tc)
		}
	}
}

func BenchmarkConvertKeysToSnakeCase(b *testing.B) {
	// Sample JSON data
	jsonData := []byte(`{"Scanner":"nmap","Args":"-sV","Start":1234567890,"StartStr":"2009-02-13 23:31:30 UTC","Version":"7.80","ScanInfo":{"Type":"syn","Protocol":"tcp","NumServices":1000},"Host":[{"StartTime":1234567890,"EndTime":1234567900,"HostAddress":[{"Address":"192.168.1.1","AddressType":"ipv4"}],"Status":{"State":"up","Reason":"echo-reply"}}]}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = convertKeysToSnakeCase(jsonData)
	}
}
