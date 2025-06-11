package zabbix

import "context"

// Host represents a Zabbix host object.
type Host struct {
	HostID            string          `json:"hostid,omitempty"`             // ID of the host (read-only; required for update operations)
	Host              string          `json:"host,omitempty"`               // Technical name of the host (required for create operations)
	Description       string          `json:"description,omitempty"`        // Description of the host
	Flags             int             `json:"flags,omitempty"`              // Origin of the host (0 - plain host; 4 - discovered host) (read-only)
	InventoryMode     int             `json:"inventory_mode,omitempty"`     // Host inventory population mode (-1 - disabled; 0 - manual; 1 - automatic)
	IPMIAuthtype      int             `json:"ipmi_authtype,omitempty"`      // IPMI authentication algorithm (-1 - default; 0 - none; 1 - MD2; 2 - MD5; 4 - straight; 5 - OEM; 6 - RMCP+)
	IPMIPassword      string          `json:"ipmi_password,omitempty"`      // IPMI password
	IPMIPrivilege     int             `json:"ipmi_privilege,omitempty"`     // IPMI privilege level (1 - callback; 2 - user; 3 - operator; 4 - admin; 5 - OEM)
	IPMIUsername      string          `json:"ipmi_username,omitempty"`      // IPMI username
	MaintenanceFrom   int64           `json:"maintenance_from,omitempty"`   // Starting time of the effective maintenance (read-only)
	MaintenanceStatus int             `json:"maintenance_status,omitempty"` // Effective maintenance status (0 - no maintenance; 1 - maintenance in effect) (read-only)
	MaintenanceType   int             `json:"maintenance_type,omitempty"`   // Effective maintenance type (0 - with data collection; 1 - without data collection) (read-only)
	MaintenanceID     string          `json:"maintenanceid,omitempty"`      // ID of the maintenance currently in effect on the host (read-only)
	Name              string          `json:"name,omitempty"`               // Visible name of the host (defaults to 'host' property value)
	MonitoredBy       int             `json:"monitored_by,omitempty"`       // Source used to monitor the host (0 - Zabbix server; 1 - Proxy; 2 - Proxy group)
	ProxyID           string          `json:"proxyid,omitempty"`            // ID of the proxy monitoring the host (required if 'monitored_by' is set to Proxy)
	ProxyGroupID      string          `json:"proxy_groupid,omitempty"`      // ID of the proxy group monitoring the host (required if 'monitored_by' is set to Proxy group)
	Status            int             `json:"status,omitempty"`             // Status and function of the host (0 - monitored; 1 - unmonitored)
	TlsConnect        int             `json:"tls_connect,omitempty"`        // Connections to host (1 - No encryption; 2 - PSK; 4 - certificate)
	TlsAccept         int             `json:"tls_accept,omitempty"`         // Connections from host (bitmask: 1 - No encryption; 2 - PSK; 4 - certificate)
	TlsIssuer         string          `json:"tls_issuer,omitempty"`         // Certificate issuer
	TlsSubject        string          `json:"tls_subject,omitempty"`        // Certificate subject
	TlsPSKIdentity    string          `json:"tls_psk_identity,omitempty"`   // PSK identity (write-only; required if 'tls_connect' is PSK or 'tls_accept' includes PSK)
	TlsPSK            string          `json:"tls_psk,omitempty"`            // Pre-shared key (PSK) (write-only; required if 'tls_connect' is PSK or 'tls_accept' includes PSK)
	ActiveAvailable   int             `json:"active_available,omitempty"`   // Host active interface availability status (0 - unknown; 1 - available; 2 - not available) (read-only)
	AssignedProxyID   string          `json:"assigned_proxyid,omitempty"`   // ID of the proxy assigned by Zabbix server if monitored by a proxy group (read-only)
	Interfaces        []HostInterface `json:"interfaces,omitempty"`         // Interfaces associated with the host
	Groups            []HostGroup     `json:"groups,omitempty"`             // Host groups to which the host belongs
	Tags              []Tag           `json:"tags,omitempty"`               // Tags associated with the host
	Templates         []Template      `json:"templates,omitempty"`          // Templates linked to the host
	Macros            []Macro         `json:"macros,omitempty"`             // User macros created for the host
	Inventory         Inventory       `json:"inventory,omitempty"`          // Inventory properties of the host
}

