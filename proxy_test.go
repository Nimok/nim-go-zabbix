package zabbix_test

import (
	"context"
	"testing"
	"time"

	zabbix "github.com/nimok/nim-go-zabbix"
)

func TestProxyGet(t *testing.T) {
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
	filter["name"] = "zabbix-proxy-sqlite3"

	proxies, err := client.ProxyGet(ctx, zabbix.ProxyGetParameters{
		GetParameters: zabbix.GetParameters{
			Filter: filter,
			Output: "extend",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(proxies) == 0 {
		t.Fatal("No proxies found")
	}

	if proxies[0].Name != "zabbix-proxy-sqlite3" {
		t.Fatal("Proxy name does not match")
	}

}
