package zabbix_test

import (
	"context"
	"testing"
	"time"

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
		SelectHosts:  "extend",
		SelectGroups: "extend",
		SelectTags:   "extend",
	})

	if err != nil {
		t.Log("Maintenance get failed:", err)
		t.FailNow()
	}
}

func TestUpdateMaintenance(t *testing.T) {

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

	now := time.Now().Unix()
	createResp, err := client.MaintenanceCreate(ctx, zabbix.MaintenanceCreateParams{
		zabbix.Maintenance{
			Name:            "Test Maintenance Update " + time.Now().Format("20060102150405"),
			ActiveSince:     now + 60,
			ActiveTill:      now + 3600,
			Description:     "before update",
			MaintenanceType: zabbix.MaintenanceTypeWithData,
			Hosts: []zabbix.HostRef{
				{
					HostID: "10084",
				},
			},
			TimePeriods: []zabbix.MaintenanceTimePeriod{
				{
					TimePeriodType: zabbix.MaintenanceTimePeriodOneTime,
					StartDate:      now + 60,
					Period:         zabbix.MaintenanceDefaultPeriodSeconds,
				},
			},
		},
	})
	if err != nil {
		t.Log("Maintenance create failed:", err)
		t.FailNow()
	}

	if len(createResp.MaintenanceIDs) == 0 {
		t.Log("No maintenance id returned by create")
		t.FailNow()
	}

	defer func() {
		_, cleanupErr := client.MaintenanceDelete(ctx, zabbix.MaintenanceDeleteParams{createResp.MaintenanceIDs[0]})
		if cleanupErr != nil {
			t.Log("Maintenance cleanup failed:", cleanupErr)
		}
	}()

	updatedDescription := "after update"
	_, err = client.MaintenanceUpdate(ctx, zabbix.MaintenanceUpdateParams{
		zabbix.Maintenance{
			MaintenanceID: createResp.MaintenanceIDs[0],
			Description:   updatedDescription,
		},
	})
	if err != nil {
		t.Log("Maintenance update failed:", err)
		t.FailNow()
	}

	maintenances, err := client.MaintenanceGet(ctx, zabbix.MaintenanceGetParams{
		MaintenanceIDs: []string{createResp.MaintenanceIDs[0]},
	})
	if err != nil {
		t.Log("Maintenance get failed:", err)
		t.FailNow()
	}

	if len(*maintenances) == 0 {
		t.Log("Maintenance not found after update")
		t.FailNow()
	}

	if (*maintenances)[0].Description != updatedDescription {
		t.Logf("Maintenance description mismatch: got=%q want=%q", (*maintenances)[0].Description, updatedDescription)
		t.FailNow()
	}
}

func TestDeleteMaintenance(t *testing.T) {

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

	now := time.Now().Unix()
	createResp, err := client.MaintenanceCreate(ctx, zabbix.MaintenanceCreateParams{
		zabbix.Maintenance{
			Name:            "Test Maintenance Delete " + time.Now().Format("20060102150405"),
			ActiveSince:     now + 60,
			ActiveTill:      now + 3600,
			Description:     "to be deleted",
			MaintenanceType: zabbix.MaintenanceTypeWithData,
			Hosts: []zabbix.HostRef{
				{
					HostID: "10084",
				},
			},
			TimePeriods: []zabbix.MaintenanceTimePeriod{
				{
					TimePeriodType: zabbix.MaintenanceTimePeriodOneTime,
					StartDate:      now + 60,
					Period:         zabbix.MaintenanceDefaultPeriodSeconds,
				},
			},
		},
	})
	if err != nil {
		t.Log("Maintenance create failed:", err)
		t.FailNow()
	}

	if len(createResp.MaintenanceIDs) == 0 {
		t.Log("No maintenance id returned by create")
		t.FailNow()
	}

	deleteResp, err := client.MaintenanceDelete(ctx, zabbix.MaintenanceDeleteParams{createResp.MaintenanceIDs[0]})
	if err != nil {
		t.Log("Maintenance delete failed:", err)
		t.FailNow()
	}

	if len(deleteResp.MaintenanceIDs) == 0 {
		t.Log("No maintenance id returned by delete")
		t.FailNow()
	}

	if deleteResp.MaintenanceIDs[0] != createResp.MaintenanceIDs[0] {
		t.Logf("Maintenance ID mismatch between create and delete: create=%q delete=%q", createResp.MaintenanceIDs[0], deleteResp.MaintenanceIDs[0])
		t.FailNow()
	}

	maintenances, err := client.MaintenanceGet(ctx, zabbix.MaintenanceGetParams{
		MaintenanceIDs: []string{createResp.MaintenanceIDs[0]},
	})
	if err != nil {
		t.Log("Maintenance get failed:", err)
		t.FailNow()
	}

	if len(*maintenances) != 0 {
		t.Logf("Expected maintenance to be deleted, but found %d result(s)", len(*maintenances))
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
	now := time.Now().Unix()
	resp, err := client.MaintenanceCreate(ctx, zabbix.MaintenanceCreateParams{
		zabbix.Maintenance{
			Name:            "Test Maintenance Create " + time.Now().Format("20060102150405"),
			ActiveSince:     now + 60,
			ActiveTill:      now + 3600,
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
					StartDate:      now + 60,
					Period:         zabbix.MaintenanceDefaultPeriodSeconds,
				},
			},
		},
	})

	if err != nil {
		t.Log("Maintenance create failed:", err)
		t.FailNow()
	}
	if len(resp.MaintenanceIDs) == 0 {
		t.Log("No maintenance id returned by create")
		t.FailNow()
	}

	defer func() {
		_, cleanupErr := client.MaintenanceDelete(ctx, zabbix.MaintenanceDeleteParams{resp.MaintenanceIDs[0]})
		if cleanupErr != nil {
			t.Log("Maintenance cleanup failed:", cleanupErr)
		}
	}()

	_, err = client.MaintenanceGet(ctx, zabbix.MaintenanceGetParams{
		MaintenanceIDs: resp.MaintenanceIDs},
	)
	if err != nil {
		t.Log("Maintenance get failed:", err)
		t.FailNow()
	}

}
