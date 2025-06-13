package zabbix_test

import (
	"context"
	"testing"

	zabbix "github.com/nimok/nim-go-zabbix"
)

func TestHostCreateAndDelete(t *testing.T) {
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
		Interfaces: []zabbix.HostInterface{
			{
				Type:  zabbix.InterfaceTypeSNMP,
				Main:  zabbix.MainInterfaceYes,
				UseIP: zabbix.UseIPOptionIP,
				IP:    "127.0.0.1",
				Port:  "10050",
				Details: zabbix.InterfaceDetails{
					Version:   zabbix.SNMPv2c,
					Bulk:      zabbix.BulkEnabled,
					Community: "{$SNMP_COMMUNITY}",
				},
			},
		},
		Groups: []zabbix.HostGroup{
			{
				GroupID: "2",
			},
		},
		Templates: []zabbix.Template{
			{
				TemplateID: "10395",
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

	hostResp, err := client.HostCreate(ctx, host)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Logf("Host created with ID: %s", hostResp.HostIDs[0])

	delResp, err := client.HostDelete(ctx, hostResp.HostIDs)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if hostResp.HostIDs[0] != delResp.HostIDs[0] {
		t.Log("Host IDs do not match")
		t.Fail()
	}

}

func TestHostCreateAndUpdate(t *testing.T) {
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
		Status:      1,
		Interfaces: []zabbix.HostInterface{
			{
				Type:  1,
				Main:  1,
				UseIP: 1,
				IP:    "127.0.0.1",
				Port:  "10050",
			},
		},
		Groups: []zabbix.HostGroup{
			{
				GroupID: "2",
			},
		},
		Templates: []zabbix.Template{
			{
				TemplateID: "10001",
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

	hostResp, err := client.HostCreate(ctx, host)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Logf("Host created with ID: %s", hostResp.HostIDs[0])

	updatedHost := zabbix.Host{
		HostID:      hostResp.HostIDs[0],
		Description: "Updated test host",
		Status:      0,
	}

	_, err = client.HostUpdate(ctx, updatedHost)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	newHost, err := client.HostGet(ctx, zabbix.HostGetParameters{
		HostIDs: hostResp.HostIDs,
	})
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	_, err = client.HostDelete(ctx, hostResp.HostIDs)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if newHost[0].Status == host.Status {
		t.Log("Status did not update")
		t.Fail()
	}

}

func TestHostCreateFailMissingPort(t *testing.T) {
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
		Interfaces: []zabbix.HostInterface{
			{
				Type:  zabbix.InterfaceTypeSNMP,
				Main:  zabbix.MainInterfaceYes,
				UseIP: zabbix.UseIPOptionIP,
				IP:    "127.0.0.1",
				//Port:  "10050",
			},
		},
		Groups: []zabbix.HostGroup{
			{
				GroupID: "2",
			},
		},
		Templates: []zabbix.Template{
			{
				TemplateID: "10001",
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

	_, err = client.HostCreate(ctx, host)
	if err == nil {
		t.Log(err)
		t.FailNow()
	}

}

func TestHostGet(t *testing.T) {
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

	filter := make(map[string]any, 0)
	filter["host"] = "Zabbix server"

	hosts, err := client.HostGet(ctx, zabbix.HostGetParameters{
		GetParameters: zabbix.GetParameters{
			Filter: filter,
			Output: "extend",
		},
	})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(hosts) == 0 {
		t.Log("No hosts found")
		t.FailNow()
	}

	if hosts[0].Host != "Zabbix server" {
		t.Log("Host name does not match")
		t.Fail()
	}

}

func TestHostCreateMonitoredByProxy(t *testing.T) {

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

	hostname := "nimgo-host"
	hostgroupId := "2"
	templateId := "10395"
	community := "public"
	ipAddress := "127.0.0.1"

	proxyId, err := createProxy(ctx, client)
	if err != nil {
		t.Fatal(err)
	}

	hostToCreate := zabbix.Host{
		Host: hostname,
		Groups: []zabbix.HostGroup{
			{
				GroupID: hostgroupId,
			},
		},
		Templates: []zabbix.Template{
			{
				TemplateID: templateId,
			},
		},
		Interfaces: []zabbix.HostInterface{
			{
				IP:    ipAddress,
				Type:  zabbix.InterfaceTypeSNMP,
				UseIP: zabbix.UseIPOptionIP,
				DNS:   "",
				Port:  "161",
				Main:  zabbix.MainInterfaceYes,
				Details: zabbix.InterfaceDetails{
					Version:   zabbix.SNMPv2c,
					Bulk:      zabbix.BulkEnabled,
					Community: "{$SNMP_COMMUNITY}",
				},
			},
		},
		Macros: []zabbix.Macro{
			{
				Macro: "{$SNMP_COMMUNITY}",
				Value: community,
			},
		},
		MonitoredBy: zabbix.MonitoredByProxy,
		ProxyID:     proxyId,
	}

	hostResp, err := client.HostCreate(ctx, hostToCreate)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Logf("Host created with ID: %s", hostResp.HostIDs[0])

	delResp, err := client.HostDelete(ctx, hostResp.HostIDs)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if hostResp.HostIDs[0] != delResp.HostIDs[0] {
		t.Log("Host IDs do not match")
		t.Fail()
	}

	if err := deleteProxy(ctx, client, proxyId); err != nil {
		t.Log(err)
		t.Fail()
	}

	ok, err := client.Logout(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if ok != true {
		t.Fatal("logout failed")
	}
}
