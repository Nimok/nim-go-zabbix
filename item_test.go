package zabbix_test

import (
	"context"
	"testing"
	"time"

	zabbix "github.com/nimok/nim-go-zabbix"
)

// createTestItem creates a trapper item on the Zabbix server host and returns its ID.
func createTestItem(t *testing.T, ctx context.Context, client zabbix.Client, name string) string {
	t.Helper()

	suffix := time.Now().Format("20060102150405.000000")
	resp, err := client.ItemCreate(ctx, zabbix.ItemCreateParams{
		zabbix.Item{
			HostID:      "10084",
			Name:        name + " " + suffix,
			Key:         "test.item[" + suffix + "]",
			Type:        zabbix.Ptr(zabbix.ItemTypeTrapper),
			ValueType:   zabbix.Ptr(zabbix.ItemValueTypeText),
			Description: "before update",
			Tags: []zabbix.Tag{
				{Tag: "component", Value: "test"},
			},
			Preprocessing: []zabbix.ItemPreprocessing{
				{Type: zabbix.ItemPreprocessingTrim, Params: " "},
			},
		},
	})
	if err != nil {
		t.Log("Item create failed:", err)
		t.FailNow()
	}

	if len(resp.ItemIDs) == 0 {
		t.Log("No item id returned by create")
		t.FailNow()
	}

	return resp.ItemIDs[0]
}

func TestGetItem(t *testing.T) {

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

	items, err := client.ItemGet(ctx, zabbix.ItemGetParams{
		GetParameters: zabbix.GetParameters{
			Output: "extend",
			Limit:  10,
		},
		HostIDs:             []string{"10084"},
		SelectTags:          "extend",
		SelectPreprocessing: "extend",
	})
	if err != nil {
		t.Log("Item get failed:", err)
		t.FailNow()
	}

	if len(items) == 0 {
		t.Log("No items found on Zabbix server host")
		t.FailNow()
	}

	if items[0].Type == nil || items[0].ValueType == nil {
		t.Log("Item type/value_type not decoded")
		t.FailNow()
	}
}

func TestCreateItem(t *testing.T) {

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

	itemID := createTestItem(t, ctx, client, "Test Item Create")

	defer func() {
		_, cleanupErr := client.ItemDelete(ctx, zabbix.ItemDeleteParams{itemID})
		if cleanupErr != nil {
			t.Log("Item cleanup failed:", cleanupErr)
		}
	}()

	items, err := client.ItemGet(ctx, zabbix.ItemGetParams{
		GetParameters: zabbix.GetParameters{
			Output: "extend",
		},
		ItemIDs:             []string{itemID},
		SelectTags:          "extend",
		SelectPreprocessing: "extend",
	})
	if err != nil {
		t.Log("Item get failed:", err)
		t.FailNow()
	}

	if len(items) != 1 {
		t.Logf("Expected 1 item, got %d", len(items))
		t.FailNow()
	}

	item := items[0]
	if item.Type == nil || *item.Type != zabbix.ItemTypeTrapper {
		t.Logf("Item type mismatch: got=%v want=%d", item.Type, zabbix.ItemTypeTrapper)
		t.FailNow()
	}
	if len(item.Tags) != 1 || item.Tags[0].Tag != "component" {
		t.Logf("Item tags mismatch: got=%+v", item.Tags)
		t.FailNow()
	}
	if len(item.Preprocessing) != 1 || item.Preprocessing[0].Type != zabbix.ItemPreprocessingTrim {
		t.Logf("Item preprocessing mismatch: got=%+v", item.Preprocessing)
		t.FailNow()
	}
}

