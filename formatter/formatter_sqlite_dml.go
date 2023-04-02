package formatter

import (
	"database/sql"
	"time"
)

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
		nf_created,
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
	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

const insertHostTracesHopsSQL = `
	INSERT INTO host_traces_hops (
		host_id,
		ttl,
		ip_address
		rtt,
		host
	) VALUES (?, ?, ?, ?, ?)`

const insertHostAddressesSQL = `
	INSERT INTO host_addresses (
		host_id,
		address,
		address_type
	)`

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
		state_state text,
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
	INSERT INTOP ports_scripts (
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

func (f *SqliteFormatter) insertHostAddresses(db *sql.DB, addresses []HostAddress) error {
	return nil
}

func (f *SqliteFormatter) insertHostNames(db *sql.DB, names *HostNames) error {
	return nil
}

func (f *SqliteFormatter) insertOSRecords(db *sql.DB, os *OS) ([]int64, error) {
	return []int64{}, nil
}

func (f *SqliteFormatter) insertPorts(db *sql.DB, n []Port) ([]int64, error) {
	return []int64{}, nil
}

func (f *SqliteFormatter) insertPortScripts(db *sql.DB, n []Script) error {
	return nil
}
