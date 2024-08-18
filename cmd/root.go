package cmd

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

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/vdjagilev/nmap-formatter/v3/formatter"
)

var config = formatter.Config{
	OutputOptions: formatter.OutputOptions{
		HTMLOptions:         formatter.HTMLOutputOptions{},
		MarkdownOptions:     formatter.MarkdownOutputOptions{},
		JSONOptions:         formatter.JSONOutputOptions{},
		CSVOptions:          formatter.CSVOutputOptions{},
		SqliteOutputOptions: formatter.SqliteOutputOptions{},
		ExcelOptions:        formatter.ExcelOutputOptions{},
	},
	ShowVersion:       false,
	CurrentVersion:    VERSION,
	SkipDownHosts:     true,
	FilterExpressions: []string{},
}

// VERSION is describing current version of the nmap-formatter
const VERSION string = "3.0.2"

var workflow formatter.Workflow

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nmap-formatter [html|csv|md|json|dot|sqlite|excel|d2] [path-to-nmap.xml]",
	Short: "Utility that can help you to convert NMAP XML application output to various other formats",
	Long:  `This utility allows you to convert NMAP XML output to various other formats like (html, csv, markdown (md), json, dot, excel, sqlite, d2)`,
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
	rootCmd.Flags().StringArrayVar(&config.CustomOptions, "x-opts", []string{}, "--x-opts=\"some_key=some_value\"")

	// Use custom templates for HTML or Markdown output
	rootCmd.Flags().StringVar(&config.TemplatePath, "html-use-template", "", "--html-use-template /path/to/template.html")
	rootCmd.Flags().StringVar(&config.TemplatePath, "md-use-template", "", "--md-use-template /path/to/template.md")

	// Some options related to the output
	// Skip hosts that are down, so they won't be listed in the output
	// Skip header information (overall meta information from the scan)
	rootCmd.Flags().BoolVar(&config.OutputOptions.HTMLOptions.SkipHeader, "html-skip-header", false, "--html-skip-header, skips header in HTML output")
	rootCmd.Flags().BoolVar(&config.OutputOptions.MarkdownOptions.SkipHeader, "md-skip-header", false, "--md-skip-header, skips header in Markdown output")

	// Skip table of contents (TOC) information
	rootCmd.Flags().BoolVar(&config.OutputOptions.HTMLOptions.SkipTOC, "html-skip-toc", false, "--html-skip-toc, skips table of contents in HTML output")
	rootCmd.Flags().BoolVar(&config.OutputOptions.MarkdownOptions.SkipTOC, "md-skip-toc", false, "--md-skip-toc, skips table of contents in Markdown output")

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

	rootCmd.Flags().BoolVar(&config.OutputOptions.HTMLOptions.DarkMode, "html-dark-mode", true, "--html-dark-mode=false, sets HTML output in dark colours")

	rootCmd.Flags().BoolVar(&config.OutputOptions.HTMLOptions.FloatingContentsTable, "html-toc-float", false, "--html-toc-float=true, Table of contents floats along with the scroll")

	// Pretty-print json
	rootCmd.Flags().BoolVar(&config.OutputOptions.JSONOptions.PrettyPrint, "json-pretty", true, "--json-pretty=false (pretty prints JSON output)")

	// Configs related to SQLite
	rootCmd.Flags().StringVar(&config.OutputOptions.SqliteOutputOptions.DSN, "sqlite-dsn", "nmap.sqlite", "--sqlite-dsn nmap.sqlite")
	rootCmd.Flags().StringVar(&config.OutputOptions.SqliteOutputOptions.ScanIdentifier, "scan-id", "", "--scan-id abc123")

	// Configs related to D2 language
	rootCmd.Flags().BoolVar(&config.SkipDownHosts, "skip-down-hosts", false, "--skip-down-hosts=true, skips hosts that are offline")

	// Multiple filter expressions supported
	rootCmd.Flags().StringArrayVar(&config.FilterExpressions, "filter", []string{}, "--filter '.Status.State == \"up\" && any(.Port, { .PortID in [80,443] })'")

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
	workflow.SetInputFile()
	workflow.SetOutputFile()

	err = workflow.Execute()
	if err != nil {
		return err
	}

	if config.Writer != nil {
		config.Writer.Close()
	}
	if config.InputFileConfig.Source != nil {
		config.InputFileConfig.Source.Close()
	}

	return nil
}
