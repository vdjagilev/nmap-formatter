package formatter

import (
	"io"
	"reflect"
	"testing"
)

func TestConfig_CustomOptionsMap(t *testing.T) {
	type fields struct {
		Writer          io.Writer
		OutputFormat    OutputFormat
		InputFileConfig InputFileConfig
		OutputFile      OutputFile
		OutputOptions   OutputOptions
		ShowVersion     bool
		TemplatePath    string
		CustomOptions   []string
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		{
			name: "No options",
			fields: fields{
				CustomOptions: []string{},
			},
			want: map[string]string{},
		},
		{
			name: "Simple key-value pairs",
			fields: fields{
				CustomOptions: []string{
					"key=value",
					"key1=value1",
					"key2=value2",
				},
			},
			want: map[string]string{
				"key":  "value",
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "Testing first split case",
			fields: fields{
				CustomOptions: []string{
					"key=value",
					"key2=value2=2",
				},
			},
			want: map[string]string{
				"key":  "value",
				"key2": "value2=2",
			},
		},
		{
			name: "Testing more split `=` signs",
			fields: fields{
				CustomOptions: []string{
					"key=value",
					"key2=value2=2=2=2",
				},
			},
			want: map[string]string{
				"key":  "value",
				"key2": "value2=2=2=2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Writer:          tt.fields.Writer,
				OutputFormat:    tt.fields.OutputFormat,
				InputFileConfig: tt.fields.InputFileConfig,
				OutputFile:      tt.fields.OutputFile,
				OutputOptions:   tt.fields.OutputOptions,
				ShowVersion:     tt.fields.ShowVersion,
				TemplatePath:    tt.fields.TemplatePath,
				CustomOptions:   tt.fields.CustomOptions,
			}
			if got := c.CustomOptionsMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.CustomOptionsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
