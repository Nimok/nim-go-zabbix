package zabbix

import "context"

// Trend represents one hourly aggregate returned by trend.get.
//
// Trends are only stored for numeric items (ItemValueTypeFloat and ItemValueTypeUnsigned).
type Trend struct {
	ItemID   string `json:"itemid,omitempty"`
	Clock    string `json:"clock,omitempty"` // start of the hour, unix ts
	Num      string `json:"num,omitempty"`   // number of values aggregated in the hour
	ValueMin string `json:"value_min,omitempty"`
	ValueAvg string `json:"value_avg,omitempty"`
	ValueMax string `json:"value_max,omitempty"`
}

// TrendGetParams are the parameters for trend.get. Of the common get parameters trend.get only
// supports output, countOutput and limit, so they are listed here instead of embedding GetParameters.
type TrendGetParams struct {
	Output      any  `json:"output,omitempty"`
	CountOutput bool `json:"countOutput,omitempty"`
	Limit       int  `json:"limit,omitempty"`

	ItemIDs []string `json:"itemids,omitempty"`

	TimeFrom *int64 `json:"time_from,omitempty"` // Unix timestamp (seconds)
	TimeTill *int64 `json:"time_till,omitempty"` // Unix timestamp (seconds)
}

func (z *zabbixClient) TrendGet(ctx context.Context, params TrendGetParams) ([]Trend, error) {

	var result []Trend

	err := z.makeRequest(ctx, "trend.get", params, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
