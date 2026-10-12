package zabbix

import (
	"context"
)

// Item represents a Zabbix item.
//
// Fields where 0 is a meaningful value (e.g. Type ItemTypeZabbixAgent, ValueType ItemValueTypeFloat,
// Status ItemStatusEnabled) are pointers, so they can be sent explicitly and are not overwritten
// unintentionally on update. Use Ptr to set them.
type Item struct {
	ItemID      string         `json:"itemid,omitempty"`      // read-only; required for update
	HostID      string         `json:"hostid,omitempty"`      // required for create
	InterfaceID string         `json:"interfaceid,omitempty"` // required for agent, SNMP, IPMI, JMX, simple check items etc.
	Name        string         `json:"name,omitempty"`        // required for create
	Key         string         `json:"key_,omitempty"`        // required for create
	Type        *ItemType      `json:"type,omitempty"`        // required for create; e.g. ItemTypeTrapper
	ValueType   *ItemValueType `json:"value_type,omitempty"`  // required for create; e.g. ItemValueTypeText
	Status      *ItemStatus    `json:"status,omitempty"`      // e.g. ItemStatusEnabled or ItemStatusDisabled

	Delay         string `json:"delay,omitempty"`   // update interval, e.g. "1m"; required for most polled types
	History       string `json:"history,omitempty"` // e.g. "31d"
	Trends        string `json:"trends,omitempty"`  // e.g. "365d"; numeric items only
	Units         string `json:"units,omitempty"`
	Description   string `json:"description,omitempty"`
	ValueMapID    string `json:"valuemapid,omitempty"`
	InventoryLink int    `json:"inventory_link,omitempty"`
	Timeout       string `json:"timeout,omitempty"`
	UUID          string `json:"uuid,omitempty"` // template items only

	MasterItemID string `json:"master_itemid,omitempty"` // required for ItemTypeDependent
	Params       string `json:"params,omitempty"`        // formula (calculated), script (script/SSH/telnet), SQL (database monitor)
	TrapperHosts string `json:"trapper_hosts,omitempty"` // allowed hosts for trapper / HTTP agent with AllowTraps
	AllowTraps   *int   `json:"allow_traps,omitempty"`   // HTTP agent; 0 no, 1 yes
	SNMPOID      string `json:"snmp_oid,omitempty"`      // required for ItemTypeSNMPAgent
	AuthType     int    `json:"authtype,omitempty"`      // e.g. ItemHTTPAuthBasic or ItemSSHAuthPublicKey
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`

	// HTTP agent
	URL             string          `json:"url,omitempty"` // required for ItemTypeHTTPAgent
	QueryFields     []ItemNameValue `json:"query_fields,omitempty"`
	Headers         []ItemNameValue `json:"headers,omitempty"`
	Posts           string          `json:"posts,omitempty"`
	PostType        int             `json:"post_type,omitempty"`      // e.g. ItemPostTypeJSON
	RequestMethod   int             `json:"request_method,omitempty"` // e.g. ItemRequestMethodPOST
	RetrieveMode    int             `json:"retrieve_mode,omitempty"`  // e.g. ItemRetrieveModeHeaders
	OutputFormat    int             `json:"output_format,omitempty"`  // e.g. ItemOutputFormatJSON
	StatusCodes     string          `json:"status_codes,omitempty"`   // e.g. "200,201-210"
	FollowRedirects *int            `json:"follow_redirects,omitempty"`
	HTTPProxy       string          `json:"http_proxy,omitempty"`
	VerifyHost      *int            `json:"verify_host,omitempty"`
	VerifyPeer      *int            `json:"verify_peer,omitempty"`
	SSLCertFile     string          `json:"ssl_cert_file,omitempty"`
	SSLKeyFile      string          `json:"ssl_key_file,omitempty"`
	SSLKeyPassword  string          `json:"ssl_key_password,omitempty"`

	// Script
	Parameters []ItemNameValue `json:"parameters,omitempty"`

	Tags          []Tag               `json:"tags,omitempty"`
	Preprocessing []ItemPreprocessing `json:"preprocessing,omitempty"`

	// Read-only
	Flags      int    `json:"flags,omitempty"`
	State      int    `json:"state,omitempty"` // e.g. ItemStateNotSupported
	Error      string `json:"error,omitempty"`
	TemplateID string `json:"templateid,omitempty"`
	LastClock  int64  `json:"lastclock,omitempty"` // unix ts
	LastNs     int    `json:"lastns,omitempty"`
	LastValue  string `json:"lastvalue,omitempty"`
	PrevValue  string `json:"prevvalue,omitempty"`
}

// ItemNameValue is used for HTTP agent headers and query fields, and script item parameters.
type ItemNameValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ItemPreprocessing represents an item preprocessing step.
type ItemPreprocessing struct {
	Type               ItemPreprocessingType         `json:"type"`                           // required
	Params             string                        `json:"params,omitempty"`               // multiple params separated by "\n"
	ErrorHandler       ItemPreprocessingErrorHandler `json:"error_handler,omitempty"`        // e.g. ItemPreprocessingErrorDiscard
	ErrorHandlerParams string                        `json:"error_handler_params,omitempty"` // for SetValue / SetError
}

// ItemGetTag is a tag filter for item.get.
type ItemGetTag struct {
	Tag      string `json:"tag"`                // tag name
	Value    string `json:"value,omitempty"`    // tag value
	Operator *int   `json:"operator,omitempty"` // 0..5 per docs
}

// Common API return payload for item.create/update/delete.
type ItemIDsResult struct {
	ItemIDs []string `json:"itemids"`
}

type ItemGetParams struct {
	GetParameters

	ItemIDs      []string `json:"itemids,omitempty"`
	GroupIDs     []string `json:"groupids,omitempty"`
	TemplateIDs  []string `json:"templateids,omitempty"`
	HostIDs      []string `json:"hostids,omitempty"`
	ProxyIDs     []string `json:"proxyids,omitempty"`
	InterfaceIDs []string `json:"interfaceids,omitempty"`
	GraphIDs     []string `json:"graphids,omitempty"`
	TriggerIDs   []string `json:"triggerids,omitempty"`

	WebItems     *bool  `json:"webitems,omitempty"`  // include web items
	Inherited    *bool  `json:"inherited,omitempty"` // true=only inherited from a template
	Templated    *bool  `json:"templated,omitempty"` // true=only items on templates
	Monitored    *bool  `json:"monitored,omitempty"` // true=only enabled items on monitored hosts
	WithTriggers *bool  `json:"with_triggers,omitempty"`
	Group        string `json:"group,omitempty"` // host group name
	Host         string `json:"host,omitempty"`  // technical host name

	// Tag search rules and tags
	EvalType *int         `json:"evaltype,omitempty"` // 0 And/Or (default), 2 Or
	Tags     []ItemGetTag `json:"tags,omitempty"`

	SelectInterfaces    any `json:"selectInterfaces,omitempty"`
	SelectTags          any `json:"selectTags,omitempty"`
	SelectPreprocessing any `json:"selectPreprocessing,omitempty"`
	SelectValueMap      any `json:"selectValueMap,omitempty"`
	SelectDiscoveryRule any `json:"selectDiscoveryRule,omitempty"`
	SelectItemDiscovery any `json:"selectItemDiscovery,omitempty"`
}

type ItemCreateParams struct {
	Item
}

type ItemUpdateParams struct {
	Item
}

type ItemDeleteParams []string

type ItemCreateResponse = ItemIDsResult

type ItemUpdateResponse = ItemIDsResult

type ItemDeleteResponse = ItemIDsResult

func (z *zabbixClient) ItemGet(ctx context.Context, params ItemGetParams) ([]Item, error) {

	var result []Item

	err := z.makeRequest(ctx, "item.get", params, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (z *zabbixClient) ItemCreate(ctx context.Context, params ItemCreateParams) (*ItemCreateResponse, error) {

	var result ItemCreateResponse

	err := z.makeRequest(ctx, "item.create", params, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (z *zabbixClient) ItemUpdate(ctx context.Context, params ItemUpdateParams) (*ItemUpdateResponse, error) {

	var result ItemUpdateResponse

	err := z.makeRequest(ctx, "item.update", params, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (z *zabbixClient) ItemDelete(ctx context.Context, params ItemDeleteParams) (*ItemDeleteResponse, error) {

	var result ItemDeleteResponse

	err := z.makeRequest(ctx, "item.delete", params, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
