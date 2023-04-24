package formatter

import (
	"database/sql"
	"strings"
	"time"
)

const sqliteStringDelimiter = ", "

const insertScanSQL = `
	INSERT INTO scans (
		nf_identifier, 
		scanner, 
		args,
		scan_info_type,
		scan_info_protocol,
		scan_info_num_services,
		scan_info_services,
		run_stats_finished_time,
		run_stats_finished_time_str,
		run_stats_finished_elapsed,
		run_stats_finished_summary,
		run_stats_finished_exit,
		run_stats_stat_hosts_up,
		run_stats_stat_hosts_down,
		run_stats_stat_hosts_total,
		verbose_level,
		debugging_level,
		start, 
		start_str, 
		nf_created
	) 
	VALUES 
	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

const insertHostsSQL = `
	INSERT INTO hosts (
		scan_id,
		nf_address_joined,
		nf_host_names_joined,
		start_time,
		end_time,
		status_state,
		status_reason,
		uptime_seconds,
		uptime_last_boot,
		distance_value,
		tcp_sequence_index,
		tcp_sequence_difficulty,
		tcp_sequence_values,
		ip_id_sequence_class,
		ip_id_sequence_values,
		tcp_ts_sequence_class,
		tcp_ts_sequence_values,
		traces_port,
		traces_protocol,
		status
	)
	VALUES
	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

const insertHostTracesHopsSQL = `
	INSERT INTO host_traces_hops (
		host_id,
		ttl,
		ip_address,
		rtt,
		host
	) VALUES (?, ?, ?, ?, ?)`

const insertHostAddressesSQL = `
	INSERT INTO host_addresses (
		host_id,
		address,
		address_type
	) VALUES (?, ?, ?)`

const insertHostNamesSQL = `
	INSERT INTO host_names (
		host_id,
		name,
		type
	) VALUES (?, ?, ?)`

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

// insertScan creates a new scan record in database and returns it's last inserted ID
func (f *SqliteFormatter) insertScan(db *sql.DB, n *NMAPRun) (int64, error) {
	id := int64(0)
	insert, err := db.Prepare(insertScanSQL)
	if err != nil {
		return id, err
	}
	defer insert.Close()

	now := time.Now()
	scanInsertResult, err := insert.Exec(
		f.getScanIdentifier(),
		n.Scanner,
		n.Args,
		n.ScanInfo.Type,
		n.ScanInfo.Protocol,
		n.ScanInfo.NumServices,
		n.ScanInfo.Services,
		n.RunStats.Finished.Time,
		n.RunStats.Finished.TimeStr,
		n.RunStats.Finished.Elapsed,
		n.RunStats.Finished.Summary,
		n.RunStats.Finished.Exit,
		n.RunStats.Hosts.Up,
		n.RunStats.Hosts.Down,
		n.RunStats.Hosts.Total,
		n.Verbose.Level,
		n.Debugging.Level,
		n.Start,
		n.StartStr,
		now.Unix(),
	)
	if err != nil {
		return id, err
	}
	return scanInsertResult.LastInsertId()
}

