package zabbix

// InterfaceType represents the type of interface.
type InterfaceType int

const (
	InterfaceTypeAgent InterfaceType = 1
	InterfaceTypeSNMP  InterfaceType = 2
	InterfaceTypeIPMI  InterfaceType = 3
	InterfaceTypeJMX   InterfaceType = 4
)

// UseIPOption indicates whether to use IP or DNS for the connection.
type UseIPOption int

const (
	UseIPOptionDNS UseIPOption = 0 // Connect using host DNS name
	UseIPOptionIP  UseIPOption = 1 // Connect using host IP address
)

// MainInterface indicates whether the interface is the default for its type.
type MainInterface int

const (
	MainInterfaceNo  MainInterface = 0 // Not default
	MainInterfaceYes MainInterface = 1 // Default
)

// InterfaceState represents the availability state of the interface.
type InterfaceState int

const (
	InterfaceStateUnknown     InterfaceState = 0 // Default
	InterfaceStateAvailable   InterfaceState = 1
	InterfaceStateUnavailable InterfaceState = 2
)

// SNMP Version Type
type SNMPVersion int

const (
	SNMPv1  SNMPVersion = 1
	SNMPv2c SNMPVersion = 2
	SNMPv3  SNMPVersion = 3
)

// Bulk Request Type
type BulkSetting int

const (
	BulkDisabled BulkSetting = 0
	BulkEnabled  BulkSetting = 1
)

// SNMPv3 Security Level Type
type SecurityLevel int

const (
	SecurityLevelNoAuthNoPriv SecurityLevel = 0
	SecurityLevelAuthNoPriv   SecurityLevel = 1
	SecurityLevelAuthPriv     SecurityLevel = 2
)

// SNMPv3 Auth Protocol Type
type AuthProtocol int

const (
	AuthProtocolMD5    AuthProtocol = 0
	AuthProtocolSHA1   AuthProtocol = 1
	AuthProtocolSHA224 AuthProtocol = 2
	AuthProtocolSHA256 AuthProtocol = 3
	AuthProtocolSHA384 AuthProtocol = 4
	AuthProtocolSHA512 AuthProtocol = 5
)

// SNMPv3 Privacy Protocol Type
type PrivProtocol int

const (
	PrivProtocolDES     PrivProtocol = 0
	PrivProtocolAES128  PrivProtocol = 1
	PrivProtocolAES192  PrivProtocol = 2
	PrivProtocolAES256  PrivProtocol = 3
	PrivProtocolAES192C PrivProtocol = 4
	PrivProtocolAES256C PrivProtocol = 5
)

// HostInterface represents a host interface in Zabbix.
type HostInterface struct {
	InterfaceID  string           `json:"interfaceid,omitempty"`   // Read-only; required for update operations
	HostID       string           `json:"hostid"`                  // Required for create operations
	Type         InterfaceType    `json:"type"`                    // Required for create operations
	IP           string           `json:"ip"`                      // Required for create operations
	DNS          string           `json:"dns"`                     // Required for create operations
	Port         string           `json:"port"`                    // Required for create operations
	UseIP        UseIPOption      `json:"useip"`                   // Required for create operations
	Main         MainInterface    `json:"main"`                    // Required for create operations
	Available    InterfaceState   `json:"available,omitempty"`     // Read-only
	DisableUntil int64            `json:"disable_until,omitempty"` // Read-only; timestamp
	Error        string           `json:"error,omitempty"`         // Read-only
	ErrorsFrom   int64            `json:"errors_from,omitempty"`   // Read-only; timestamp
	Details      InterfaceDetails `json:"details,omitempty"`       // Required if Type is SNMP
}

// InterfaceDetails represents the details of a host interface, particularly for SNMP.
type InterfaceDetails struct {
	Version        SNMPVersion   `json:"version"`                   // Required
	Bulk           BulkSetting   `json:"bulk,omitempty"`            // Optional
	Community      string        `json:"community,omitempty"`       // Required if Version is SNMPv1 or SNMPv2c
	MaxRepetitions int           `json:"max_repetitions,omitempty"` // Default: 10
	SecurityName   string        `json:"securityname,omitempty"`    // SNMPv3 only
	SecurityLevel  SecurityLevel `json:"securitylevel,omitempty"`   // SNMPv3 only
	AuthPassphrase string        `json:"authpassphrase,omitempty"`  // SNMPv3 only
	PrivPassphrase string        `json:"privpassphrase,omitempty"`  // SNMPv3 only
	AuthProtocol   AuthProtocol  `json:"authprotocol,omitempty"`    // SNMPv3 only
	PrivProtocol   PrivProtocol  `json:"privprotocol,omitempty"`    // SNMPv3 only
	ContextName    string        `json:"contextname,omitempty"`     // SNMPv3 only
}
