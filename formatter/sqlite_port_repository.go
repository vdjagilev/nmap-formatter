package formatter

import (
	"database/sql"
	"strings"
)

// PortRepository is responsible for populating database with a data
// related to ports
type PortRepository struct {
	sqlite *SqliteDB
	conn   *sql.DB
	hostID int64
}

const insertPortsSQL = `
	INSERT INTO ports (
		host_id,
		port_id,
		state_state,
		state_reason,
		state_reason_ttl,
		service_name,
		service_product,
		service_version,
		service_extra_info,
		service_method,
		service_conf,
		service_cpe
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

const insertPortsScriptsSQL = `
	INSERT INTO ports_scripts (
		ports_id,
		script_id,
		script_output
	) VALUES (?, ?, ?)`

func (p *PortRepository) insertRecords(host *Host) error {
	for _, port := range host.Port {
		id, err := p.insertPort(&port)
		if err != nil {
			return err
		}
		err = p.insertPortScripts(id, port.Script)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PortRepository) insertPort(port *Port) (int64, error) {
	return p.sqlite.insertReturnID(
		insertPortsSQL,
		p.hostID,
		port.PortID,
		port.State.State,
		port.State.Reason,
		port.State.ReasonTTL,
		port.Service.Name,
		port.Service.Product,
		port.Service.Version,
		port.Service.ExtraInfo,
		port.Service.Method,
		port.Service.Conf,
		strings.Join(port.Service.CPE, sqliteStringDelimiter),
	)
}

func (p *PortRepository) insertPortScripts(portID int64, scripts []Script) error {
	for _, script := range scripts {
		err := p.sqlite.insert(
			insertPortsScriptsSQL,
			portID,
			script.ID,
			script.Output,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
