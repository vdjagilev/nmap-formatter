package formatter

import (
	"os"
	"path"
	"testing"
)

func TestInputFileConfig_ExistsOpen(t *testing.T) {
	type fields struct {
		Path    string
		IsStdin bool
	}
	beforeFunc := func(path string, t *testing.T) {
		f, err := os.Create(path)
		if err != nil {
			t.Errorf("error creating temporary file: %s", err)
		}
		defer f.Close()
	}
	afterFunc := func(name string) {
		os.Remove(name)
	}
	tests := []struct {
		name      string
		fields    fields
		wantErr   bool
		file      string
		runBefore bool
		runAfter  bool
		before    func(path string, t *testing.T)
		after     func(path string)
	}{
		{
			name: "File does not exist",
			fields: fields{
				Path: "",
			},
			wantErr: true,
		},
		{
			name: "File exists",
			fields: fields{
				Path: path.Join(os.TempDir(), "inputfile_config_test_exists_2.txt"),
			},
			wantErr:   false,
			file:      path.Join(os.TempDir(), "inputfile_config_test_exists_2.txt"),
			before:    beforeFunc,
			after:     afterFunc,
			runBefore: true,
			runAfter:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.runBefore {
				tt.before(tt.file, t)
			}
			if tt.runAfter {
				defer tt.after(tt.file)
			}
			i := &InputFileConfig{
				Path:    tt.fields.Path,
				IsStdin: tt.fields.IsStdin,
			}
			if err := i.ExistsOpen(); (err != nil) != tt.wantErr {
				t.Errorf("InputFileConfig.ExistsOpen() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