type HostGetParameters struct {
	GetParameters

	HostIDs                []string          `json:"hostids,omitempty"`
	GroupIDs               []string          `json:"groupids,omitempty"`
	ApplicationIDs         []string          `json:"applicationids,omitempty"`
	DServiceIDs            []string          `json:"dserviceids,omitempty"`
	GraphIDs               []string          `json:"graphids,omitempty"`
	HttpTestIDs            []string          `json:"httptestids,omitempty"`
	InterfaceIDs           []string          `json:"interfaceids,omitempty"`
	ItemIDs                []string          `json:"itemids,omitempty"`
	MaintenanceIDs         []string          `json:"maintenanceids,omitempty"`
	MonitoredHosts         bool              `json:"monitored_hosts,omitempty"`
	ProxyHosts             bool              `json:"proxy_hosts,omitempty"`
	ProxyIDs               []string          `json:"proxyids,omitempty"`
	TemplatedHosts         bool              `json:"templated_hosts,omitempty"`
	TemplateIDs            []string          `json:"templateids,omitempty"`
	TriggerIDs             []string          `json:"triggerids,omitempty"`
	WithItems              bool              `json:"with_items,omitempty"`
	WithApplications       bool              `json:"with_applications,omitempty"`
	WithGraphs             bool              `json:"with_graphs,omitempty"`
	WithHttpTests          bool              `json:"with_httptests,omitempty"`
	WithMonitoredHttpTests bool              `json:"with_monitored_httptests,omitempty"`
	WithMonitoredItems     bool              `json:"with_monitored_items,omitempty"`
	WithMonitoredTriggers  bool              `json:"with_monitored_triggers,omitempty"`
	WithSimpleGraphItems   bool              `json:"with_simple_graph_items,omitempty"`
	WithTriggers           bool              `json:"with_triggers,omitempty"`
	SelectGroups           any               `json:"selectGroups,omitempty"`
	SelectApplications     any               `json:"selectApplications,omitempty"`
	SelectDiscoveries      any               `json:"selectDiscoveries,omitempty"`
	SelectDiscoveryRule    any               `json:"selectDiscoveryRule,omitempty"`
	SelectGraphs           any               `json:"selectGraphs,omitempty"`
	SelectHostDiscovery    any               `json:"selectHostDiscovery,omitempty"`
	SelectHttpTests        any               `json:"selectHttpTests,omitempty"`
	SelectInterfaces       any               `json:"selectInterfaces,omitempty"`
	SelectInventory        any               `json:"selectInventory,omitempty"`
	SelectMacros           any               `json:"selectMacros,omitempty"`
	SelectParentTemplates  any               `json:"selectParentTemplates,omitempty"`
	SelectScreens          any               `json:"selectScreens,omitempty"`
	LimitSelects           int               `json:"limitSelects,omitempty"`
	SearchInventory        map[string]string `json:"searchInventory,omitempty"`
}

type HostCreateResponse struct {
	HostIDs []string `json:"hostids"` // IDs of the created hosts
}

type HostDeleteResponse struct {
	HostIDs []string `json:"hostids"` // IDs of the deleted hosts
}

type HostUpdateResponse struct {
	HostIDs []string `json:"hostids"` // IDs of the updated host
}

func (z *zabbixClient) HostCreate(ctx context.Context, params Host) (*HostCreateResponse, error) {

	var result HostCreateResponse

	err := z.makeRequest(ctx, "host.create", params, &result)
	if err != nil {
		return nil, err
	}

	return &result, err
}

func (z *zabbixClient) HostDelete(ctx context.Context, params []string) (*HostDeleteResponse, error) {

	var result HostDeleteResponse

	err := z.makeRequest(ctx, "host.delete", params, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (z *zabbixClient) HostUpdate(ctx context.Context, params Host) (*HostUpdateResponse, error) {

	var result HostUpdateResponse

	err := z.makeRequest(ctx, "host.update", params, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (z *zabbixClient) HostGet(ctx context.Context, params HostGetParameters) ([]Host, error) {

	var result []Host

	err := z.makeRequest(ctx, "host.get", params, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
