package formatter

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/vdjagilev/nmap-formatter/types"
)

// New returns new instance of formatter the exact struct
// of formatter would depend on provided config
func New(config *Config) Formatter {
	switch config.OutputFormat {
	case types.JSONOutput:
		return &JSONFormatter{
			config,
		}
	case types.HTMLOutput:
		return &HTMLFormatter{
			config,
		}
	case types.MarkdownOutput:
		return &MarkdownFormatter{
			config,
		}
	case types.CSVOutput:
		return &CSVFormatter{
			config,
		}
	}
	return nil
}

type Workflow struct {
	Config *Config
}

// Execute is the core of the application which executes required steps
// one-by-one to achieve formatting from input -> output.
func (w *Workflow) Execute() (err error) {
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
	templateData := types.TemplateData{
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
func (w *Workflow) parse() (NMAPRun types.NMAPRun, err error) {
	input, err := os.ReadFile(string(w.Config.InputFile))
	if err != nil {
		return
	}
	if err = xml.Unmarshal(input, &NMAPRun); err != nil {
		return
	}
	return NMAPRun, nil
}

type Formatter interface {
	// Format the data and output it to appropriate io.Writer
	Format(td *types.TemplateData) error
}
