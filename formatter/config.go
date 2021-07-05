package formatter

import (
	"io"
)

type Config struct {
	Writer        io.Writer
	OutputFormat  OutputFormat
	InputFile     InputFile
	OutputFile    OutputFile
	OutputOptions OutputOptions
}