func (f *SqliteFormatter) insertReturnID(db *sql.DB, sql string, args ...any) (int64, error) {
	insert, err := db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer insert.Close()
	result, err := insert.Exec(args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (f *SqliteFormatter) insert(db *sql.DB, sql string, args ...any) error {
	insert, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	defer insert.Close()
	_, err = insert.Exec(args)
	return err
}

func (f *SqliteFormatter) insertHosts(db *sql.DB, scanID int64, hosts []Host) ([]int64, error) {
	var ids []int64
	insert, err := db.Prepare(insertHostsSQL)
	if err != nil {
		return ids, err
	}
	defer insert.Close()

	for _, host := range hosts {
		insertResult, err := insert.Exec(
			scanID,
			"TODO",
			"TODO",
			host.StartTime,
			host.EndTime,
			host.Status.State,
			host.Status.Reason,
			host.Uptime.Seconds,
			host.Uptime.LastBoot,
			host.Distance.Value,
			host.TCPSequence.Index,
			host.TCPSequence.Difficulty,
			host.TCPSequence.Values,
			host.IPIDSequence.Class,
			host.IPIDSequence.Values,
			host.TCPTSSequence.Class,
			host.TCPTSSequence.Values,
			host.Trace.Port,
			host.Trace.Protocol,
			host.Status,
		)
		if err != nil {
			return []int64{}, err
		}
		id, err := insertResult.LastInsertId()
		if err != nil {
			return []int64{}, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (f *SqliteFormatter) insertHostTracesHops(db *sql.DB, hostID int64, hops []Hop) error {
	var err error
	insert, err := db.Prepare(insertHostTracesHopsSQL)
	if err != nil {
		return err
	}
	defer insert.Close()
	for _, hop := range hops {
		_, err := insert.Exec(
			hostID,
			hop.TTL,
			hop.IPAddr,
			hop.RTT,
			hop.Host,
		)
		if err != nil {
			return err
		}
	}
	return err
}

func (f *SqliteFormatter) insertHostAddresses(db *sql.DB, hostID int64, addresses []HostAddress) error {
	var err error
	insert, err := db.Prepare(insertHostAddressesSQL)
	if err != nil {
		return err
	}
	defer insert.Close()
	for _, hostAddress := range addresses {
		_, err := insert.Exec(
			hostID,
			hostAddress.Address,
			hostAddress.AddressType,
		)
		if err != nil {
			return err
		}
	}
	return err
}

func (f *SqliteFormatter) insertHostNames(db *sql.DB, hostID int64, names *HostNames) error {
	var err error
	insert, err := db.Prepare(insertHostNamesSQL)
	if err != nil {
		return err
	}
	defer insert.Close()
	for _, hostName := range names.HostName {
		_, err := insert.Exec(
			hostID,
			hostName.Name,
			hostName.Type,
		)
		if err != nil {
			return err
		}
	}
	return err
}

func (f *SqliteFormatter) insertOSRecords(db *sql.DB, hostID int64, os *OS) (int64, error) {
	var err error
	insert, err := db.Prepare(insertHostOSSQL)
	if err != nil {
		return 0, err
	}
	defer insert.Close()
	result, err := insert.Exec(
		hostID,
		os.OSClass.Type,
		os.OSClass.Vendor,
		os.OSClass.OSFamily,
		os.OSClass.OSGen,
		os.OSClass.Accuracy,
		strings.Join(os.OSClass.CPE, sqliteStringDelimiter),
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (f *SqliteFormatter) insertOSPortUsed(db *sql.DB, osID int64, portUsed []OSPortUsed) error {
	insert, err := db.Prepare(insertHostOSPortUsedSQL)
	if err != nil {
		return err
	}
	defer insert.Close()
	for _, port := range portUsed {
		_, err := insert.Exec(
			osID,
			port.State,
			port.Protocol,
			port.PortID,
		)
		if err != nil {
			return err
		}
	}
	return err
}

func (f *SqliteFormatter) insertOSMatch(db *sql.DB, osID int64, match []OSMatch) error {
	insert, err := db.Prepare(insertHostOSMatchSQL)
	if err != nil {
		return err
	}
	defer insert.Close()
	for _, matchRecord := range match {
		_, err = insert.Exec(
			osID,
			matchRecord.Name,
			matchRecord.Accuracy,
			matchRecord.Line,
		)
		if err != nil {
			return err
		}
	}
	return err
}

func (f *SqliteFormatter) insertPorts(db *sql.DB, hostID int64, n []Port) ([]int64, error) {
	var portIDs []int64
	insert, err := db.Prepare(insertPortsSQL)
	if err != nil {
		return portIDs, err
	}
	defer insert.Close()
	for _, portRecord := range n {
		result, err := insert.Exec(
			hostID,
			portRecord.PortID,
			portRecord.State.State,
			portRecord.State.Reason,
			portRecord.State.ReasonTTL,
			portRecord.Service.Name,
			portRecord.Service.Product,
			portRecord.Service.Version,
			portRecord.Service.ExtraInfo,
			portRecord.Service.Method,
			portRecord.Service.Conf,
			strings.Join(portRecord.Service.CPE, sqliteStringDelimiter),
		)
		if err != nil {
			return portIDs, err
		}
		portID, err := result.LastInsertId()
		if err != nil {
			return portIDs, err
		}
		portIDs = append(portIDs, portID)
	}
	return portIDs, err
}

func (f *SqliteFormatter) insertPortScripts(db *sql.DB, portsID int64, n []Script) error {
	insert, err := db.Prepare(insertPortsScriptsSQL)
	if err != nil {
		return err
	}
	defer insert.Close()
	for _, script := range n {
		_, err := insert.Exec(
			portsID,
			script.ID,
			script.Output,
		)
		if err != nil {
			return err
		}
	}
	return err
}
