package zabbix_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	zabbix "github.com/nimok/nim-go-zabbix"
)

// pushTestHistory pushes values until the server accepts them. A newly created trapper item only
// accepts values once the server has loaded it into its configuration cache. history.push reports
// acceptance per value, so only rejected values are retried to avoid storing duplicates.
// The returned results are in the same order as values.
func pushTestHistory(t *testing.T, ctx context.Context, client zabbix.Client, values zabbix.HistoryPushParams) []zabbix.HistoryPushResult {
	t.Helper()

	results := make([]zabbix.HistoryPushResult, len(values))
	pending := make([]int, len(values)) // indexes into values not yet accepted
	for i := range values {
		pending[i] = i
	}

	deadline := time.Now().Add(90 * time.Second)
	for {
		batch := make(zabbix.HistoryPushParams, len(pending))
		for i, idx := range pending {
			batch[i] = values[idx]
		}

		resp, err := client.HistoryPush(ctx, batch)
		if err != nil {
			t.Log("History push failed:", err)
			t.FailNow()
		}

		if resp.Response != "success" || len(resp.Data) != len(batch) {
			t.Logf("Unexpected push response: %+v", resp)
			t.FailNow()
		}

		var rejected []int
		var lastErr string
		for i, d := range resp.Data {
			results[pending[i]] = d
			if d.Error != "" {
				rejected = append(rejected, pending[i])
				lastErr = d.Error
			}
		}
		if len(rejected) == 0 {
			return results
		}

		if time.Now().After(deadline) {
			t.Log("History push rejected:", lastErr)
			t.FailNow()
		}
		pending = rejected
		time.Sleep(2 * time.Second)
	}
}

// waitForHistory polls history.get until the item has at least want values of the given type.
func waitForHistory(t *testing.T, ctx context.Context, client zabbix.Client, itemID string, valueType zabbix.ItemValueType, want int) []zabbix.History {
	t.Helper()

	deadline := time.Now().Add(30 * time.Second)
	for {
		history, err := client.HistoryGet(ctx, zabbix.HistoryGetParams{
			GetParameters: zabbix.GetParameters{
				Output:    "extend",
				Sortfield: []string{"clock"},
			},
			History: zabbix.Ptr(valueType),
			ItemIDs: []string{itemID},
		})
		if err != nil {
			t.Log("History get failed:", err)
			t.FailNow()
		}

		if len(history) >= want {
			return history
		}

		if time.Now().After(deadline) {
			t.Logf("Expected %d history values, got %d", want, len(history))
			t.FailNow()
		}
		time.Sleep(time.Second)
	}
}

func TestGetHistory(t *testing.T) {

	ctx := context.Background()

	client, err := zabbix.NewClient(url, zabbix.WithUserPass(user, passwd))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if err := client.Authenticate(); err != nil {
		t.Log("Initial auth failed:", err)
		t.FailNow()
	}

	// Push our own float values instead of relying on data collected by the server,
	// which a freshly started instance may not have yet.
	suffix := time.Now().Format("20060102150405.000000")
	itemResp, err := client.ItemCreate(ctx, zabbix.ItemCreateParams{
		zabbix.Item{
			HostID:    "10084",
			Name:      "Test History Get " + suffix,
			Key:       "test.history.float[" + suffix + "]",
			Type:      zabbix.Ptr(zabbix.ItemTypeTrapper),
			ValueType: zabbix.Ptr(zabbix.ItemValueTypeFloat),
		},
	})
	if err != nil || len(itemResp.ItemIDs) == 0 {
		t.Log("Item create failed:", err)
		t.FailNow()
	}
	itemID := itemResp.ItemIDs[0]

	defer func() {
		_, cleanupErr := client.ItemDelete(ctx, zabbix.ItemDeleteParams{itemID})
		if cleanupErr != nil {
			t.Log("Item cleanup failed:", cleanupErr)
		}
	}()

	now := time.Now().Unix()
	pushTestHistory(t, ctx, client, zabbix.HistoryPushParams{
		{ItemID: itemID, Value: "1.5", Clock: now - 120},
		{ItemID: itemID, Value: "2.5", Clock: now},
	})
	waitForHistory(t, ctx, client, itemID, zabbix.ItemValueTypeFloat, 2)

	// Only the newest value is within the time range.
	history, err := client.HistoryGet(ctx, zabbix.HistoryGetParams{
		GetParameters: zabbix.GetParameters{
			Output:    "extend",
			Sortfield: []string{"clock"},
			Sortorder: "DESC",
		},
		History:  zabbix.Ptr(zabbix.ItemValueTypeFloat),
		HostIDs:  []string{"10084"},
		ItemIDs:  []string{itemID},
		TimeFrom: zabbix.Ptr(now - 60),
	})
	if err != nil {
		t.Log("History get failed:", err)
		t.FailNow()
	}

	if len(history) != 1 {
		t.Logf("Expected 1 history value, got %d: %+v", len(history), history)
		t.FailNow()
	}

	value, err := strconv.ParseFloat(history[0].Value, 64)
	if history[0].ItemID != itemID || history[0].Clock != strconv.FormatInt(now, 10) || err != nil || value != 2.5 {
		t.Logf("History not decoded: got=%+v", history[0])
		t.FailNow()
	}
}

