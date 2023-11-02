package formatter

import (
	"bufio"
	"errors"
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

// ReadContents reads content from stdin or provided file-path
func (i *InputFileConfig) ReadContents() ([]byte, error) {
	var err error
	var content []byte
	if i.Source == nil {
		return nil, errors.New("no reading source is defined")
	}
	scanner := bufio.NewScanner(i.Source)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		content = append(content, scanner.Bytes()...)
	}
	err = scanner.Err()
	return content, err
}

// ExistsOpen tries to open a file for reading, returning an error if it fails
func (i *InputFileConfig) ExistsOpen() error {
	f, err := os.Open(i.Path)
	f.Close()
	return err
}
