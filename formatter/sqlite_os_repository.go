package formatter

import (
	"database/sql"
	"strings"
)

type OSRepository struct {
	sqlite *SqliteDB
	conn   *sql.DB
	hostID int64
}

const insertHostOSSQL = `
	INSERT INTO host_os (
		host_id,
		class_type,
		class_vendor,
		class_osfamily,
		class_osgen,
		class_accuracy,
		cpe
	) VALUES (?, ?, ?, ?, ?, ?, ?)`

const insertHostOSPortUsedSQL = `
	INSERT INTO host_os_port_used (
		host_os_id,
		state,
		protocol,
		port_id
	) VALUES (?, ?, ?, ?)`

const insertHostOSMatchSQL = `
	INSERT INTO host_os_match (
		host_os_id,
		name,
		accuracy,
		line
	) VALUES (?, ?, ?, ?)`

func (o *OSRepository) insertOSRecords(os *OS) error {
	id, err := o.sqlite.insertReturnID(
		insertHostOSSQL,
		o.hostID,
		os.OSClass.Type,
		os.OSClass.Vendor,
		os.OSClass.OSFamily,
		os.OSClass.OSGen,
		os.OSClass.Accuracy,
		strings.Join(os.OSClass.CPE, sqliteStringDelimiter),
	)
	if err != nil {
		return err
	}
	err = o.insertOSPortUsed(id, os.OSPortUsed)
	if err != nil {
		return err
	}
	return o.insertOSMatch(id, os.OSMatch)
}

func (o *OSRepository) insertOSPortUsed(osID int64, osPortUsed []OSPortUsed) error {
	for _, port := range osPortUsed {
		err := o.sqlite.insert(
			insertHostOSPortUsedSQL,
			osID,
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

func (o *OSRepository) insertOSMatch(osID int64, osMatch []OSMatch) error {
	for _, match := range osMatch {
		err := o.sqlite.insert(
			insertHostOSMatchSQL,
			osID,
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
