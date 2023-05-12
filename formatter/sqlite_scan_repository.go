package formatter

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ScanRepository struct {
	sqlite *SqliteDB
	conn   *sql.DB
	config *Config
}

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

func (s *ScanRepository) insertScan(n *NMAPRun) error {
	// The nf_ prefix in tables are related to nmap-formatter
	// either the creation date or passed options (identifier)
	// Identifiers are needed to help users to differentiate between scans
	now := time.Now()
	id, err := s.sqlite.insertReturnID(
		insertScanSQL,
		s.getScanIdentifier(),
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
		return err
	}

	hostRepository := &HostRepository{
		conn:   s.conn,
		scanID: id,
		sqlite: s.sqlite,
	}

	return hostRepository.insertHosts(n.Host)
}

// getScanIdentifier returns a unique string provided either by a user or generates new random uuid
func (s *ScanRepository) getScanIdentifier() string {
	if s.config.OutputOptions.SqliteOutputOptions.ScanIdentifier == "" {
		return generateUUID()
	}
	return s.config.OutputOptions.SqliteOutputOptions.ScanIdentifier
}

// generateUUID generates new random uuid string
func generateUUID() string {
	return uuid.New().String()
}
