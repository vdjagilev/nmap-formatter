package formatter

import (
	// Written this way to avoid automatic removal by text editor
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// sqliteStringDelimiter is a hardcoded string that is used to join
// certain values together in database
const sqliteStringDelimiter = ", "

// SqliteFormatter is a main struct to handle output for Sqlite
type SqliteFormatter struct {
	config *Config
}

// Format the data to sqlite format and insert it into appropriate output (Sqlite DSN)
func (f *SqliteFormatter) Format(td *TemplateData, templateContent string) error {
	db, err := NewSqliteDB(f.config)
	if err != nil {
		return fmt.Errorf("could not create new sqlite instance: %v", err)
	}

	err = db.prepare()
	if err != nil {
		return fmt.Errorf("failed to prepare db: %v", err)
	}

	err = db.populate(&td.NMAPRun)
	return db.finish(err)
}

// defaultTemplateContent does not return anything
func (f *SqliteFormatter) defaultTemplateContent() string {
	return ""
}
