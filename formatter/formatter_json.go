package formatter

import (
	"bytes"
	"encoding/json"
)

// JSONFormatter is struct defined for JSON Output use-case
type JSONFormatter struct {
	config *Config
}

// Format the data and output it to appropriate io.Writer
func (f *JSONFormatter) Format(td *TemplateData, templateContent string) (err error) {
	jsonData := new(bytes.Buffer)
	encoder := json.NewEncoder(jsonData)
	if td.OutputOptions.JSONOptions.PrettyPrint {
		// space size = 2
		encoder.SetIndent("", "  ")
	}
	err = encoder.Encode(td.NMAPRun)
	if err != nil {
		return err
	}
	_, err = f.config.Writer.Write(jsonData.Bytes())
	return
}

func (f *JSONFormatter) defaultTemplateContent() string {
	return ""
}
