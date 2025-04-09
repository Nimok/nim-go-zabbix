package zabbix_test

import (
	"context"
	"testing"
	"time"

	zabbix "github.com/nimok/nim-go-zabbix"
)

func TestHostCreateAndDelete(t *testing.T) {
	ctx := context.Background()

	client, err := zabbix.NewZabbixClient(url, zabbix.WithUserPass(user, passwd), zabbix.WithBearerTokenTTL(1*time.Hour))
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

	client, err := zabbix.NewZabbixClient(url, zabbix.WithUserPass(user, passwd), zabbix.WithBearerTokenTTL(1*time.Hour))
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
	}

	_, err = client.HostUpdate(ctx, updatedHost)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	_, err = client.HostDelete(ctx, hostResp.HostIDs)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

}

func TestHostCreateFailMissingPort(t *testing.T) {
	ctx := context.Background()

	client, err := zabbix.NewZabbixClient(url, zabbix.WithUserPass(user, passwd), zabbix.WithBearerTokenTTL(1*time.Hour))
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
				Type:  1,
				Main:  1,
				UseIP: 1,
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

	client, err := zabbix.NewZabbixClient(url, zabbix.WithUserPass(user, passwd), zabbix.WithBearerTokenTTL(1*time.Hour))
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
		t.Fail()
	}

	if hosts[0].Host != "Zabbix server" {
		t.Log("Host name does not match")
		t.Fail()
	}

}
