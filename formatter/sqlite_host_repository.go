package formatter

import "database/sql"

// HostRepository main responsibility is to populated database
// with host related data
type HostRepository struct {
	sqlite *SqliteDB
	conn   *sql.DB
	scanID int64
}

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

const insertHostTracesHopsSQL = `
	INSERT INTO host_traces_hops (
		host_id,
		ttl,
		ip_address,
		rtt,
		host
	) VALUES (?, ?, ?, ?, ?)`

func (h *HostRepository) insertHosts(hosts []Host) error {
	for _, host := range hosts {
		id, err := h.insertHost(host)
		if err != nil {
			return err
		}
		err = h.insertHostNetInfo(id, &host)
		if err != nil {
			return err
		}
		err = h.insertHostInfo(id, &host)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *HostRepository) insertHost(host Host) (int64, error) {
	return h.sqlite.insertReturnID(
		insertHostsSQL,
		h.scanID,
		host.JoinedAddresses(sqliteStringDelimiter),
		host.JoinedHostNames(sqliteStringDelimiter),
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
		host.Status.State,
	)
}

func (h *HostRepository) insertHostNetInfo(hostID int64, host *Host) error {
	err := h.insertTracesHops(hostID, host.Trace.Hops)
	if err != nil {
		return err
	}
	err = h.insertHostAddresses(hostID, host.HostAddress)
	if err != nil {
		return err
	}
	return h.insertHostNames(hostID, host.HostNames)
}

func (h *HostRepository) insertTracesHops(hostID int64, hops []Hop) error {
	for _, hop := range hops {
		err := h.sqlite.insert(
			insertHostTracesHopsSQL,
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
	return nil
}

func (h *HostRepository) insertHostAddresses(hostID int64, addresses []HostAddress) error {
	for _, addr := range addresses {
		err := h.sqlite.insert(
			insertHostAddressesSQL,
			hostID,
			addr.Address,
			addr.AddressType,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *HostRepository) insertHostNames(hostID int64, names HostNames) error {
	for _, name := range names.HostName {
		err := h.sqlite.insert(
			insertHostNamesSQL,
			hostID,
			name.Name,
			name.Type,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *HostRepository) insertHostInfo(hostID int64, host *Host) error {
	os := &OSRepository{
		sqlite: h.sqlite,
		conn:   h.conn,
		hostID: hostID,
	}
	err := os.insertOSRecords(&host.OS)
	if err != nil {
		return err
	}
	ports := &PortRepository{
		sqlite: h.sqlite,
		conn:   h.conn,
		hostID: hostID,
	}
	return ports.insertRecords(host)
}
