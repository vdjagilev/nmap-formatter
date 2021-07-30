package formatter

import (
	"io"
)

// Config defines main application configs (requirements from user), like:
// where output will be delivered, desired output format, input file path, output file path
// and different output options
type Config struct {
	Writer        io.Writer
	OutputFormat  OutputFormat
	InputFile     InputFile
	OutputFile    OutputFile
	OutputOptions OutputOptions
	ShowVersion   bool
}
