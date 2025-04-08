package zabbix

// Tag represents a host tag in Zabbix.
type Tag struct {
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}
