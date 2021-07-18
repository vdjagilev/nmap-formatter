package formatter

import (
	"encoding/xml"
	"fmt"
	"os"
)

// Workflow interface that describes the main functions that are used in nmap-formatter
type Workflow interface {
	Execute() (err error)
	SetConfig(c *Config)
}

// MainWorkflow is main workflow implementation struct
type MainWorkflow struct {
	Config *Config
}

// SetConfig is a simple setter-function that sets the configuration
func (w *MainWorkflow) SetConfig(c *Config) {
	w.Config = c
}

// Execute is the core of the application which executes required steps
// one-by-one to achieve formatting from input -> output.
func (w *MainWorkflow) Execute() (err error) {
	// If no output file has been provided all content
	// goes to the STDOUT
	if w.Config.OutputFile == "" {
		w.Config.Writer = os.Stdout
	} else {
		// Open output file for writing, produces an error if file already exists
		// This won't work if user redirects output to some file using ">" or ">>"
		f, err := os.OpenFile(string(w.Config.OutputFile), os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return err
		}
		defer f.Close()
		w.Config.Writer = f
	}

	// Reading & parsing the input file
	NMAPRun, err := w.parse()
	if err != nil {
		return
	}

	// Build template data with NMAPRun entry & various output options
	templateData := TemplateData{
		NMAPRun:       NMAPRun,
		OutputOptions: w.Config.OutputOptions,
	}

	// Getting new instance of formatter based on provided config
	formatter := New(w.Config)

	// This part usually should not happen
	if formatter == nil {
		return fmt.Errorf("no formatter is defined")
	}

	err = formatter.Format(&templateData)
	return
}

// parse reads & unmarshalles the input file into NMAPRun struct
func (w *MainWorkflow) parse() (NMAPRun NMAPRun, err error) {
	input, err := os.ReadFile(string(w.Config.InputFile))
	if err != nil {
		return
	}
	if err = xml.Unmarshal(input, &NMAPRun); err != nil {
		return
	}
	return NMAPRun, nil
}