func TestPushHistory(t *testing.T) {

	ctx := context.Background()

	client, err := zabbix.NewClient(url, zabbix.WithUserPass(user, passwd))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if err := client.Authenticate(); err != nil {
		t.Log("Initial auth failed:", err)
		t.FailNow()
	}

	itemID := createTestItem(t, ctx, client, "Test History Push")

	defer func() {
		_, cleanupErr := client.ItemDelete(ctx, zabbix.ItemDeleteParams{itemID})
		if cleanupErr != nil {
			t.Log("Item cleanup failed:", cleanupErr)
		}
	}()

	items, err := client.ItemGet(ctx, zabbix.ItemGetParams{
		GetParameters: zabbix.GetParameters{
			Output: []string{"key_"},
		},
		ItemIDs: []string{itemID},
	})
	if err != nil || len(items) != 1 {
		t.Log("Item get failed:", err)
		t.FailNow()
	}

	now := time.Now().Unix()
	results := pushTestHistory(t, ctx, client, zabbix.HistoryPushParams{
		{ItemID: itemID, Value: "first", Clock: now - 1},
		{Host: "Zabbix server", Key: items[0].Key, Value: "second", Clock: now},
	})

	for _, d := range results {
		if d.ItemID != itemID {
			t.Logf("Push result item mismatch: got=%q want=%q", d.ItemID, itemID)
			t.FailNow()
		}
	}

	history := waitForHistory(t, ctx, client, itemID, zabbix.ItemValueTypeText, 2)
	if history[0].Value != "first" || history[1].Value != "second" {
		t.Logf("History values mismatch: got=%+v", history)
		t.FailNow()
	}
}

func TestClearHistory(t *testing.T) {

	ctx := context.Background()

	client, err := zabbix.NewClient(url, zabbix.WithUserPass(user, passwd))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if err := client.Authenticate(); err != nil {
		t.Log("Initial auth failed:", err)
		t.FailNow()
	}

	itemID := createTestItem(t, ctx, client, "Test History Clear")

	defer func() {
		_, cleanupErr := client.ItemDelete(ctx, zabbix.ItemDeleteParams{itemID})
		if cleanupErr != nil {
			t.Log("Item cleanup failed:", cleanupErr)
		}
	}()

	pushTestHistory(t, ctx, client, zabbix.HistoryPushParams{
		{ItemID: itemID, Value: "to be cleared"},
	})
	waitForHistory(t, ctx, client, itemID, zabbix.ItemValueTypeText, 1)

	clearResp, err := client.HistoryClear(ctx, zabbix.HistoryClearParams{itemID})
	if err != nil {
		t.Log("History clear failed:", err)
		t.FailNow()
	}

	if len(clearResp.ItemIDs) != 1 || clearResp.ItemIDs[0] != itemID {
		t.Logf("History clear item mismatch: got=%v want=%s", clearResp.ItemIDs, itemID)
		t.FailNow()
	}

	history, err := client.HistoryGet(ctx, zabbix.HistoryGetParams{
		GetParameters: zabbix.GetParameters{
			Output: "extend",
		},
		History: zabbix.Ptr(zabbix.ItemValueTypeText),
		ItemIDs: []string{itemID},
	})
	if err != nil {
		t.Log("History get failed:", err)
		t.FailNow()
	}

	if len(history) != 0 {
		t.Logf("Expected no history after clear, got %d", len(history))
		t.FailNow()
	}
}
