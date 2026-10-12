package zabbix

import "context"

// History represents one value returned by history.get.
//
// The fields returned depend on the requested history type: all types return ItemID, Clock, Ns
// and Value, log history (ItemValueTypeLog) also returns ID, Timestamp, Source, Severity and
// LogEventID. Binary history values are base64 encoded.
type History struct {
	ItemID string `json:"itemid,omitempty"`
	Clock  string `json:"clock,omitempty"` // unix ts
	Ns     string `json:"ns,omitempty"`
	Value  string `json:"value,omitempty"`

	// Log history only
	ID         string `json:"id,omitempty"`
	Timestamp  string `json:"timestamp,omitempty"` // log entry time, unix ts
	Source     string `json:"source,omitempty"`
	Severity   string `json:"severity,omitempty"`
	LogEventID string `json:"logeventid,omitempty"`
}

type HistoryGetParams struct {
	GetParameters

	History *ItemValueType `json:"history,omitempty"` // value type to return; Zabbix default is ItemValueTypeUnsigned
	HostIDs []string       `json:"hostids,omitempty"`
	ItemIDs []string       `json:"itemids,omitempty"`

	TimeFrom *int64 `json:"time_from,omitempty"` // Unix timestamp (seconds)
	TimeTill *int64 `json:"time_till,omitempty"` // Unix timestamp (seconds)
}

// HistoryClearParams is a list of item IDs whose history should be cleared.
type HistoryClearParams []string

type HistoryClearResponse struct {
	ItemIDs []string `json:"itemids"`
}

// HistoryPushValue is a single value sent with history.push. Identify the item either by
// ItemID or by Host and Key. Only trapper and HTTP agent items with AllowTraps accept pushed values.
type HistoryPushValue struct {
	ItemID string `json:"itemid,omitempty"`
	Host   string `json:"host,omitempty"`  // technical host name; used together with Key
	Key    string `json:"key,omitempty"`   // item key; used together with Host
	Value  string `json:"value"`           // required
	Clock  int64  `json:"clock,omitempty"` // unix ts; defaults to the current time
	Ns     int    `json:"ns,omitempty"`
}

type HistoryPushParams []HistoryPushValue

// HistoryPushResult is the per-value result of history.push, in the same order as the pushed values.
type HistoryPushResult struct {
	ItemID string `json:"itemid,omitempty"`
	Error  string `json:"error,omitempty"` // set if the value was rejected
}

type HistoryPushResponse struct {
	Response string              `json:"response"` // "success"
	Data     []HistoryPushResult `json:"data"`
}

func (z *zabbixClient) HistoryGet(ctx context.Context, params HistoryGetParams) ([]History, error) {

	var result []History

	err := z.makeRequest(ctx, "history.get", params, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (z *zabbixClient) HistoryClear(ctx context.Context, params HistoryClearParams) (*HistoryClearResponse, error) {

	var result HistoryClearResponse

	err := z.makeRequest(ctx, "history.clear", params, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (z *zabbixClient) HistoryPush(ctx context.Context, params HistoryPushParams) (*HistoryPushResponse, error) {

	var result HistoryPushResponse

	err := z.makeRequest(ctx, "history.push", params, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
