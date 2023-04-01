package formatter

import (
	"database/sql"
	"fmt"
)

// SqliteDdl describes the whole database schema for sqlite
const SqliteDdl = `
CREATE TABLE IF NOT EXISTS scans (
	id integer not null primary key,
	nf_identifier text,
	scanner text,
	args text,
	scan_info_type text,
	scan_info_protocol text,
	scan_info_num_services integer,
	scan_info_services text,
	run_stats_finished_time integer,
	run_stats_finished_time_str text,
	run_stats_finished_elapsed real,
	run_stats_finished_summary text,
	run_stats_finished_exit text,
	run_stats_stat_hosts_up integer,
	run_stats_stat_hosts_down integer,
	run_stats_stat_hosts_total integer,
	verbose_level integer,
	debugging_level integer,
	start integer,
	start_str text,
	nf_created integer
);
CREATE TABLE IF NOT EXISTS hosts (
	id integer not null primary key,
	scan_id integer not null,
	nf_address_joined text,
	nf_host_names_joined text,
	nf_created integer,
	start_time integer,
	end_time integer,
	status_state text,
	status_reason text,
	uptime_seconds integer,
	uptime_last_boot string,
	distance_value integer,
	tcp_sequence_index text,
	tcp_sequence_difficulty text,
	tcp_sequence_values text,
	ip_id_sequence_class text,
	ip_id_sequence_values text,
	tcp_ts_sequence_class text,
	tcp_ts_sequence_values text,
	status text
);
CREATE TABLE IF NOT EXISTS host_traces (
	id integer not null primary key,
	host_id integer not null,
	port integer,
	protocol text
);
CREATE TABLE IF NOT EXISTS host_traces_hops (
	id integer not null primary key,
	host_traces_id integer not null,
	ttl integer,
	ip_address text,
	rtt real,
	host text
);
CREATE TABLE IF NOT EXISTS host_addresses (
	id integer not null primary key,
	host_id integer not null,
	address text,
	address_type text
);
CREATE TABLE IF NOT EXISTS host_names (
	id integer not null primary key,
	host_id integer not null,
	name text,
	type text
);
CREATE TABLE IF NOT EXISTS host_os (
	id integer not null primary key,
	host_id integer not null,
	class_type text,
	class_vendor text,
	class_osfamily text,
	class_osgen text,
	class_accuracy text,
	cpe text
);
CREATE TABLE IF NOT EXISTS host_os_port_used (
	id integer not null primary key,
	host_os_id integer not null,
	state text,
	protocol text,
	port_id integer
);
CREATE TABLE IF NOT EXISTS host_os_match (
	id integer not null primary key,
	host_os_id integer not null,
	name text,
	accuracy text,
	line text
);
CREATE TABLE IF NOT EXISTS ports (
	id integer not null primary key,
	host_id integer not null,
	port_id integer,
	state_state text,
	state_reason text,
	state_reason_ttl text,
	service_name text,
	service_product text,
	service_version text,
	service_extra_info text,
	service_method text,
	service_conf text,
	service_cpe text
);
CREATE TABLE IF NOT EXISTS ports_scripts (
	id integer not null primary key,
	ports_id integer not null,
	script_id text,
	script_output text
);
CREATE TABLE IF NOT EXISTS nf_schema (
	version text
);`

func (f *SqliteFormatter) generateSchema(db *sql.DB) error {
	// Create migrate schema
	_, err := db.Exec(SqliteDdl)
	if err != nil {
		return err
	}

	// Set schema version by truncating table and inserting new version
	_, err = db.Exec(`DELETE FROM nf_schema;`)
	if err != nil {
		return fmt.Errorf("could not clean nf_schema table: %v", err)
	}

	_, err = db.Exec(`INSERT INTO nf_schema VALUES (?);`, f.config.CurrentVersion)
	if err != nil {
		return fmt.Errorf("could not insert new nf_schema version: %v", err)
	}

	return nil
}

// schemaExists queries database and tries to select data from nf_schema
// database table, if that fails it's a clear indicator that schema
// does not exist in a database
func (f *SqliteFormatter) schemaExists(db *sql.DB) bool {
	// Try to get data from the database
	rows, err := db.Query(`SELECT version FROM nf_schema LIMIT 1`)
	if err != nil {
		return false
	}
	defer rows.Close()
	return true
}
