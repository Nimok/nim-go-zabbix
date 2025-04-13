package zabbix_test

import (
	"context"
	"testing"
	"time"

	zabbix "github.com/nimok/nim-go-zabbix"
)

func TestHostgroupGet(t *testing.T) {
	ctx := context.Background()

	client, err := zabbix.NewZabbixClient(url, zabbix.WithUserPass(user, passwd), zabbix.WithBearerTokenTTL(1*time.Hour))
	if err != nil {
		t.Fatal(err)
	}

	// Authenticate
	if err := client.Authenticate(); err != nil {
		t.Fatal("Initial auth failed:", err)
	}

	filter := make(map[string]any, 0)
	filter["name"] = "Zabbix servers"

	hostgroups, err := client.HostgroupGet(ctx, zabbix.HostGroupGetParameters{
		GetParameters: zabbix.GetParameters{
			Filter: filter,
			Output: "extend",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(hostgroups) == 0 {
		t.Fatal("No hostgroups found")
	}

	if hostgroups[0].Name != "Zabbix servers" {
		t.Fatal("Hostgroup name does not match")
	}

}
