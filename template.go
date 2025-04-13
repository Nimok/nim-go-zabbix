package zabbix

import "context"

// Template represents a template to be linked to a host in Zabbix.
type Template struct {
	TemplateID    string `json:"templateid,omitempty"`     // ID of the template; read-only, required for update operations
	Host          string `json:"host,omitempty"`           // Technical name of the template; required for create operations
	Name          string `json:"name,omitempty"`           // Visible name of the template; defaults to 'host' if not set
	Description   string `json:"description,omitempty"`    // Description of the template
	UUID          string `json:"uuid,omitempty"`           // Universal unique identifier; auto-generated if not provided
	VendorName    string `json:"vendor_name,omitempty"`    // Template vendor name; both vendor_name and vendor_version should be set or left empty for create operations
	VendorVersion string `json:"vendor_version,omitempty"` // Template vendor version; both vendor_name and vendor_version should be set or left empty for create operations
}

type TemplateGetParameters struct {
	GetParameters

	TemplateIDs           []string            `json:"templateids,omitempty"`
	GroupIDs              []string            `json:"groupids,omitempty"`
	ParentTemplateIDs     []string            `json:"parentTemplateids,omitempty"`
	HostIDs               []string            `json:"hostids,omitempty"`
	GraphIDs              []string            `json:"graphids,omitempty"`
	ItemIDs               []string            `json:"itemids,omitempty"`
	TriggerIDs            []string            `json:"triggerids,omitempty"`
	WithItems             bool                `json:"with_items,omitempty"`
	WithTriggers          bool                `json:"with_triggers,omitempty"`
	WithGraphs            bool                `json:"with_graphs,omitempty"`
	WithHTTPTests         bool                `json:"with_httptests,omitempty"`
	EvalType              int                 `json:"evaltype,omitempty"`
	Tags                  []map[string]string `json:"tags,omitempty"`
	SelectTags            any                 `json:"selectTags,omitempty"`
	SelectHosts           any                 `json:"selectHosts,omitempty"`
	SelectTemplateGroups  any                 `json:"selectTemplateGroups,omitempty"`
	SelectTemplates       any                 `json:"selectTemplates,omitempty"`
	SelectParentTemplates any                 `json:"selectParentTemplates,omitempty"`
	SelectHttpTests       any                 `json:"selectHttpTests,omitempty"`
	SelectItems           any                 `json:"selectItems,omitempty"`
	SelectDiscoveries     any                 `json:"selectDiscoveries,omitempty"`
	SelectTriggers        any                 `json:"selectTriggers,omitempty"`
	SelectGraphs          any                 `json:"selectGraphs,omitempty"`
	SelectMacros          any                 `json:"selectMacros,omitempty"`
	SelectDashboards      any                 `json:"selectDashboards,omitempty"`
	SelectValueMaps       any                 `json:"selectValueMaps,omitempty"`
	LimitSelects          int                 `json:"limitSelects,omitempty"`
	SortField             any                 `json:"sortfield,omitempty"`
}

func (z *zabbixClient) TemplateGet(ctx context.Context, params TemplateGetParameters) ([]Template, error) {

	var result []Template

	err := z.makeRequest(ctx, "template.get", params, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
