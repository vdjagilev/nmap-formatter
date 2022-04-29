package formatter

import (
	"io"
	"log"
	"strings"
)

// Config defines main application configs (requirements from user), like:
// where output will be delivered, desired output format, input file path, output file path
// and different output options
type Config struct {
	Writer          io.Writer
	OutputFormat    OutputFormat
	InputFileConfig InputFileConfig
	OutputFile      OutputFile
	OutputOptions   OutputOptions
	ShowVersion     bool
	TemplatePath    string
	CustomOptions   []string
}

// CustomOptionsMap returns custom options provided in the CLI
func (c *Config) CustomOptionsMap() map[string]string {
	m := map[string]string{}
	for _, o := range c.CustomOptions {
		split := strings.SplitN(o, "=", 2)
		if len(split) != 2 {
			log.Printf("custom option %s has wrong format", o)
			continue
		}
		// An attempt to replace all problematic characters with underscores
		r := strings.NewReplacer(
			"-", "_",
			" ", "_",
		)
		m[r.Replace(split[0])] = split[1]
	}
	return m
}
