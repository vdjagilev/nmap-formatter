package formatter

// TemplateData is a struct that is used in the template, where
// NMAPRun is containing all parsed data from XML & OutputOptions contain all the information
// about certain output preferences.
type TemplateData struct {
	NMAPRun       NMAPRun
	OutputOptions OutputOptions
}
