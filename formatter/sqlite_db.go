package formatter

import (
	"database/sql"
	"fmt"

	// a package needed to embed files into a go runtime
	_ "embed"
)

// SqliteDB holds sqlite connection and transaction and performs
// main database function - preparation, population and finishing
// data migration
type SqliteDB struct {
	db             *sql.DB
	tx             *sql.Tx
	config         *Config
	scanRepository *ScanRepository
	stmtCache      map[string]*sql.Stmt // Cache for prepared statements
}

// SqliteDDL contains database schema definition
//
//go:embed resources/sql/sqlite_ddl.sql
var SqliteDDL string

// NewSqliteDB attempts to create new instance of SqliteDB struct
// and attempts to make a connection to the database, once it's successful
// it assigns variables to a ScanRepository struct and returns a pointer
func NewSqliteDB(c *Config) (*SqliteDB, error) {
	db, err := sql.Open("sqlite3", c.OutputOptions.SqliteOutputOptions.DSN)
	if err != nil {
		return nil, err
	}

	// Apply performance optimizations via PRAGMA settings
	// These settings significantly improve write performance for bulk inserts
	pragmas := []string{
		"PRAGMA synchronous = NORMAL",    // Balance between safety and speed
		"PRAGMA journal_mode = WAL",      // Write-Ahead Logging for better concurrency
		"PRAGMA cache_size = -64000",     // Use 64MB cache (negative = KB)
		"PRAGMA temp_store = MEMORY",     // Store temp tables/indices in memory
		"PRAGMA busy_timeout = 5000",     // Wait up to 5s if database is locked
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("failed to set %s: %v", pragma, err)
		}
	}

	sqlite := &SqliteDB{
		db:        db,
		tx:        nil,
		config:    c,
		stmtCache: make(map[string]*sql.Stmt),
		scanRepository: &ScanRepository{
			conn:   db,
			config: c,
		},
	}
	sqlite.scanRepository.sqlite = sqlite
	return sqlite, nil
}

func (s *SqliteDB) prepare() error {
	if !s.schemaExists() {
		err := s.generateSchema()
		if err != nil {
			return fmt.Errorf("could not generate schema: %v", err)
		}
	}
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction: %v", err)
	}
	s.tx = tx
	return nil
}

// schemaExists queries database and tries to select data from nf_schema
// database table, if that fails it's a clear indicator that schema does not exist in a database
func (s *SqliteDB) schemaExists() bool {
	// Try to get data from the database
	rows, err := s.db.Query(`SELECT version FROM nf_schema LIMIT 1`)
	if err != nil {
		return false
	}
	defer func() {
		_ = rows.Close()
	}()
	return true
}

// generateSchema creates
func (s *SqliteDB) generateSchema() error {
	// Create schema from SQL
	_, err := s.db.Exec(SqliteDDL)
	if err != nil {
		return err
	}

	// Set schema version by truncating table and inserting new version
	_, err = s.db.Exec(`DELETE FROM nf_schema;`)
	if err != nil {
		return fmt.Errorf("could not clean nf_schema table: %v", err)
	}

	_, err = s.db.Exec(`INSERT INTO nf_schema VALUES (?);`, s.config.CurrentVersion)
	if err != nil {
		return fmt.Errorf("could not insert new nf_schema version: %v", err)
	}
	return nil
}

// populate function starts populating database with scan results
func (s *SqliteDB) populate(n *NMAPRun) error {
	return s.scanRepository.insertScan(n)
}

// finish is a place where commit or rollback should happen and database connection is closed
func (s *SqliteDB) finish(err error) error {
	// Close all cached prepared statements before closing DB
	defer s.closeAllStmts()
	defer func() {
		_ = s.db.Close()
	}()
	if err != nil {
		rollbackErr := s.tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("failed rollback: %v: failed: %v", rollbackErr, err)
		}
		return err
	}

	err = s.tx.Commit()
	if err != nil {
		return fmt.Errorf("failed commit: %v", err)
	}
	return err
}

// getOrPrepareStmt retrieves a prepared statement from cache or prepares it if not cached
func (s *SqliteDB) getOrPrepareStmt(sql string) (*sql.Stmt, error) {
	// Check if statement is already cached
	if stmt, exists := s.stmtCache[sql]; exists {
		return stmt, nil
	}

	// Prepare new statement and cache it
	stmt, err := s.db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	s.stmtCache[sql] = stmt
	return stmt, nil
}

// closeAllStmts closes all cached prepared statements
func (s *SqliteDB) closeAllStmts() {
	for _, stmt := range s.stmtCache {
		_ = stmt.Close()
	}
	s.stmtCache = make(map[string]*sql.Stmt)
}

// insertReturnID is a generic function to execute INSERT SQL statement with arguments and return
// ID of the last element inserted and error in case it fails
func (s *SqliteDB) insertReturnID(sql string, args ...any) (int64, error) {
	stmt, err := s.getOrPrepareStmt(sql)
	if err != nil {
		return 0, err
	}
	result, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// insert is a generic function to execute INSERT SQL statement and return error
// if it fails
func (s *SqliteDB) insert(sql string, args ...any) error {
	stmt, err := s.getOrPrepareStmt(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(args...)
	return err
}
