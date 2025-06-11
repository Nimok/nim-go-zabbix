package zabbix

import "context"

type Proxy struct {
	ProxyID              string `json:"proxyid,omitempty"`                // ID of the proxy; read-only, required for update operations
	Name                 string `json:"name,omitempty"`                   // Name of the proxy; required for create operations
	ProxyGroupID         string `json:"proxy_groupid,omitempty"`          // ID of the proxy group; 0 if not assigned to any group
	LocalAddress         string `json:"local_address,omitempty"`          // Address for active agents; required if proxy_groupid is not 0
	LocalPort            string `json:"local_port,omitempty"`             // Local proxy port number; default is 10051
	OperatingMode        int    `json:"operating_mode"`                   // Type of proxy; 0 for active, 1 for passive; required for create operations
	Description          string `json:"description,omitempty"`            // Description of the proxy
	LastAccess           int64  `json:"lastaccess,omitempty"`             // Time when the proxy last connected to the server; read-only
	Address              string `json:"address,omitempty"`                // IP address or DNS name to connect to; required if operating_mode is passive
	Port                 string `json:"port,omitempty"`                   // Port number to connect to; default is 10051
	AllowedAddresses     string `json:"allowed_addresses,omitempty"`      // Comma-delimited IP addresses or DNS names of active Zabbix proxy
	TLSConnect           int    `json:"tls_connect,omitempty"`            // Connections to host; 1 (default) No encryption, 2 PSK, 4 certificate
	TLSAccept            int    `json:"tls_accept,omitempty"`             // Connections from host; bitmask: 1 (default) No encryption, 2 PSK, 4 certificate
	TLSIssuer            string `json:"tls_issuer,omitempty"`             // Certificate issuer
	TLSSubject           string `json:"tls_subject,omitempty"`            // Certificate subject
	TLSPskIdentity       string `json:"tls_psk_identity,omitempty"`       // PSK identity; write-only, required if TLSConnect or TLSAccept includes PSK
	TLSPsk               string `json:"tls_psk,omitempty"`                // Pre-shared key (PSK); write-only, required if TLSConnect or TLSAccept includes PSK
	CustomTimeouts       int    `json:"custom_timeouts,omitempty"`        // Whether to override global item timeouts; 0 (default) use global settings, 1 override timeouts
	TimeoutZabbixAgent   string `json:"timeout_zabbix_agent,omitempty"`   // Timeout for Zabbix agent checks; required if CustomTimeouts is 1
	TimeoutSimpleCheck   string `json:"timeout_simple_check,omitempty"`   // Timeout for simple checks; required if CustomTimeouts is 1
	TimeoutSnmpAgent     string `json:"timeout_snmp_agent,omitempty"`     // Timeout for SNMP agent checks; required if CustomTimeouts is 1
	TimeoutExternalCheck string `json:"timeout_external_check,omitempty"` // Timeout for external checks; required if CustomTimeouts is 1
	TimeoutDbMonitor     string `json:"timeout_db_monitor,omitempty"`     // Timeout for database monitoring; required if CustomTimeouts is 1
	TimeoutHttpAgent     string `json:"timeout_http_agent,omitempty"`     // Timeout for HTTP agent checks; required if CustomTimeouts is 1
	TimeoutSshAgent      string `json:"timeout_ssh_agent,omitempty"`      // Timeout for SSH agent checks; required if CustomTimeouts is 1
	TimeoutTelnetAgent   string `json:"timeout_telnet_agent,omitempty"`   // Timeout for Telnet agent checks; required if CustomTimeouts is 1
	TimeoutScript        string `json:"timeout_script,omitempty"`         // Timeout for script checks; required if CustomTimeouts is 1
	TimeoutBrowser       string `json:"timeout_browser,omitempty"`        // Timeout for browser checks; required if CustomTimeouts is 1
	Version              int    `json:"version,omitempty"`                // Version of proxy; read-only
	Compatibility        int    `json:"compatibility,omitempty"`          // Version compatibility with Zabbix server; read-only
	State                int    `json:"state,omitempty"`                  // State of the proxy; read-only
}

type ProxyGetParameters struct {
	GetParameters

	ProxyIDs            []string `json:"proxyids,omitempty"`
	ProxyGroupIDs       []string `json:"proxy_groupids,omitempty"`
	SelectAssignedHosts any      `json:"selectAssignedHosts,omitempty"`
	SelectHosts         any      `json:"selectHosts,omitempty"`
	SelectProxyGroup    any      `json:"selectProxyGroup,omitempty"`
	SortField           any      `json:"sortfield,omitempty"`
}

type ProxyCreateParameters struct {
	Proxy

	Hosts []Host `json:"hosts,omitempty"`
}

type ProxyCreateResponse struct {
	ProxyIDs []string `json:"proxyids"` // IDs of the created proxies
}

type ProxyDeleteResponse struct {
	ProxyIDs []string `json:"proxyids"` // IDs of the deleted proxies
}

func (z *zabbixClient) ProxyGet(ctx context.Context, params ProxyGetParameters) ([]Proxy, error) {

	var result []Proxy

	err := z.makeRequest(ctx, "proxy.get", params, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (z *zabbixClient) ProxyCreate(ctx context.Context, params ProxyCreateParameters) (*ProxyCreateResponse, error) {

	var result ProxyCreateResponse

	err := z.makeRequest(ctx, "proxy.create", params, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (z *zabbixClient) ProxyDelete(ctx context.Context, params []string) (*ProxyDeleteResponse, error) {

	var result ProxyDeleteResponse

	err := z.makeRequest(ctx, "proxy.delete", params, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
