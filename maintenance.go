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
	MaintenanceType int    `json:"maintenance_type,omitempty"` // e.g. MaintenanceTypeWithData or MaintenanceTypeNoData
	TagsEvalType    int    `json:"tags_evaltype,omitempty"`    // e.g. MaintenanceTagsEvalTypeAndOr or MaintenanceTagsEvalTypeOr

	HostGroups  []HostGroupRef          `json:"groups,omitempty"`      // required if hosts not set
	Hosts       []HostRef               `json:"hosts,omitempty"`       // required if groups not set
	TimePeriods []MaintenanceTimePeriod `json:"timeperiods,omitempty"` // required for create
	Tags        []MaintenanceProblemTag `json:"tags,omitempty"`        // only for MaintenanceTypeWithData
}

type MaintenanceTimePeriod struct {
	Period         int   `json:"period,omitempty"`          // seconds; default MaintenanceDefaultPeriodSeconds
	TimePeriodType int   `json:"timeperiod_type,omitempty"` // e.g. MaintenanceTimePeriodOneTime, Daily, Weekly, Monthly
	StartDate      int64 `json:"start_date,omitempty"`      // only for one time only; unix ts
	StartTime      int   `json:"start_time,omitempty"`      // seconds since 00:00; daily/weekly/monthly

	Every     int `json:"every,omitempty"`     // monthly week selector (MaintenanceWeekFirst..MaintenanceWeekLast) or interval
	DayOfWeek int `json:"dayofweek,omitempty"` // day bitmask (MaintenanceDayMonday..MaintenanceDaySunday)
	Day       int `json:"day,omitempty"`       // day of month; required for monthly if dayofweek not set
	Month     int `json:"month,omitempty"`     // month bitmask (MaintenanceMonthJanuary..MaintenanceMonthDecember)
}

// MaintenanceProblemTag represents the Zabbix maintenance "problem tag" object.
type MaintenanceProblemTag struct {
	Tag      string `json:"tag"`                // required
	Operator int    `json:"operator,omitempty"` // e.g. MaintenanceTagOperatorEquals or MaintenanceTagOperatorContains
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
