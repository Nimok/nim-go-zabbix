package zabbix

import (
	"context"
)

type Maintenance struct {
	MaintenanceID string `json:"maintenanceid,omitempty"` // read-only; required for update
	Name          string `json:"name,omitempty"`          // required for create

	ActiveSince int64 `json:"active_since,omitempty"` // required for create; unix ts
	ActiveTill  int64 `json:"active_till,omitempty"`  // required for create; unix ts

	Description     string `json:"description,omitempty"`
	MaintenanceType int    `json:"maintenance_type,omitempty"` // 0=with data collection, 1=without
	TagsEvalType    int    `json:"tags_evaltype,omitempty"`    // 0=And/Or, 2=Or

	HostGroups  []HostGroupRef          `json:"groups,omitempty"`      // required if hosts not set
	Hosts       []HostRef               `json:"hosts,omitempty"`       // required if groups not set
	TimePeriods []MaintenanceTimePeriod `json:"timeperiods,omitempty"` // required for create
	Tags        []MaintenanceProblemTag `json:"tags,omitempty"`        // only for maintenance_type=0
}

type MaintenanceTimePeriod struct {
	Period         int   `json:"period,omitempty"`          // seconds; default 3600
	TimePeriodType int   `json:"timeperiod_type,omitempty"` // 0=one time, 2=daily, 3=weekly, 4=monthly
	StartDate      int64 `json:"start_date,omitempty"`      // only for one time only; unix ts
	StartTime      int   `json:"start_time,omitempty"`      // seconds since 00:00; daily/weekly/monthly

	Every     int `json:"every,omitempty"`     // interval/week-of-month/day-of-month semantics
	DayOfWeek int `json:"dayofweek,omitempty"` // bitmask; required for weekly (and some monthly)
	Day       int `json:"day,omitempty"`       // day of month; required for monthly if dayofweek not set
	Month     int `json:"month,omitempty"`     // bitmask; required for monthly
}

// MaintenanceProblemTag represents the Zabbix maintenance "problem tag" object.
type MaintenanceProblemTag struct {
	Tag      string `json:"tag"`                // required
	Operator int    `json:"operator,omitempty"` // 0=Equals, 2=Contains (default)
	Value    string `json:"value,omitempty"`    // tag value
}

// HostGroupRef is the minimal object accepted by maintenance.create/update for groups.
type HostGroupRef struct {
	GroupID string `json:"groupid"`
}

// HostRef is the minimal object accepted by maintenance.create/update for hosts.
type HostRef struct {
	HostID string `json:"hostid"`
}

// Common API return payload for maintenance.create/update.
type MaintenanceIDsResult struct {
	MaintenanceIDs []string `json:"maintenanceids"`
}

type MaintenanceGetParams struct {
	GetParameters

	GroupIDs       []string `json:"groupids,omitempty"`
	HostIDs        []string `json:"hostids,omitempty"`
	MaintenanceIDs []string `json:"maintenanceids,omitempty"`

	SelectHostGroups  any `json:"selectHostGroups,omitempty"`
	SelectHosts       any `json:"selectHosts,omitempty"`
	SelectTags        any `json:"selectTags,omitempty"`
	SelectTimeperiods any `json:"selectTimeperiods,omitempty"`
}

type MaintenanceCreateParams struct {
	Maintenance
}

type MaintenanceCreateResponse struct {
	MaintenanceIDs []string `json:"maintenanceids"` // IDs of the created maintenances
}

func (z *zabbixClient) MaintenanceGet(ctx context.Context, params MaintenanceGetParams) (*[]Maintenance, error) {

	var result []Maintenance

	err := z.makeRequest(ctx, "maintenance.get", params, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (z *zabbixClient) MaintenanceCreate(ctx context.Context, params MaintenanceCreateParams) (*MaintenanceCreateResponse, error) {

	var result MaintenanceCreateResponse

	err := z.makeRequest(ctx, "maintenance.create", params, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
