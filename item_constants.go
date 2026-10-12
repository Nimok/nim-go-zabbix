package zabbix

// ItemType represents the item type (`type`).
type ItemType int

const (
	ItemTypeZabbixAgent       ItemType = 0
	ItemTypeTrapper           ItemType = 2
	ItemTypeSimpleCheck       ItemType = 3
	ItemTypeInternal          ItemType = 5
	ItemTypeZabbixAgentActive ItemType = 7
	ItemTypeWebItem           ItemType = 9 // read-only; created by web scenarios
	ItemTypeExternalCheck     ItemType = 10
	ItemTypeDatabaseMonitor   ItemType = 11
	ItemTypeIPMIAgent         ItemType = 12
	ItemTypeSSHAgent          ItemType = 13
	ItemTypeTelnetAgent       ItemType = 14
	ItemTypeCalculated        ItemType = 15
	ItemTypeJMXAgent          ItemType = 16
	ItemTypeSNMPTrap          ItemType = 17
	ItemTypeDependent         ItemType = 18
	ItemTypeHTTPAgent         ItemType = 19
	ItemTypeSNMPAgent         ItemType = 20
	ItemTypeScript            ItemType = 21
	ItemTypeBrowser           ItemType = 22
)

// ItemValueType represents the type of information of the item (`value_type`).
type ItemValueType int

const (
	ItemValueTypeFloat    ItemValueType = 0
	ItemValueTypeChar     ItemValueType = 1
	ItemValueTypeLog      ItemValueType = 2
	ItemValueTypeUnsigned ItemValueType = 3
	ItemValueTypeText     ItemValueType = 4
	ItemValueTypeBinary   ItemValueType = 5
)

// ItemStatus represents whether the item is enabled (`status`).
type ItemStatus int

const (
	ItemStatusEnabled  ItemStatus = 0
	ItemStatusDisabled ItemStatus = 1
)

// Item state values (`state`, read-only).
const (
	ItemStateNormal       = 0
	ItemStateNotSupported = 1
)

// HTTP agent request method values (`request_method`).
const (
	ItemRequestMethodGET  = 0
	ItemRequestMethodPOST = 1
	ItemRequestMethodPUT  = 2
	ItemRequestMethodHEAD = 3
)

// HTTP agent post body type values (`post_type`).
const (
	ItemPostTypeRaw  = 0
	ItemPostTypeJSON = 2
	ItemPostTypeXML  = 3
)

// HTTP agent retrieve mode values (`retrieve_mode`).
const (
	ItemRetrieveModeBody    = 0
	ItemRetrieveModeHeaders = 1
	ItemRetrieveModeBoth    = 2
)

// HTTP agent output format values (`output_format`).
const (
	ItemOutputFormatRaw  = 0
	ItemOutputFormatJSON = 1
)

// HTTP agent authentication method values (`authtype`).
const (
	ItemHTTPAuthNone     = 0
	ItemHTTPAuthBasic    = 1
	ItemHTTPAuthNTLM     = 2
	ItemHTTPAuthKerberos = 3
	ItemHTTPAuthDigest   = 4
)

// SSH agent authentication method values (`authtype`).
const (
	ItemSSHAuthPassword  = 0
	ItemSSHAuthPublicKey = 1
)

// ItemPreprocessingType represents the preprocessing step type (`type`).
type ItemPreprocessingType int

const (
	ItemPreprocessingCustomMultiplier          ItemPreprocessingType = 1
	ItemPreprocessingRightTrim                 ItemPreprocessingType = 2
	ItemPreprocessingLeftTrim                  ItemPreprocessingType = 3
	ItemPreprocessingTrim                      ItemPreprocessingType = 4
	ItemPreprocessingRegex                     ItemPreprocessingType = 5
	ItemPreprocessingBoolToDecimal             ItemPreprocessingType = 6
	ItemPreprocessingOctalToDecimal            ItemPreprocessingType = 7
	ItemPreprocessingHexToDecimal              ItemPreprocessingType = 8
	ItemPreprocessingSimpleChange              ItemPreprocessingType = 9
	ItemPreprocessingChangePerSecond           ItemPreprocessingType = 10
	ItemPreprocessingXMLXPath                  ItemPreprocessingType = 11
	ItemPreprocessingJSONPath                  ItemPreprocessingType = 12
	ItemPreprocessingInRange                   ItemPreprocessingType = 13
	ItemPreprocessingMatchesRegex              ItemPreprocessingType = 14
	ItemPreprocessingNotMatchesRegex           ItemPreprocessingType = 15
	ItemPreprocessingCheckJSONError            ItemPreprocessingType = 16
	ItemPreprocessingCheckXMLError             ItemPreprocessingType = 17
	ItemPreprocessingCheckRegexError           ItemPreprocessingType = 18
	ItemPreprocessingDiscardUnchanged          ItemPreprocessingType = 19
	ItemPreprocessingDiscardUnchangedHeartbeat ItemPreprocessingType = 20
	ItemPreprocessingJavaScript                ItemPreprocessingType = 21
	ItemPreprocessingPrometheusPattern         ItemPreprocessingType = 22
	ItemPreprocessingPrometheusToJSON          ItemPreprocessingType = 23
	ItemPreprocessingCSVToJSON                 ItemPreprocessingType = 24
	ItemPreprocessingReplace                   ItemPreprocessingType = 25
	ItemPreprocessingCheckUnsupported          ItemPreprocessingType = 26
	ItemPreprocessingXMLToJSON                 ItemPreprocessingType = 27
	ItemPreprocessingSNMPWalkValue             ItemPreprocessingType = 28
	ItemPreprocessingSNMPWalkToJSON            ItemPreprocessingType = 29
	ItemPreprocessingSNMPGetValue              ItemPreprocessingType = 30
)

// ItemPreprocessingErrorHandler represents the action taken when a preprocessing step fails (`error_handler`).
type ItemPreprocessingErrorHandler int

const (
	ItemPreprocessingErrorDefault  ItemPreprocessingErrorHandler = 0 // set item state to "not supported"
	ItemPreprocessingErrorDiscard  ItemPreprocessingErrorHandler = 1
	ItemPreprocessingErrorSetValue ItemPreprocessingErrorHandler = 2
	ItemPreprocessingErrorSetError ItemPreprocessingErrorHandler = 3
)
