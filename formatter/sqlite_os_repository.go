package formatter

import (
	"database/sql"
	"strings"
)

// OSRepository is responsible for populating database with OS related data on
// the scanned host
type OSRepository struct {
	sqlite *SqliteDB
	conn   *sql.DB
	hostID int64
}

const insertHostOSClassSQL = `
	INSERT INTO host_os_class (
		host_id,
		type,
		vendor,
		osfamily,
		osgen,
		accuracy,
		cpe
	) VALUES (?, ?, ?, ?, ?, ?, ?)`

const insertHostOSPortUsedSQL = `
	INSERT INTO host_os_port_used (
		host_id,
		state,
		protocol,
		port_id
	) VALUES (?, ?, ?, ?)`

const insertHostOSMatchSQL = `
	INSERT INTO host_os_match (
		host_id,
		name,
		accuracy,
		line
	) VALUES (?, ?, ?, ?)`

func (o *OSRepository) insertOSRecords(os *OS) error {
	err := o.insertOSClass(os.OSClass)
	if err != nil {
		return err
	}
	err = o.insertOSPortUsed(os.OSPortUsed)
	if err != nil {
		return err
	}
	return o.insertOSMatch(os.OSMatch)
}

func (o *OSRepository) insertOSClass(osClass []OSClass) error {
	for _, class := range osClass {
		err := o.sqlite.insert(
			insertHostOSClassSQL,
			o.hostID,
			class.Type,
			class.Vendor,
			class.OSFamily,
			class.OSGen,
			class.Accuracy,
			strings.Join(class.CPE, sqliteStringDelimiter),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *OSRepository) insertOSPortUsed(osPortUsed []OSPortUsed) error {
	for _, port := range osPortUsed {
		err := o.sqlite.insert(
			insertHostOSPortUsedSQL,
			o.hostID,
			port.State,
			port.Protocol,
			port.PortID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *OSRepository) insertOSMatch(osMatch []OSMatch) error {
	for _, match := range osMatch {
		err := o.sqlite.insert(
			insertHostOSMatchSQL,
			o.hostID,
			match.Name,
			match.Accuracy,
			match.Line,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
