package zabbix_test

import (
	"context"
	"testing"

	zabbix "github.com/nimok/nim-go-zabbix"
)

func TestHostInterfaceCreateAndDelete(t *testing.T) {
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

	// Create a host
	host := zabbix.Host{
		Host:        "test-host",
		Description: "Test host",
		Interfaces:  []zabbix.HostInterface{},
		Groups: []zabbix.HostGroup{
			{
				GroupID: "2",
			},
		},
		Macros: []zabbix.Macro{
			{
				Macro:       "{$MACRO_NAME}",
				Value:       "macro value",
				Description: "Test macro",
			},
		},
	}

	hostResp, err := client.HostCreate(ctx, []zabbix.Host{host})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	interfaceCreateResp, err := client.HostInterfaceCreate(ctx, zabbix.HostInterface{
		HostID: hostResp.HostIDs[0],
		Type:   zabbix.InterfaceTypeSNMP,
		Main:   zabbix.MainInterfaceYes,
		UseIP:  zabbix.UseIPOptionIP,
		IP:     "127.0.0.1",
		Port:   "7777",
		Details: zabbix.InterfaceDetails{
			Version:   zabbix.SNMPv2c,
			Bulk:      zabbix.BulkEnabled,
			Community: "{$SNMP_COMMUNITY}",
		},
	})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	_, err = client.HostInterfaceDelete(ctx, []string{interfaceCreateResp.HostInterfaceIDs[0]})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	interfaceGetResp, err := client.HostInterfaceGet(ctx, zabbix.HostInterfaceGetParams{
		HostIDs: []string{hostResp.HostIDs[0]},
	})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	_, err = client.HostDelete(ctx, hostResp.HostIDs)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if len(interfaceGetResp) != 0 {
		t.Logf("Interfaces should be empty")
		t.FailNow()
	}

}