func TestCreateDependentItem(t *testing.T) {

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

	masterID := createTestItem(t, ctx, client, "Test Item Master")

	// Deleting the master item also deletes its dependent items.
	defer func() {
		_, cleanupErr := client.ItemDelete(ctx, zabbix.ItemDeleteParams{masterID})
		if cleanupErr != nil {
			t.Log("Item cleanup failed:", cleanupErr)
		}
	}()

	suffix := time.Now().Format("20060102150405.000000")
	resp, err := client.ItemCreate(ctx, zabbix.ItemCreateParams{
		zabbix.Item{
			HostID:       "10084",
			Name:         "Test Item Dependent " + suffix,
			Key:          "test.dependent[" + suffix + "]",
			Type:         zabbix.Ptr(zabbix.ItemTypeDependent),
			ValueType:    zabbix.Ptr(zabbix.ItemValueTypeFloat),
			MasterItemID: masterID,
			Preprocessing: []zabbix.ItemPreprocessing{
				{
					Type:               zabbix.ItemPreprocessingJSONPath,
					Params:             "$.value",
					ErrorHandler:       zabbix.ItemPreprocessingErrorSetValue,
					ErrorHandlerParams: "0",
				},
			},
		},
	})
	if err != nil {
		t.Log("Dependent item create failed:", err)
		t.FailNow()
	}

	if len(resp.ItemIDs) == 0 {
		t.Log("No item id returned by create")
		t.FailNow()
	}

	items, err := client.ItemGet(ctx, zabbix.ItemGetParams{
		GetParameters: zabbix.GetParameters{
			Output: "extend",
		},
		ItemIDs: resp.ItemIDs,
	})
	if err != nil {
		t.Log("Item get failed:", err)
		t.FailNow()
	}

	if len(items) != 1 || items[0].MasterItemID != masterID {
		t.Logf("Dependent item mismatch: got=%+v", items)
		t.FailNow()
	}
}

func TestUpdateItem(t *testing.T) {

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

	itemID := createTestItem(t, ctx, client, "Test Item Update")

	defer func() {
		_, cleanupErr := client.ItemDelete(ctx, zabbix.ItemDeleteParams{itemID})
		if cleanupErr != nil {
			t.Log("Item cleanup failed:", cleanupErr)
		}
	}()

	updatedDescription := "after update"
	_, err = client.ItemUpdate(ctx, zabbix.ItemUpdateParams{
		zabbix.Item{
			ItemID:      itemID,
			Description: updatedDescription,
			Status:      zabbix.Ptr(zabbix.ItemStatusDisabled),
		},
	})
	if err != nil {
		t.Log("Item update failed:", err)
		t.FailNow()
	}

	items, err := client.ItemGet(ctx, zabbix.ItemGetParams{
		GetParameters: zabbix.GetParameters{
			Output: "extend",
		},
		ItemIDs: []string{itemID},
	})
	if err != nil {
		t.Log("Item get failed:", err)
		t.FailNow()
	}

	if len(items) == 0 {
		t.Log("Item not found after update")
		t.FailNow()
	}

	if items[0].Description != updatedDescription {
		t.Logf("Item description mismatch: got=%q want=%q", items[0].Description, updatedDescription)
		t.FailNow()
	}

	if items[0].Status == nil || *items[0].Status != zabbix.ItemStatusDisabled {
		t.Logf("Item status mismatch: got=%v want=%d", items[0].Status, zabbix.ItemStatusDisabled)
		t.FailNow()
	}
}

func TestDeleteItem(t *testing.T) {

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

	itemID := createTestItem(t, ctx, client, "Test Item Delete")

	deleteResp, err := client.ItemDelete(ctx, zabbix.ItemDeleteParams{itemID})
	if err != nil {
		t.Log("Item delete failed:", err)
		t.FailNow()
	}

	if len(deleteResp.ItemIDs) == 0 || deleteResp.ItemIDs[0] != itemID {
		t.Logf("Item ID mismatch between create and delete: create=%q delete=%v", itemID, deleteResp.ItemIDs)
		t.FailNow()
	}

	items, err := client.ItemGet(ctx, zabbix.ItemGetParams{
		ItemIDs: []string{itemID},
	})
	if err != nil {
		t.Log("Item get failed:", err)
		t.FailNow()
	}

	if len(items) != 0 {
		t.Logf("Expected item to be deleted, but found %d result(s)", len(items))
		t.FailNow()
	}
}
