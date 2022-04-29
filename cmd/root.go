package cmd

/*
Copyright © 2021-2022 vdjagilev

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

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/vdjagilev/nmap-formatter/formatter"
)

var config = formatter.Config{
	OutputOptions: formatter.OutputOptions{
		HTMLOptions:     formatter.HTMLOutputOptions{},
		MarkdownOptions: formatter.MarkdownOutputOptions{},
		JSONOptions:     formatter.JSONOutputOptions{},
		CSVOptions:      formatter.CSVOutputOptions{},
	},
	ShowVersion: false,
}

// VERSION is describing current version of the nmap-formatter
const VERSION string = "0.3.2"

var workflow formatter.Workflow

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nmap-formatter [html|csv|md|json] [path-to-nmap.xml]",
	Short: "Utility that can help you to convert NMAP XML application output to various other formats",
	Long:  `This utility allows you to convert NMAP XML output to various other formats like (html, csv, markdown (md), json)`,
	Args:  arguments,
	RunE:  run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Logging entries go to stderr, while stdout still can be used to save output to the file
	log.SetOutput(os.Stderr)

	rootCmd.Flags().StringVarP((*string)(&config.OutputFile), "file", "f", "", "-f output-file (by default \"\" will output to STDOUT)")
	rootCmd.Flags().BoolVar(&config.ShowVersion, "version", false, "--version, will show you the current version of the app")

	// Use custom templates for HTML or Markdown output
	rootCmd.Flags().StringVar(&config.TemplatePath, "html-use-template", "", "--html-use-template /path/to/template.html")
	rootCmd.Flags().StringVar(&config.TemplatePath, "md-use-template", "", "--md-use-template /path/to/template.md")

	// Some options related to the output
	// Skip hosts that are down, so they won't be listed in the output
	rootCmd.Flags().BoolVar(&config.OutputOptions.HTMLOptions.SkipDownHosts, "html-skip-down-hosts", true, "--html-skip-down-hosts=false, would print all hosts that are offline in HTML output")
	rootCmd.Flags().BoolVar(&config.OutputOptions.MarkdownOptions.SkipDownHosts, "md-skip-down-hosts", true, "--md-skip-down-hosts=false, would print all hosts that are offline in Markdown output")
	rootCmd.Flags().BoolVar(&config.OutputOptions.CSVOptions.SkipDownHosts, "csv-skip-down-hosts", true, "--csv-skip-down-hosts=false, would print all hosts that are offline in CSV output")

	// Skip summary (overall meta information from the scan)
	rootCmd.Flags().BoolVar(&config.OutputOptions.HTMLOptions.SkipSummary, "html-skip-summary", false, "--html-skip-summary=true, skips summary in HTML output")
	rootCmd.Flags().BoolVar(&config.OutputOptions.MarkdownOptions.SkipSummary, "md-skip-summary", false, "--md-skip-summary=true, skips summary in Markdown output")

	// Skip traceroute information (from scan machine to the target)
	rootCmd.Flags().BoolVar(&config.OutputOptions.HTMLOptions.SkipTraceroute, "html-skip-traceroute", false, "--html-skip-traceroute=true, skips traceroute information in HTML output")
	rootCmd.Flags().BoolVar(&config.OutputOptions.MarkdownOptions.SkipTraceroute, "md-skip-traceroute", false, "--md-skip-traceroute=true, skips traceroute information in Markdown output")

	// Skip metrics related information
	rootCmd.Flags().BoolVar(&config.OutputOptions.HTMLOptions.SkipMetrics, "html-skip-metrics", false, "--html-skip-metrics=true, skips metrics information in HTML output")
	rootCmd.Flags().BoolVar(&config.OutputOptions.MarkdownOptions.SkipMetrics, "md-skip-metrics", false, "--md-skip-metrics=true, skips metrics information in Markdown output")

	// Skip information from port scripts (nse-scripts)
	rootCmd.Flags().BoolVar(&config.OutputOptions.HTMLOptions.SkipPortScripts, "html-skip-port-scripts", false, "--html-skip-port-scripts=true, skips port scripts information in HTML output")
	rootCmd.Flags().BoolVar(&config.OutputOptions.MarkdownOptions.SkipPortScripts, "md-skip-port-scripts", false, "--md-skip-port-scripts=true, skips port scripts information in Markdown output")

	// Pretty-print json
	rootCmd.Flags().BoolVar(&config.OutputOptions.JSONOptions.PrettyPrint, "json-pretty", true, "--json-pretty=false (pretty prints JSON output)")

	workflow = &formatter.MainWorkflow{}
}

// arguments function validates the arguments passed to the application
// and sets configurations
func arguments(cmd *cobra.Command, args []string) error {
	if shouldShowVersion(&config, args) {
		return nil
	}
	if len(args) < 1 {
		return errors.New("requires output format argument")
	}

	config.OutputFormat = formatter.OutputFormat(args[0])
	config.InputFileConfig = formatter.InputFileConfig{}

	if len(args) > 1 {
		config.InputFileConfig.Path = args[1]
	} else {
		config.InputFileConfig.IsStdin = true
	}
	return nil
}

// version just prints the current version of nmap-formatter
func version() {
	fmt.Printf("nmap-formatter version: %s\n", VERSION)
}

// shouldShowVersion returns boolean whether app should show current version or not
// based on flag passed to the app or first argument
func shouldShowVersion(c *formatter.Config, args []string) bool {
	return c.ShowVersion || (len(args) == 1 && args[0] == "version")
}

// run executes the main application workflow and finishes fatally if there is some error
func run(cmd *cobra.Command, args []string) error {
	if shouldShowVersion(&config, args) {
		version()
		return nil
	}
	err := validate(config)
	if err != nil {
		return err
	}

	workflow.SetConfig(&config)

	err = workflow.Execute()
	if err != nil {
		return err
	}
	return nil
}

// validate is checking input from the command line
func validate(config formatter.Config) error {
	if !config.OutputFormat.IsValid() {
		return fmt.Errorf("not valid format: %s, please choose html/json/md/csv", config.OutputFormat)
	}

	// Checking if xml file is readable
	if !config.InputFileConfig.IsStdin {
		err := config.InputFileConfig.ExistsOpen()
		if err != nil {
			return fmt.Errorf("could not open XML file: %v", err)
		}
	}

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
