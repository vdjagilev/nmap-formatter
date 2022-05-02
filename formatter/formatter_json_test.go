package formatter

import (
	"errors"
	"testing"
)

func TestJSONFormatter_Format(t *testing.T) {
	writer := jsonMockedWriter{}
	type args struct {
		td *TemplateData
	}
	tests := []struct {
		name       string
		f          *JSONFormatter
		args       args
		wantErr    bool
		err        error
		wantOutput string
	}{
		{
			name: "Empty output",
			f: &JSONFormatter{
				&Config{
					Writer: &writer,
				},
			},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{},
				},
			},
			wantErr:    false,
			err:        nil,
			wantOutput: "{\"Scanner\":\"\",\"Args\":\"\",\"Start\":0,\"StartStr\":\"\",\"Version\":\"\",\"ScanInfo\":{\"Type\":\"\",\"Protocol\":\"\",\"NumServices\":0,\"Services\":\"\"},\"Host\":null,\"Verbose\":{\"Level\":0},\"Debugging\":{\"Level\":0},\"RunStats\":{\"Finished\":{\"Time\":0,\"TimeStr\":\"\",\"Elapsed\":0,\"Summary\":\"\",\"Exit\":\"\"},\"Hosts\":{\"Up\":0,\"Down\":0,\"Total\":0}}}\n",
		},
		{
			name: "Empty output (with intend)",
			f: &JSONFormatter{
				&Config{
					Writer: &writer,
				},
			},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{},
					OutputOptions: OutputOptions{
						JSONOptions: JSONOutputOptions{
							PrettyPrint: true,
						},
					},
				},
			},
			wantErr: false,
			err:     nil,
			wantOutput: `{
  "Scanner": "",
  "Args": "",
  "Start": 0,
  "StartStr": "",
  "Version": "",
  "ScanInfo": {
    "Type": "",
    "Protocol": "",
    "NumServices": 0,
    "Services": ""
  },
  "Host": null,
  "Verbose": {
    "Level": 0
  },
  "Debugging": {
    "Level": 0
  },
  "RunStats": {
    "Finished": {
      "Time": 0,
      "TimeStr": "",
      "Elapsed": 0,
      "Summary": "",
      "Exit": ""
    },
    "Hosts": {
      "Up": 0,
      "Down": 0,
      "Total": 0
    }
  }
}
`,
		},
		{
			name: "Error",
			f: &JSONFormatter{
				&Config{
					Writer: &writer,
				},
			},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{},
				},
			},
			wantErr:    true,
			err:        errors.New("Some error happened"),
			wantOutput: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer.err = tt.err
			if err := tt.f.Format(tt.args.td, ""); (err != nil) != tt.wantErr {
				t.Errorf("JSONFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.err == nil && string(tt.wantOutput) != string(writer.data) {
				t.Errorf("JSONFormatter.Format output = \n%v, wantOutput = \n%v", string(writer.data), tt.wantOutput)
			}
		})
	}
}

type jsonMockedWriter struct {
	data []byte
	err  error
}

func (w *jsonMockedWriter) Write(p []byte) (n int, err error) {
	if w.err != nil {
		return 0, w.err
	}
	w.data = p
	return len(p), nil
}

func (w *jsonMockedWriter) Close() error {
	return nil
}
