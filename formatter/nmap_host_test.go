package formatter

import (
	"testing"
)

func TestHost_JoinedAddresses(t *testing.T) {
	type fields struct {
		StartTime     int
		EndTime       int
		Port          []Port
		HostAddress   []HostAddress
		HostNames     HostNames
		Status        HostStatus
		OS            OS
		Trace         Trace
		Uptime        Uptime
		Distance      Distance
		TCPSequence   TCPSequence
		IPIDSequence  IPIDSequence
		TCPTSSequence TCPTSSequence
	}
	type args struct {
		delimiter string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Single address",
			fields: fields{
				HostAddress: []HostAddress{
					{
						Address:     "192.168.1.1",
						AddressType: "ipv4",
					},
				},
			},
			args: args{
				delimiter: "/",
			},
			want: "192.168.1.1",
		},
		{
			name: "Two addresses",
			fields: fields{
				HostAddress: []HostAddress{
					{
						Address:     "192.168.1.1",
						AddressType: "ipv4",
					},
					{
						Address:     "FF:FF:FF:FF:FF",
						AddressType: "mac",
					},
				},
			},
			args: args{
				delimiter: "/",
			},
			want: "192.168.1.1 / FF:FF:FF:FF:FF",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Host{
				StartTime:     tt.fields.StartTime,
				EndTime:       tt.fields.EndTime,
				Port:          tt.fields.Port,
				HostAddress:   tt.fields.HostAddress,
				HostNames:     tt.fields.HostNames,
				Status:        tt.fields.Status,
				OS:            tt.fields.OS,
				Trace:         tt.fields.Trace,
				Uptime:        tt.fields.Uptime,
				Distance:      tt.fields.Distance,
				TCPSequence:   tt.fields.TCPSequence,
				IPIDSequence:  tt.fields.IPIDSequence,
				TCPTSSequence: tt.fields.TCPTSSequence,
			}
			if got := h.JoinedAddresses(tt.args.delimiter); got != tt.want {
				t.Errorf("Host.JoinedAddresses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHost_JoinedHostNames(t *testing.T) {
	type fields struct {
		StartTime     int
		EndTime       int
		Port          []Port
		HostAddress   []HostAddress
		HostNames     HostNames
		Status        HostStatus
		OS            OS
		Trace         Trace
		Uptime        Uptime
		Distance      Distance
		TCPSequence   TCPSequence
		IPIDSequence  IPIDSequence
		TCPTSSequence TCPTSSequence
	}
	type args struct {
		delimiter string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Empty string",
			fields: fields{
				StartTime:     0,
				EndTime:       0,
				Port:          []Port{},
				HostAddress:   []HostAddress{},
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
			args: args{
				delimiter: ", ",
			},
			want: "",
		},
		{
			name: "1 name",
			fields: fields{
				StartTime:   0,
				EndTime:     0,
				Port:        []Port{},
				HostAddress: []HostAddress{},
				HostNames: HostNames{
					HostName: []HostName{
						{
							Name: "example.com",
						},
					},
				},
				Status:        HostStatus{},
				OS:            OS{},
				Trace:         Trace{},
				Uptime:        Uptime{},
				Distance:      Distance{},
				TCPSequence:   TCPSequence{},
				IPIDSequence:  IPIDSequence{},
				TCPTSSequence: TCPTSSequence{},
			},
			args: args{
				delimiter: ", ",
			},
			want: "example.com",
		},
		{
			name: "2 names",
			fields: fields{
				StartTime:   0,
				EndTime:     0,
				Port:        []Port{},
				HostAddress: []HostAddress{},
				HostNames: HostNames{
					HostName: []HostName{
						{
							Name: "example.com",
						},
						{
							Name: "github.com",
						},
					},
				},
				Status:        HostStatus{},
				OS:            OS{},
				Trace:         Trace{},
				Uptime:        Uptime{},
				Distance:      Distance{},
				TCPSequence:   TCPSequence{},
				IPIDSequence:  IPIDSequence{},
				TCPTSSequence: TCPTSSequence{},
			},
			args: args{
				delimiter: ", ",
			},
			want: "example.com, github.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Host{
				StartTime:     tt.fields.StartTime,
				EndTime:       tt.fields.EndTime,
				Port:          tt.fields.Port,
				HostAddress:   tt.fields.HostAddress,
				HostNames:     tt.fields.HostNames,
				Status:        tt.fields.Status,
				OS:            tt.fields.OS,
				Trace:         tt.fields.Trace,
				Uptime:        tt.fields.Uptime,
				Distance:      tt.fields.Distance,
				TCPSequence:   tt.fields.TCPSequence,
				IPIDSequence:  tt.fields.IPIDSequence,
				TCPTSSequence: tt.fields.TCPTSSequence,
			}
			if got := h.JoinedHostNames(tt.args.delimiter); got != tt.want {
				t.Errorf("Host.JoinedHostNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHostStatus_IsUp(t *testing.T) {
	type fields struct {
		State  string
		Reason string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Host is down",
			fields: fields{
				State:  "down",
				Reason: "ping",
			},
			want: false,
		},
		{
			name: "Host is up",
			fields: fields{
				State:  "up",
				Reason: "no reason",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hs := &HostStatus{
				State:  tt.fields.State,
				Reason: tt.fields.Reason,
			}
			if got := hs.IsUp(); got != tt.want {
				t.Errorf("HostStatus.IsUp() = %v, want %v", got, tt.want)
			}
		})
	}
}
