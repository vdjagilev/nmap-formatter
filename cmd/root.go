/*
Copyright © 2021 vdjagilev

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
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/vdjagilev/nmap-formatter/formatter"
)

var config = formatter.Config{
	OutputOptions: formatter.OutputOptions{},
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nmap-formatter [path-to-nmap.xml] [html|csv|md|json]",
	Short: "Utility that can help you to convert NMAP XML application output to various other formats",
	Long:  `This utility allows you to convert NMAP XML output to various other formats like (html, csv, markdown (md), json)`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires an xml file argument")
		}
		if len(args) < 2 {
			return errors.New("requires output format argument")
		}

		config.InputFile = formatter.InputFile(args[0])
		config.OutputFormat = formatter.OutputFormat(args[1])

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := validate(config)
		if err != nil {
			log.Fatal(err)
		}

		workflow := formatter.Workflow{
			Config: &config,
		}

		err = workflow.Execute()
		if err != nil {
			log.Fatalf("failed to run: %v", err)
		}
	},
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
	rootCmd.Flags().BoolVar(&config.OutputOptions.DisplayDownHosts, "down-hosts", false, "--down-hosts=true")
}

// validate is checking input from the command line
func validate(config formatter.Config) error {
	// Validating the format
	if !config.OutputFormat.IsValid() {
		return fmt.Errorf("not valid format: %s, please choose html/json/md/csv", config.OutputFormat)
	}

	// Checking if xml file is readable
	f, err := os.Open(string(config.InputFile))
	if err != nil {
		return fmt.Errorf("could not open XML file: %v", err)
	}
	defer f.Close()
	return nil
}
