package zabbix

// ---------------------------
// Host status and inventory
// ---------------------------

const (
	HostStatusMonitored   = 0
	HostStatusUnmonitored = 1

	InventoryDisabled = -1
	InventoryManual   = 0
	InventoryAuto     = 1
)

// ---------------------------
// TLS connect and accept
// ---------------------------

const (
	TLSNoEncryption = 1
	TLSPSK          = 2
	TLSCert         = 4
)

// Bitmask for TlsAccept:
// TLSNoEncryption | TLSPSK | TLSCert

// ---------------------------
// IPMI settings
// ---------------------------

const (
	IPMIAuthDefault = -1
	IPMINoAuth      = 0
	IPMIMD2         = 1
	IPMIMD5         = 2
	IPMIStraight    = 4
	IPMIOEM         = 5
	IPMIRMCPPlus    = 6

	IPMIPrivilegeCallback = 1
	IPMIPrivilegeUser     = 2
	IPMIPrivilegeOperator = 3
	IPMIPrivilegeAdmin    = 4
	IPMIPrivilegeOEM      = 5
)

// ---------------------------
// SNMP security level
// ---------------------------

const (
	SNMPSecurityNoAuthNoPriv = 0
	SNMPSecurityAuthNoPriv   = 1
	SNMPSecurityAuthPriv     = 2
)

// ---------------------------
// Host maintenance status
// ---------------------------

const (
	MaintenanceNone   = 0
	MaintenanceActive = 1
)

const (
	MaintenanceWithData = 0
	MaintenanceNoData   = 1
)

// ---------------------------
// Monitoring method
// ---------------------------

const (
	MonitoredByServer     = 0
	MonitoredByProxy      = 1
	MonitoredByProxyGroup = 2
)

// ---------------------------
// Availability (read-only)
// ---------------------------

const (
	AvailabilityUnknown     = 0
	AvailabilityAvailable   = 1
	AvailabilityUnavailable = 2
)
