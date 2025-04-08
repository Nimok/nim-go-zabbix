package zabbix

type GetParameters struct {
	Output                 any               `json:"output,omitempty"`
	CountOutput            bool              `json:"countOutput,omitempty"`
	Editable               bool              `json:"editable,omitempty"`
	ExcludeSearch          bool              `json:"excludeSearch,omitempty"`
	Filter                 map[string]any    `json:"filter,omitempty"`
	Limit                  int               `json:"limit,omitempty"`
	Search                 map[string]string `json:"search,omitempty"`
	SearchByAny            bool              `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool              `json:"searchWildcardsEnabled,omitempty"`
	Sortfield              []string          `json:"sortfield,omitempty"`
	Sortorder              any               `json:"sortorder,omitempty"`
	StartSearch            bool              `json:"startSearch,omitempty"`
	PreserveKeys           bool              `json:"preservekeys,omitempty"`
	SelectHosts            any               `json:"selectHosts,omitempty"`
	SelectItems            any               `json:"selectItems,omitempty"`
	SelectTriggers         any               `json:"selectTriggers,omitempty"`
	// Add other select parameters as needed
}
