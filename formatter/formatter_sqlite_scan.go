package formatter

import "github.com/google/uuid"

// getScanIdentifier returns a unique string provided either by a user or generates new random uuid
func (f *SqliteFormatter) getScanIdentifier() string {
	if f.config.OutputOptions.SqliteOutputOptions.ScanIdentifier == "" {
		return generateUUID()
	}
	return f.config.OutputOptions.SqliteOutputOptions.ScanIdentifier
}

// generateUUID generates new random uuid string
func generateUUID() string {
	return uuid.New().String()
}
