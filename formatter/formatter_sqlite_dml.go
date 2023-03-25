package formatter

import (
	"database/sql"
	"time"
)

const INSERT_SCAN_SQL = `
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
	insert, err := db.Prepare(INSERT_SCAN_SQL)
	if err != nil {
		return id, err
	}
	defer insert.Close()

	now := time.Now()
	scan_insert_result, err := insert.Exec(
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
	return scan_insert_result.LastInsertId()
}
