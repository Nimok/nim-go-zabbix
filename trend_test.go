package zabbix_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	zabbix "github.com/nimok/nim-go-zabbix"
)

// waitForTrends polls trend.get until the item has at least want trends up to timeTill.
func waitForTrends(t *testing.T, ctx context.Context, client zabbix.Client, itemID string, timeTill int64, want int) []zabbix.Trend {
	t.Helper()

	deadline := time.Now().Add(60 * time.Second)
	for {
		trends, err := client.TrendGet(ctx, zabbix.TrendGetParams{
			Output:   "extend",
			ItemIDs:  []string{itemID},
			TimeTill: zabbix.Ptr(timeTill),
		})
		if err != nil {
			t.Log("Trend get failed:", err)
			t.FailNow()
		}

		if len(trends) >= want {
			return trends
		}

		if time.Now().After(deadline) {
			t.Logf("Expected %d trends, got %d", want, len(trends))
			t.FailNow()
		}
		time.Sleep(time.Second)
	}
}

func TestGetTrend(t *testing.T) {

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

	suffix := time.Now().Format("20060102150405.000000")
	itemResp, err := client.ItemCreate(ctx, zabbix.ItemCreateParams{
		zabbix.Item{
			HostID:    "10084",
			Name:      "Test Trend Get " + suffix,
			Key:       "test.trend.float[" + suffix + "]",
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

	// The server writes an hour's trend once it receives a value for a later hour of the item,
	// so push values for two past hours followed by a value for the current hour.
	now := time.Now().Unix()
	hourStart := now - now%3600
	firstHour := hourStart - 2*3600
	secondHour := hourStart - 3600
	pushTestHistory(t, ctx, client, zabbix.HistoryPushParams{
		{ItemID: itemID, Value: "1", Clock: firstHour + 10},
		{ItemID: itemID, Value: "3", Clock: firstHour + 20},
		{ItemID: itemID, Value: "10", Clock: secondHour + 10},
		{ItemID: itemID, Value: "20", Clock: now},
	})

	trends := waitForTrends(t, ctx, client, itemID, hourStart-1, 2)
	if len(trends) != 2 {
		t.Logf("Expected 2 trends, got %d: %+v", len(trends), trends)
		t.FailNow()
	}

	var first *zabbix.Trend
	for i := range trends {
		if trends[i].Clock == strconv.FormatInt(firstHour, 10) {
			first = &trends[i]
		}
	}
	if first == nil {
		t.Logf("No trend for hour %d: %+v", firstHour, trends)
		t.FailNow()
	}

	minValue, errMin := strconv.ParseFloat(first.ValueMin, 64)
	avgValue, errAvg := strconv.ParseFloat(first.ValueAvg, 64)
	maxValue, errMax := strconv.ParseFloat(first.ValueMax, 64)
	if first.ItemID != itemID || first.Num != "2" || errMin != nil || errAvg != nil || errMax != nil ||
		minValue != 1 || avgValue != 2 || maxValue != 3 {
		t.Logf("Trend not decoded: got=%+v", *first)
		t.FailNow()
	}

	// Only the second hour is within the time range.
	filtered, err := client.TrendGet(ctx, zabbix.TrendGetParams{
		Output:   []string{"itemid", "clock", "num"},
		ItemIDs:  []string{itemID},
		TimeFrom: zabbix.Ptr(secondHour),
		TimeTill: zabbix.Ptr(hourStart - 1),
	})
	if err != nil {
		t.Log("Trend get failed:", err)
		t.FailNow()
	}

	if len(filtered) != 1 || filtered[0].Clock != strconv.FormatInt(secondHour, 10) || filtered[0].Num != "1" {
		t.Logf("Unexpected filtered trends: %+v", filtered)
		t.FailNow()
	}

	limited, err := client.TrendGet(ctx, zabbix.TrendGetParams{
		Output:   "extend",
		Limit:    1,
		ItemIDs:  []string{itemID},
		TimeTill: zabbix.Ptr(hourStart - 1),
	})
	if err != nil {
		t.Log("Trend get failed:", err)
		t.FailNow()
	}

	if len(limited) != 1 {
		t.Logf("Expected 1 trend with limit, got %d", len(limited))
		t.FailNow()
	}
}
