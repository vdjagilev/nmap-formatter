package formatter

import (
	"bytes"
	"encoding/json"
)

type JSONFormatter struct {
	config *Config
}

// Format the data and output it to appropriate io.Writer
func (f *JSONFormatter) Format(td *TemplateData) (err error) {
	jsonData := new(bytes.Buffer)
	err = json.NewEncoder(jsonData).Encode(td.NMAPRun)
	if err != nil {
		return err
	}
	_, err = f.config.Writer.Write(jsonData.Bytes())
	return
}
