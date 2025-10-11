package zabbix

import "context"

type ProblemGetTag struct {
	Tag      string `json:"tag"`                // tag name (exact match)
	Value    string `json:"value,omitempty"`    // tag value (operator-dependent match)
	Operator *int   `json:"operator,omitempty"` // 0..5 per docs
}

type ProblemGetParams struct {
	GetParameters

	EventIDs   []string `json:"eventids,omitempty"`
	GroupIDs   []string `json:"groupids,omitempty"`
	HostIDs    []string `json:"hostids,omitempty"`
	ObjectIDs  []string `json:"objectids,omitempty"`
	ActionUser []string `json:"action_userids,omitempty"`

	// Basic filters
	Source       *int  `json:"source,omitempty"`       // default 0 (trigger)
	Object       *int  `json:"object,omitempty"`       // default 0 (trigger)
	Acknowledged *bool `json:"acknowledged,omitempty"` // true=only acked, false=only unacked
	Action       *int  `json:"action,omitempty"`       // bitmap of event update actions
	Suppressed   *bool `json:"suppressed,omitempty"`   // true=only suppressed
	Symptom      *bool `json:"symptom,omitempty"`      // true=symptom, false=cause
	Severities   []int `json:"severities,omitempty"`   // applies only if object=trigger

	// Tag search rules and tags
	EvalType *int            `json:"evaltype,omitempty"` // 0 And/Or (default), 2 Or
	Tags     []ProblemGetTag `json:"tags,omitempty"`

	// Time / range filters
	Recent      *bool  `json:"recent,omitempty"`       // include recently RESOLVED
	EventIDFrom string `json:"eventid_from,omitempty"` // >= given ID
	EventIDTill string `json:"eventid_till,omitempty"` // <= given ID
	TimeFrom    *int64 `json:"time_from,omitempty"`    // Unix timestamp (seconds)
	TimeTill    *int64 `json:"time_till,omitempty"`    // Unix timestamp (seconds)

	// Select/expand related data (query type: "extend", "count", or []string)
	SelectAcknowledges    any `json:"selectAcknowledges,omitempty"`    // e.g. "extend" or []string
	SelectTags            any `json:"selectTags,omitempty"`            // e.g. "extend" or []string
	SelectSuppressionData any `json:"selectSuppressionData,omitempty"` // e.g. "extend" or []string

}

// Problem represents one entry returned by problem.get.
type Problem struct {
	// Core problem fields (IDs and timestamps are strings in Zabbix JSON)
	EventID       string `json:"eventid"`       // ID
	Source        string `json:"source"`        // "0" trigger, "3" internal, "4" service-status update
	Object        string `json:"object"`        // depends on source
	ObjectID      string `json:"objectid"`      // related object ID
	Clock         string `json:"clock"`         // timestamp (Unix seconds, as string)
	Ns            string `json:"ns"`            // creation nanoseconds (as string)
	REventID      string `json:"r_eventid"`     // recovery event ID
	RClock        string `json:"r_clock"`       // recovery time (Unix seconds, as string)
	RNs           string `json:"r_ns"`          // recovery nanoseconds (as string)
	CauseEventID  string `json:"cause_eventid"` // ID of the cause event
	CorrelationID string `json:"correlationid"` // correlation rule ID (if recovered by rule)
	UserID        string `json:"userid"`        // user who manually closed the problem (if any)
	Name          string `json:"name"`          // resolved problem name (may be empty for unresolved)
	Acknowledged  string `json:"acknowledged"`  // "0" or "1"
	Severity      string `json:"severity"`      // "0".."5"
	Suppressed    string `json:"suppressed"`    // "0" or "1"
	OpData        string `json:"opdata"`        // operational data with expanded macros

	// Added when requested:
	URLs            []ProblemMediaURL       `json:"urls,omitempty"`             // media-type URLs (active only)
	Acknowledges    []ProblemAcknowledge    `json:"acknowledges,omitempty"`     // selectAcknowledges
	Tags            []ProblemTag            `json:"tags,omitempty"`             // selectTags
	SuppressionData []ProblemSuppressionRef `json:"suppression_data,omitempty"` // selectSuppressionData
}

// ProblemMediaURL corresponds to entries in Problem.URLs.
type ProblemMediaURL struct {
	Name string `json:"name"` // defined URL name
	URL  string `json:"url"`  // actual URL
}

// ProblemAcknowledge is returned when selectAcknowledges is used.
type ProblemAcknowledge struct {
	AcknowledgeID string `json:"acknowledgeid"`  // ID
	UserID        string `json:"userid"`         // who updated
	EventID       string `json:"eventid"`        // event updated
	Clock         string `json:"clock"`          // when (Unix seconds, as string)
	Message       string `json:"message"`        // update message
	Action        string `json:"action"`         // update action type
	OldSeverity   string `json:"old_severity"`   // previous severity
	NewSeverity   string `json:"new_severity"`   // new severity
	SuppressUntil string `json:"suppress_until"` // suppression end time (Unix seconds, as string)
	TaskID        string `json:"taskid"`         // task id (rank change)
}

// ProblemTag is returned when selectTags is used.
type ProblemTag struct {
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

// ProblemSuppressionRef is returned when selectSuppressionData is used.
type ProblemSuppressionRef struct {
	MaintenanceID string `json:"maintenanceid"`  // active maintenance ID (if any)
	UserID        string `json:"userid"`         // who suppressed the problem
	SuppressUntil string `json:"suppress_until"` // suppression end time (Unix seconds, as string)
}

func (z *zabbixClient) ProblemGet(ctx context.Context, params ProblemGetParams) (*[]Problem, error) {

	var result []Problem

	err := z.makeRequest(ctx, "problem.get", params, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
