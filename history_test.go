package zabbix_test

import (
	"context"
	"testing"
	"time"

	zabbix "github.com/nimok/nim-go-zabbix"
)

// pushTestHistory pushes values until the server accepts them. A newly created trapper item only
// accepts values once the server has loaded it into its configuration cache.
func pushTestHistory(t *testing.T, ctx context.Context, client zabbix.Client, values zabbix.HistoryPushParams) *zabbix.HistoryPushResponse {
	t.Helper()

	deadline := time.Now().Add(90 * time.Second)
	for {
		resp, err := client.HistoryPush(ctx, values)
		if err != nil {
			t.Log("History push failed:", err)
			t.FailNow()
		}

		var rejected string
		for _, d := range resp.Data {
			if d.Error != "" {
				rejected = d.Error
			}
		}
		if rejected == "" {
			return resp
		}

		if time.Now().After(deadline) {
			t.Log("History push rejected:", rejected)
			t.FailNow()
		}
		time.Sleep(2 * time.Second)
	}
}

// waitForHistory polls history.get until the item has at least want text values.
func waitForHistory(t *testing.T, ctx context.Context, client zabbix.Client, itemID string, want int) []zabbix.History {
	t.Helper()

	deadline := time.Now().Add(30 * time.Second)
	for {
		history, err := client.HistoryGet(ctx, zabbix.HistoryGetParams{
			GetParameters: zabbix.GetParameters{
				Output:    "extend",
				Sortfield: []string{"clock"},
			},
			History: zabbix.Ptr(zabbix.ItemValueTypeText),
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

	history, err := client.HistoryGet(ctx, zabbix.HistoryGetParams{
		GetParameters: zabbix.GetParameters{
			Output:    "extend",
			Limit:     10,
			Sortfield: []string{"clock"},
			Sortorder: "DESC",
		},
		History: zabbix.Ptr(zabbix.ItemValueTypeFloat),
		HostIDs: []string{"10084"},
	})
	if err != nil {
		t.Log("History get failed:", err)
		t.FailNow()
	}

	if len(history) == 0 {
		t.Log("No float history found on Zabbix server host")
		t.FailNow()
	}

	if history[0].ItemID == "" || history[0].Clock == "" || history[0].Value == "" {
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
	resp := pushTestHistory(t, ctx, client, zabbix.HistoryPushParams{
		{ItemID: itemID, Value: "first", Clock: now - 1},
		{Host: "Zabbix server", Key: items[0].Key, Value: "second", Clock: now},
	})

	if resp.Response != "success" || len(resp.Data) != 2 {
		t.Logf("Unexpected push response: %+v", resp)
		t.FailNow()
	}
	for _, d := range resp.Data {
		if d.ItemID != itemID {
			t.Logf("Push result item mismatch: got=%q want=%q", d.ItemID, itemID)
			t.FailNow()
		}
	}

	history := waitForHistory(t, ctx, client, itemID, 2)
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
	waitForHistory(t, ctx, client, itemID, 1)

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
