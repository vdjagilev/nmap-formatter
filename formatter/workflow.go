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
	var inputFile *os.File
	// This one is read in `parse()` function, we can close it here
	defer inputFile.Close()
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

	// Set InputFileConfig source to stdin or specific file
	if w.Config.InputFileConfig.IsStdin {
		inputFile = os.Stdin
	} else {
		inputFile, err = os.Open(w.Config.InputFileConfig.Path)
		if err != nil {
			return
		}
	}
	w.Config.InputFileConfig.Source = inputFile

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

	// Setting custom options for template if they exist
	if len(w.Config.CustomOptions) > 0 {
		templateData.CustomOptions = w.Config.CustomOptionsMap()
	}

	// Getting new instance of formatter based on provided config
	formatter := New(w.Config)

	// This part usually should not happen
	if formatter == nil {
		return fmt.Errorf("no formatter is defined")
	}

	// Trying to read template content (read a file, or get default in case where no option was used)
	templateContent, err := TemplateContent(formatter, w.Config)
	if err != nil {
		return fmt.Errorf("error getting template content: %v", err)
	}

	err = formatter.Format(&templateData, templateContent)
	return
}

// parse reads & unmarshalles the input file into NMAPRun struct
func (w *MainWorkflow) parse() (run NMAPRun, err error) {
	input, err := w.Config.InputFileConfig.ReadContents()
	if err != nil {
		return
	}
	if err = xml.Unmarshal(input, &run); err != nil {
		return
	}
	return run, nil
}
