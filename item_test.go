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

func TestGetItemSelects(t *testing.T) {

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

	// Zabbix returns [] instead of an object for an empty valuemap/discoveryRule/itemDiscovery,
	// so this also checks that those decode to nil.
	items, err := client.ItemGet(ctx, zabbix.ItemGetParams{
		GetParameters: zabbix.GetParameters{
			Output: "extend",
		},
		HostIDs:             []string{"10084"},
		SelectInterfaces:    "extend",
		SelectValueMap:      "extend",
		SelectDiscoveryRule: "extend",
		SelectItemDiscovery: "extend",
	})
	if err != nil {
		t.Log("Item get failed:", err)
		t.FailNow()
	}

	var withInterface, withValueMap bool
	for _, item := range items {
		if len(item.Interfaces) > 0 && item.Interfaces[0].InterfaceID == item.InterfaceID {
			withInterface = true
		}

		hasValueMap := item.ValueMapID != "" && item.ValueMapID != "0"
		if hasValueMap != (item.ValueMap != nil) {
			t.Logf("Item %s valuemap mismatch: valuemapid=%q valuemap=%+v", item.ItemID, item.ValueMapID, item.ValueMap)
			t.FailNow()
		}
		if hasValueMap {
			if item.ValueMap.ValueMapID != item.ValueMapID || len(item.ValueMap.Mappings) == 0 {
				t.Logf("Item %s valuemap not decoded: got=%+v", item.ItemID, item.ValueMap)
				t.FailNow()
			}
			withValueMap = true
		}

		if item.Flags != 4 && // 4 = discovered item
			(item.DiscoveryRule != nil || item.ItemDiscovery != nil) {
			t.Logf("Item %s is not discovered but has discovery data: rule=%+v discovery=%+v", item.ItemID, item.DiscoveryRule, item.ItemDiscovery)
			t.FailNow()
		}
	}

	if !withInterface {
		t.Log("No item with decoded interfaces found on Zabbix server host")
		t.FailNow()
	}
	if !withValueMap {
		t.Log("No item with decoded value map found on Zabbix server host")
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

func TestCreateSSHItemPublicKey(t *testing.T) {

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
	resp, err := client.ItemCreate(ctx, zabbix.ItemCreateParams{
		zabbix.Item{
			HostID:      "10084",
			InterfaceID: "1",
			Name:        "Test Item SSH " + suffix,
			Key:         "ssh.run[test," + suffix + "]",
			Type:        zabbix.Ptr(zabbix.ItemTypeSSHAgent),
			ValueType:   zabbix.Ptr(zabbix.ItemValueTypeText),
			Delay:       "1h",
			Params:      "uptime",
			AuthType:    zabbix.ItemSSHAuthPublicKey,
			Username:    "zabbix",
			PublicKey:   "id_ed25519.pub",
			PrivateKey:  "id_ed25519",
			Status:      zabbix.Ptr(zabbix.ItemStatusDisabled),
		},
	})
	if err != nil {
		t.Log("SSH item create failed:", err)
		t.FailNow()
	}

	defer func() {
		_, cleanupErr := client.ItemDelete(ctx, zabbix.ItemDeleteParams(resp.ItemIDs))
		if cleanupErr != nil {
			t.Log("Item cleanup failed:", cleanupErr)
		}
	}()

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

	if len(items) != 1 || items[0].PublicKey != "id_ed25519.pub" || items[0].PrivateKey != "id_ed25519" {
		t.Logf("SSH item key files mismatch: got=%+v", items)
		t.FailNow()
	}
}

func TestCreateIPMIItem(t *testing.T) {

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

	// IPMI items need an IPMI interface, so use a dedicated host.
	suffix := time.Now().Format("20060102150405.000000")
	hostResp, err := client.HostCreate(ctx, []zabbix.Host{
		{
			Host:   "test-ipmi-item-" + suffix,
			Groups: []zabbix.HostGroup{{GroupID: "2"}},
			Interfaces: []zabbix.HostInterface{
				{
					Type:  zabbix.InterfaceTypeIPMI,
					Main:  zabbix.MainInterfaceYes,
					UseIP: zabbix.UseIPOptionIP,
					IP:    "127.0.0.1",
					Port:  "623",
				},
			},
		},
	})
	if err != nil {
		t.Log("Host create failed:", err)
		t.FailNow()
	}

	// Deleting the host also deletes its items.
	defer func() {
		_, cleanupErr := client.HostDelete(ctx, hostResp.HostIDs)
		if cleanupErr != nil {
			t.Log("Host cleanup failed:", cleanupErr)
		}
	}()

	interfaces, err := client.HostInterfaceGet(ctx, zabbix.HostInterfaceGetParams{
		HostIDs: hostResp.HostIDs,
	})
	if err != nil || len(interfaces) != 1 {
		t.Logf("Host interface get failed: err=%v interfaces=%+v", err, interfaces)
		t.FailNow()
	}

	resp, err := client.ItemCreate(ctx, zabbix.ItemCreateParams{
		zabbix.Item{
			HostID:      hostResp.HostIDs[0],
			InterfaceID: interfaces[0].InterfaceID,
			Name:        "Test Item IPMI",
			Key:         "ipmi.test",
			Type:        zabbix.Ptr(zabbix.ItemTypeIPMIAgent),
			ValueType:   zabbix.Ptr(zabbix.ItemValueTypeFloat),
			Delay:       "1m",
			IPMISensor:  "CPU Temp",
		},
	})
	if err != nil {
		t.Log("IPMI item create failed:", err)
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

	if len(items) != 1 || items[0].IPMISensor != "CPU Temp" {
		t.Logf("IPMI item sensor mismatch: got=%+v", items)
		t.FailNow()
	}
}
