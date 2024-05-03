package formatter

import (
	"reflect"
	"testing"
)

func Test_filterExpr(t *testing.T) {
	type args struct {
		nmapRUN NMAPRun
		code    string
	}
	tests := []struct {
		name    string
		args    args
		want    NMAPRun
		wantErr bool
	}{
		{
			name: "Basic test",
			args: args{
				nmapRUN: NMAPRun{
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
									Protocol: "http",
									PortID:   80,
									State:    PortState{},
									Service:  PortService{},
									Script:   []Script{},
								},
							},
							HostAddress: []HostAddress{
								{
									Address:     "10.10.10.1",
									AddressType: "ipv4",
									Vendor:      "",
								},
							},
							HostNames: HostNames{},
							Status: HostStatus{
								State:  "up",
								Reason: "",
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
									Protocol: "",
									PortID:   22,
									State:    PortState{},
									Service:  PortService{},
									Script:   []Script{},
								},
							},
							HostAddress: []HostAddress{
								{
									Address:     "10.10.10.20",
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
				code: `.Status.State == "up" && any(.Port, { .PortID in [80] })`,
			},
			want: NMAPRun{
				Scanner:  "",
				Args:     "",
				Start:    0,
				StartStr: "",
				Version:  "",
				ScanInfo: ScanInfo{},
				Host: []Host{{StartTime: 0, EndTime: 0, Port: []Port{
					{
						Protocol: "http",
						PortID:   80,
						State:    PortState{},
						Service:  PortService{},
						Script:   []Script{},
					},
				}, HostAddress: []HostAddress{{Address: "10.10.10.1", AddressType: "ipv4", Vendor: ""}}, HostNames: HostNames{}, Status: HostStatus{State: "up", Reason: ""}, OS: OS{}, Trace: Trace{}, Uptime: Uptime{}, Distance: Distance{}, TCPSequence: TCPSequence{}, IPIDSequence: IPIDSequence{}, TCPTSSequence: TCPTSSequence{}}},
				Verbose:   Verbose{},
				Debugging: Debugging{},
				RunStats:  RunStats{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filterExpr(tt.args.nmapRUN, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("filterExpr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterExpr() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
