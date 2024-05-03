package formatter

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

// Workflow interface that describes the main functions that are used in nmap-formatter
type Workflow interface {
	Execute() (err error)
	SetConfig(c *Config)
	SetInputFile()
	SetOutputFile()
}

// MainWorkflow is main workflow implementation struct
type MainWorkflow struct {
	Config *Config
}

// SetConfig is a simple setter-function that sets the configuration
func (w *MainWorkflow) SetConfig(c *Config) {
	w.Config = c
}

// SetOutputFile sets output file (file descriptor) depending on config and returns error
// if there is output file reading issue
func (w *MainWorkflow) SetOutputFile() {
	if w.Config.OutputFile == "" {
		w.Config.Writer = os.Stdout
	} else {
		// Error has been checked before executing this function
		f, _ := os.OpenFile(string(w.Config.OutputFile), os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.ModePerm)
		w.Config.Writer = f
	}
}

// SetInputFile sets an input file (file descriptor) in the config
func (w *MainWorkflow) SetInputFile() {
	var inputFile *os.File
	if w.Config.InputFileConfig.IsStdin {
		inputFile = os.Stdin
	} else {
		// Error has been checked before executing this function
		inputFile, _ = os.Open(w.Config.InputFileConfig.Path)
	}
	w.Config.InputFileConfig.Source = inputFile
}

// Execute is the core of the application which executes required steps
// one-by-one to achieve formatting from input -> output.
func (w *MainWorkflow) Execute() (err error) {
	// Reading & parsing the input file
	NMAPRun, err := w.parse()
	if err != nil {
		return
	}

	filteredRun := NMAPRun
	for _, expr := range w.Config.FilterExpressions {
		log.Printf("filtering with expression: %s", expr)
		filteredRun, err = filterExpr(filteredRun, expr)
		if err != nil {
			return fmt.Errorf("error filtering: %v", err)
		}
	}

	// Build template data with NMAPRun entry & various output options
	templateData := TemplateData{
		NMAPRun:       filteredRun,
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
	return formatter.Format(&templateData, templateContent)
}

// parse reads & unmarshalles the input file into NMAPRun struct
func (w *MainWorkflow) parse() (run NMAPRun, err error) {
	if w.Config.InputFileConfig.Source == nil {
		return run, fmt.Errorf("no input file is defined")
	}
	d := xml.NewDecoder(w.Config.InputFileConfig.Source)
	_, err = d.Token()
	if err != nil {
		return
	}

	err = d.Decode(&run)
	return
}
