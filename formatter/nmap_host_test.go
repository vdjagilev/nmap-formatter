package formatter

import "testing"

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
