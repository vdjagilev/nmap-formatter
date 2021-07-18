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
			wantOutput: "{\"Scanner\":\"\",\"Args\":\"\",\"Start\":\"\",\"StartStr\":\"\",\"Version\":\"\",\"ScanInfo\":{\"Type\":\"\",\"Protocol\":\"\",\"NumServices\":\"\",\"Services\":\"\"},\"Host\":null,\"Verbose\":{\"Level\":\"\"},\"Debugging\":{\"Level\":\"\"},\"RunStats\":{\"Finished\":{\"Time\":\"\",\"TimeStr\":\"\",\"Elapsed\":\"\",\"Summary\":\"\",\"Exit\":\"\"},\"Hosts\":{\"Up\":\"\",\"Down\":\"\",\"Total\":\"\"}}}\n",
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
			if err := tt.f.Format(tt.args.td); (err != nil) != tt.wantErr {
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
