package formatter

import (
	"io"
	"os"
)

// OutputFile describes output file name (full path)
type OutputFile string

// InputFile describes input file (nmap XML full path)
type InputFile string

// InputFileConfig stores all options related to nmap XML (path, is content is taken from stdin and io reader)
type InputFileConfig struct {
	Path    string
	IsStdin bool
	Source  io.ReadCloser
}

// ExistsOpen tries to open a file for reading, returning an error if it fails
func (i *InputFileConfig) ExistsOpen() error {
	f, err := os.Open(i.Path)
	if err != nil {
		return err
	}
	return f.Close()
}
