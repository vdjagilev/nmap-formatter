package formatter

import (
	"io"

	"github.com/vdjagilev/nmap-formatter/types"
)

type Config struct {
	Writer        io.Writer
	OutputFormat  types.OutputFormat
	InputFile     types.InputFile
	OutputFile    types.OutputFile
	OutputOptions types.OutputOptions
}
