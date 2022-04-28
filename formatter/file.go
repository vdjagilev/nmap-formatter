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

type InputFileConfig struct {
	Path    string
	IsStdin bool
	Source  io.Reader
}

// ReadContents reads content from stdin or provided file-path
func (i *InputFileConfig) ReadContents() ([]byte, error) {
	var err error
	var content []byte
	if i.Source == nil {
		return nil, errors.New("no reading source is defined")
	}
	scanner := bufio.NewScanner(i.Source)
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
