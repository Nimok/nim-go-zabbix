package zabbix

// Macro represents a user macro in Zabbix.
type Macro struct {
	Macro       string `json:"macro"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
}
