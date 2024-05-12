/*
Copyright © 2021-2024 vdjagilev

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	_ "embed"
	"errors"
	"os"
	"path"
	"testing"

	"github.com/spf13/cobra"
	"github.com/vdjagilev/nmap-formatter/v3/formatter"
)

func Test_validate(t *testing.T) {
	type args struct {
		config formatter.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		before  func(t *testing.T)
		after   func(t *testing.T)
	}{
		{
			name:    "Wrong output format",
			args:    args{config: formatter.Config{OutputFormat: formatter.OutputFormat("test")}},
			wantErr: true,
			before:  func(t *testing.T) {},
			after:   func(t *testing.T) {},
		},
		{
			name: "Missing input file",
			args: args{
				config: formatter.Config{
					OutputFormat: formatter.CSVOutput,
				},
			},
			wantErr: true,
			before:  func(t *testing.T) {},
			after:   func(t *testing.T) {},
		},
		{
			name: "Successful validation",
			args: args{
				config: formatter.Config{
					OutputFormat: formatter.CSVOutput,
					InputFileConfig: formatter.InputFileConfig{
						Path: path.Join(os.TempDir(), "formatter_cmd_valid_2"),
					},
				},
			},
			wantErr: false,
			before: func(t *testing.T) {
				path := path.Join(os.TempDir(), "formatter_cmd_valid_2")
				_, err := os.Create(path)
				if err != nil {
					t.Errorf("could not create temporary file: %s", path)
				}
			},
			after: func(t *testing.T) {
				path := path.Join(os.TempDir(), "formatter_cmd_valid_2")
				err := os.Remove(path)
				if err != nil {
					t.Logf("could not remove temporary file: %s", path)
				}
			},
		},
		{
			name: "Successful validation output file",
			args: args{
				config: formatter.Config{
					OutputFormat: formatter.ExcelOutput,
					InputFileConfig: formatter.InputFileConfig{
						Path: path.Join(os.TempDir(), "formatter_cmd_valid_output"),
					},
					OutputFile: formatter.OutputFile("output.xlsx"),
				},
			},
			wantErr: false,
			before: func(t *testing.T) {
				path := path.Join(os.TempDir(), "formatter_cmd_valid_output")
				_, err := os.Create(path)
				if err != nil {
					t.Errorf("could not create input file: %s", path)
				}
			},
			after: func(t *testing.T) {},
		},
		{
			name: "Successful validation template",
			args: args{
				config: formatter.Config{
					OutputFormat: formatter.MarkdownOutput,
					TemplatePath: path.Join(os.TempDir(), "formatter_template_valid"),
					InputFileConfig: formatter.InputFileConfig{
						Path: path.Join(os.TempDir(), "formatter_cmd_valid_template_3"),
					},
				},
			},
			wantErr: false,
			before: func(t *testing.T) {
				templatePath := path.Join(os.TempDir(), "formatter_template_valid")
				filePath := path.Join(os.TempDir(), "formatter_cmd_valid_template_3")
				_, err := os.Create(templatePath)
				if err != nil {
					t.Errorf("could not create temporary file: %s", templatePath)
				}
				_, err = os.Create(filePath)
				if err != nil {
					t.Errorf("could not create temporary file: %s", filePath)
				}
			},
			after: func(t *testing.T) {
				templatePath := path.Join(os.TempDir(), "formatter_template_valid")
				filePath := path.Join(os.TempDir(), "formatter_cmd_valid_template_3")
				err := os.Remove(templatePath)
				if err != nil {
					t.Logf("could not remove temporary file: %s", templatePath)
				}
				err = os.Remove(filePath)
				if err != nil {
					t.Logf("could not remove temporary file: %s", filePath)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before(t)
			defer tt.after(t)
			if err := validate(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_arguments(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "No XML path argument provided",
			args:    args{},
			wantErr: true,
		},
		{
			name: "No Output format argument provided (reading from stdin)",
			args: args{
				args: []string{"html"},
			},
			wantErr: false,
		},
		{
			name: "Version argument provided",
			args: args{
				args: []string{"version"},
			},
			wantErr: false,
		},
		{
			name: "2 arguments provided",
			args: args{
				args: []string{
					"file.xml",
					"html",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := arguments(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("arguments() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_run(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	before := func(file string, testWorkflow formatter.Workflow, testConfig formatter.Config, t *testing.T) {
		var err error
		if len(file) != 0 {
			_, err = os.Create(file)
		}
		if err != nil {
			t.Errorf("could not create temporary file: %s", file)
		}
		workflow = testWorkflow
		config = testConfig
		config.InputFileConfig = formatter.InputFileConfig{
			Path: file,
		}
	}
	after := func(file string, t *testing.T) {
		var err error
		if len(file) != 0 {
			err = os.Remove(file)
		}
		if err != nil {
			t.Logf("could not remove temporary file: %s", file)
		}
		workflow = nil
		config = formatter.Config{}
	}
	tests := []struct {
		name      string
		input     string
		workflow  formatter.Workflow
		config    formatter.Config
		args      args
		runBefore bool
		wantErr   bool
	}{
		{
			name: "Fails validation during the run (no settings at all, will fail)",
			args: args{},
			config: formatter.Config{
				ShowVersion: false,
			},
			wantErr: true,
		},
		{
			name:      "Workflow execution fails",
			input:     path.Join(os.TempDir(), "formatter_cmd_run_1"),
			runBefore: true,
			workflow: &testWorkflow{
				executeResult: errors.New("Bad failure"),
			},
			config: formatter.Config{
				OutputFormat: "csv",
				ShowVersion:  false,
			},
			args:    args{},
			wantErr: true,
		},
		{
			name:      "Shows version using flag",
			runBefore: true,
			config: formatter.Config{
				ShowVersion: true,
			},
			args:    args{},
			wantErr: false,
		},
		{
			name:      "Successful workflow execution",
			input:     path.Join(os.TempDir(), "formatter_cmd_run_2"),
			runBefore: true,
			workflow:  &testWorkflow{},
			config: formatter.Config{
				OutputFormat: "html",
				ShowVersion:  false, // false by default
				Writer:       &rootMockedWriter{},
				InputFileConfig: formatter.InputFileConfig{
					Source: os.Stdin,
				},
			},
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.runBefore {
				before(tt.input, tt.workflow, tt.config, t)
				defer after(tt.input, t)
			}
			if err := run(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type testWorkflow struct {
	executeResult error
}

func (w *testWorkflow) Execute() (err error) {
	return w.executeResult
}

func (w *testWorkflow) SetConfig(c *formatter.Config) {
}

func (w *testWorkflow) SetInputFile() {
}

func (w *testWorkflow) SetOutputFile() {
}

type rootMockedWriter struct {
	data []byte
}

func (w *rootMockedWriter) Write(p []byte) (n int, err error) {
	w.data = p
	return len(p), nil
}

func (w *rootMockedWriter) Close() error {
	return nil
}

func Test_shouldShowVersion(t *testing.T) {
	type args struct {
		c    *formatter.Config
		args []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Don't show version (html)",
			args: args{
				c: &formatter.Config{
					ShowVersion: false,
				},
				args: []string{
					"html",
					"path/to/file.xml",
				},
			},
			want: false,
		},
		{
			name: "Don't show version (arguments are used incorrectly)",
			args: args{
				c: &formatter.Config{
					ShowVersion: false,
				},
				args: []string{
					"version",
					"path/to/file.xml",
				},
			},
			want: false,
		},
		{
			name: "Show version (flag)",
			args: args{
				c: &formatter.Config{
					ShowVersion: true,
				},
				args: []string{},
			},
			want: true,
		},
		{
			name: "Show version (argument)",
			args: args{
				c: &formatter.Config{
					ShowVersion: false,
				},
				args: []string{
					"version",
				},
			},
			want: true,
		},
		{
			name: "Show version (both)",
			args: args{
				c: &formatter.Config{
					ShowVersion: true,
				},
				args: []string{
					"version",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldShowVersion(tt.args.c, tt.args.args); got != tt.want {
				t.Errorf("shouldShowVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
