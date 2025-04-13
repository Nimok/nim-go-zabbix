package zabbix

import "context"

// Hostgroup represents a host group in Zabbix.
type HostGroup struct {
	GroupID string `json:"groupid,omitempty"` // ID of the host group; read-only, required for update operations
	Name    string `json:"name,omitempty"`    // Name of the host group; required for create operations
	Flags   int    `json:"flags,omitempty"`   // Origin of the host group; read-only
	UUID    string `json:"uuid,omitempty"`    // Universal unique identifier; auto-generated if not provided
}

type HostGroupGetParameters struct {
	GetParameters

	GraphIDs                      []string `json:"graphids,omitempty"`
	GroupIDs                      []string `json:"groupids,omitempty"`
	HostIDs                       []string `json:"hostids,omitempty"`
	MaintenanceIDs                []string `json:"maintenanceids,omitempty"`
	TriggerIDs                    []string `json:"triggerids,omitempty"`
	WithGraphs                    bool     `json:"with_graphs,omitempty"`
	WithGraphPrototypes           bool     `json:"with_graph_prototypes,omitempty"`
	WithHosts                     bool     `json:"with_hosts,omitempty"`
	WithHTTPTests                 bool     `json:"with_httptests,omitempty"`
	WithItems                     bool     `json:"with_items,omitempty"`
	WithItemPrototypes            bool     `json:"with_item_prototypes,omitempty"`
	WithSimpleGraphItemPrototypes bool     `json:"with_simple_graph_item_prototypes,omitempty"`
	WithMonitoredHTTPTests        bool     `json:"with_monitored_httptests,omitempty"`
	WithMonitoredHosts            bool     `json:"with_monitored_hosts,omitempty"`
	WithMonitoredItems            bool     `json:"with_monitored_items,omitempty"`
	WithMonitoredTriggers         bool     `json:"with_monitored_triggers,omitempty"`
	WithSimpleGraphItems          bool     `json:"with_simple_graph_items,omitempty"`
	WithTriggers                  bool     `json:"with_triggers,omitempty"`
	SelectDiscoveryRules          any      `json:"selectDiscoveryRules,omitempty"`
	SelectGroupDiscoveries        any      `json:"selectGroupDiscoveries,omitempty"`
	SelectHostPrototypes          any      `json:"selectHostPrototypes,omitempty"`
	SelectHosts                   any      `json:"selectHosts,omitempty"`
	LimitSelects                  int      `json:"limitSelects,omitempty"`
	SortField                     any      `json:"sortfield,omitempty"`
}

func (z *zabbixClient) HostgroupGet(ctx context.Context, params HostGroupGetParameters) ([]HostGroup, error) {

	var result []HostGroup

	err := z.makeRequest(ctx, "hostgroup.get", params, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
