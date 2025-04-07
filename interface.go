package zabbix

// ---------------------------
// Host interface types
// ---------------------------

const (
	InterfaceAgent = 1
	InterfaceSNMP  = 2
	InterfaceIPMI  = 3
	InterfaceJMX   = 4
)

// HostInterface represents a host interface in Zabbix.
type HostInterface struct {
	Type    int              `json:"type"`
	Main    int              `json:"main"`
	UseIP   int              `json:"useip"`
	IP      string           `json:"ip"`
	DNS     string           `json:"dns"`
	Port    string           `json:"port"`
	Details InterfaceDetails `json:"details,omitempty"`
}

// InterfaceDetails represents the details of a host interface, particularly for SNMP.
type InterfaceDetails struct {
	Version       int    `json:"version,omitempty"`
	Bulk          int    `json:"bulk,omitempty"`
	SecurityName  string `json:"securityname,omitempty"`
	ContextName   string `json:"contextname,omitempty"`
	SecurityLevel int    `json:"securitylevel,omitempty"`
}
