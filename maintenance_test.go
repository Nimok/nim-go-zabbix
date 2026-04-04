package zabbix_test

import (
	"context"
	"testing"

	zabbix "github.com/nimok/nim-go-zabbix"
)

func TestGetMaintenance(t *testing.T) {

	ctx := context.Background()

	client, err := zabbix.NewClient(url, zabbix.WithUserPass(user, passwd))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// Authenticate
	if err := client.Authenticate(); err != nil {
		t.Log("Initial auth failed:", err)
		t.FailNow()
	}
	_, err = client.MaintenanceGet(ctx, zabbix.MaintenanceGetParams{
		GetParameters: zabbix.GetParameters{
			Output: "extend",
			Limit:  10,
		},
		SelectHosts:      "extend",
		SelectHostGroups: "extend",
		SelectTags:       "extend",
	})

	if err != nil {
		t.Log("Maintenance get failed:", err)
		t.FailNow()
	}
}

func TestCreateMaintenance(t *testing.T) {

	ctx := context.Background()

	client, err := zabbix.NewClient(url, zabbix.WithUserPass(user, passwd))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// Authenticate
	if err := client.Authenticate(); err != nil {
		t.Log("Initial auth failed:", err)
		t.FailNow()
	}
	resp, err := client.MaintenanceCreate(ctx, zabbix.MaintenanceCreateParams{
		zabbix.Maintenance{
			Name:            "Test Maintenance",
			ActiveSince:     1769457981,
			ActiveTill:      1769458981,
			Description:     "some descr",
			MaintenanceType: zabbix.MaintenanceTypeWithData,
			Hosts: []zabbix.HostRef{
				{
					HostID: "10084",
				},
			},
			TimePeriods: []zabbix.MaintenanceTimePeriod{
				{
					TimePeriodType: zabbix.MaintenanceTimePeriodOneTime,
				},
			},
		},
	})

	if err != nil {
		t.Log("Maintenance create failed:", err)
		t.FailNow()
	}

	_, err = client.MaintenanceGet(ctx, zabbix.MaintenanceGetParams{
		MaintenanceIDs: resp.MaintenanceIDs},
	)
	if err != nil {
		t.Log("Maintenance get failed:", err)
		t.FailNow()
	}

}
