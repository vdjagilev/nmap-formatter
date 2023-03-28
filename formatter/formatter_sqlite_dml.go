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
		run_stats_stat_hosts_total,
		verbose_level,
		debugging_level,
		start, 
		start_str, 
		nf_created
	) 
	VALUES 
	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

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

func (f *SqliteFormatter) insertHosts(db *sql.DB, hosts []Host) error {
	return nil
}

func (f *SqliteFormatter) insertHostTraces(db *sql.DB, trace *Trace) error {
	return nil
}

func (f *SqliteFormatter) insertHostTracesHops(db *sql.DB, hops []Hop) error {
	return nil
}

func (f *SqliteFormatter) insertHostAddresses(db *sql.DB, addresses []HostAddress) error {
	return nil
}

func (f *SqliteFormatter) insertHostNames(db *sql.DB, names *HostNames) error {
	return nil
}

func (f *SqliteFormatter) insertOSRecords(db *sql.DB, os *OS) error {
	return nil
}

func (f *SqliteFormatter) insertPorts(db *sql.DB, n []Port) error {
	return nil
}

func (f *SqliteFormatter) insertPortScripts(db *sql.DB, n []Script) error {
	return nil
}
