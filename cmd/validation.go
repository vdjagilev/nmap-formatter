package cmd

import (
	"fmt"
	"os"

	"github.com/vdjagilev/nmap-formatter/v2/formatter"
)

// validate is checking input from the command line
func validate(config formatter.Config) error {
	if !config.OutputFormat.IsValid() {
		return fmt.Errorf("not valid format: %s, please choose html/json/md/csv", config.OutputFormat)
	}

	err := validateIOFiles(config)
	if err != nil {
		return err
	}

	return validateTemplateConfig(config)
}

// validateIOFiles validates whether Input files and output files exists/have permissions to be created
func validateIOFiles(config formatter.Config) error {
	// Checking if xml file is readable
	if !config.InputFileConfig.IsStdin {
		err := config.InputFileConfig.ExistsOpen()
		if err != nil {
			return fmt.Errorf("could not open XML file: %v", err)
		}
	}
	// Checking if output file can be created and does not exist already
	// If OutputFile == "", it means that all output goes to stdout, no check needed
	if config.OutputFile != "" {
		outputFile, err := os.OpenFile(string(config.OutputFile), os.O_EXCL|os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			return fmt.Errorf("unable to create output file: %s", err)
		}
		outputFile.Close()
		os.Remove(string(config.OutputFile))
	}
	return nil
}

// validateTemplateConfig validates template config, if it has adequate output format configs and is readable
func validateTemplateConfig(config formatter.Config) error {
	// Checking if custom template is existing and readable and
	if config.TemplatePath != "" {
		switch config.OutputFormat {
		case formatter.CSVOutput:
		case formatter.JSONOutput:
			return fmt.Errorf("cannot set templates for the formats other than HTML or Markdown")
		}
		file, err := os.Open(config.TemplatePath)
		if err != nil {
			return fmt.Errorf("could not read template file: %v", err)
		}
		defer file.Close()
	}
	return nil
}
