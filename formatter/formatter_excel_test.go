package formatter

import "testing"

type excelMockedWriter struct {
	data []byte
}

func (w *excelMockedWriter) Write(p []byte) (n int, err error) {
	w.data = p
	return len(p), nil
}

func (w *excelMockedWriter) Close() error {
	return nil
}

func TestExcelFormatter_Format(t *testing.T) {
	writer := &excelMockedWriter{}
	type fields struct {
		config *Config
	}
	type args struct {
		td *TemplateData
	}
	tests := []struct {
		name    string
		f       *ExcelFormatter
		args    args
		wantErr bool
	}{
		{
			name: "Empty Test ExcelFormatter.Format",
			f: &ExcelFormatter{
				config: &Config{
					Writer: writer,
				},
			},
			args: args{
				td: &TemplateData{
					NMAPRun: NMAPRun{
						Host: []Host{},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.Format(tt.args.td, ""); (err != nil) != tt.wantErr {
				t.Errorf("ExcelFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
