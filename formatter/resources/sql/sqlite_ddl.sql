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
	traces_port integer,
	traces_protocol text,
	status text
);
CREATE TABLE IF NOT EXISTS host_traces_hops (
	id integer not null primary key,
	host_id integer not null,
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
CREATE TABLE IF NOT EXISTS host_os_class (
	id integer not null primary key,
	host_id integer not null,
	type text,
	vendor text,
	osfamily text,
	osgen text,
	accuracy text,
	cpe text
);
CREATE TABLE IF NOT EXISTS host_os_port_used (
	id integer not null primary key,
	host_id integer not null,
	state text,
	protocol text,
	port_id integer
);
CREATE TABLE IF NOT EXISTS host_os_match (
	id integer not null primary key,
	host_id integer not null,
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
);

-- Performance indexes for foreign key columns
-- These significantly improve query performance and JOIN operations
CREATE INDEX IF NOT EXISTS idx_hosts_scan_id ON hosts(scan_id);
CREATE INDEX IF NOT EXISTS idx_ports_host_id ON ports(host_id);
CREATE INDEX IF NOT EXISTS idx_host_addresses_host_id ON host_addresses(host_id);
CREATE INDEX IF NOT EXISTS idx_host_names_host_id ON host_names(host_id);
CREATE INDEX IF NOT EXISTS idx_host_traces_hops_host_id ON host_traces_hops(host_id);
CREATE INDEX IF NOT EXISTS idx_host_os_class_host_id ON host_os_class(host_id);
CREATE INDEX IF NOT EXISTS idx_host_os_port_used_host_id ON host_os_port_used(host_id);
CREATE INDEX IF NOT EXISTS idx_host_os_match_host_id ON host_os_match(host_id);
CREATE INDEX IF NOT EXISTS idx_ports_scripts_ports_id ON ports_scripts(ports_id);

-- Additional useful indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_ports_state ON ports(state_state);
CREATE INDEX IF NOT EXISTS idx_ports_service ON ports(service_name);
CREATE INDEX IF NOT EXISTS idx_hosts_status ON hosts(status_state);