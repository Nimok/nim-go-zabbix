package zabbix

// Inventory represents the host inventory properties in Zabbix.
type Inventory struct {
	MacAddressA string `json:"macaddress_a,omitempty"`
	MacAddressB string `json:"macaddress_b,omitempty"`
	// Add other inventory fields as needed
}
