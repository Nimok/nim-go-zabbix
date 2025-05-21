package zabbix_test

import (
	"context"
	"testing"
	"time"

	zabbix "github.com/nimok/nim-go-zabbix"
)

func TestProxyGet(t *testing.T) {
	ctx := context.Background()
	proxyName := "some-proxy"

	client, err := zabbix.NewZabbixClient(url, zabbix.WithUserPass(user, passwd), zabbix.WithBearerTokenTTL(1*time.Hour))
	if err != nil {
		t.Fatal(err)
	}

	// Authenticate
	if err := client.Authenticate(); err != nil {
		t.Fatal("Initial auth failed:", err)
	}

	params := zabbix.ProxyCreateParameters{
		Proxy: zabbix.Proxy{
			Name:          proxyName,
			OperatingMode: 0,
		},
	}

	createResp, err := client.ProxyCreate(ctx, params)
	if err != nil {
		t.Fatal(err)
	}

	filter := make(map[string]any, 0)
	filter["name"] = proxyName

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

	if proxies[0].Name != proxyName {
		t.Fatal("Proxy name does not match")
	}

	_, err = client.ProxyDelete(ctx, createResp.ProxyIDs)
	if err != nil {
		t.Fatal(err)
	}

}

func TestProxyCreateAndDelete(t *testing.T) {
	ctx := context.Background()

	client, err := zabbix.NewZabbixClient(url, zabbix.WithUserPass(user, passwd), zabbix.WithBearerTokenTTL(1*time.Hour))
	if err != nil {
		t.Fatal(err)
	}

	// Authenticate
	if err := client.Authenticate(); err != nil {
		t.Fatal("Initial auth failed:", err)
	}

	params := zabbix.ProxyCreateParameters{
		Proxy: zabbix.Proxy{
			Name:          "my-proxy",
			OperatingMode: 0,
		},
	}

	createResp, err := client.ProxyCreate(ctx, params)
	if err != nil {
		t.Fatal(err)
	}

	deleteResp, err := client.ProxyDelete(ctx, createResp.ProxyIDs)
	if err != nil {
		t.Fatal(err)
	}

	if deleteResp.ProxyIDs[0] != createResp.ProxyIDs[0] {
		t.Fatal("proxy id mismatch")
	}

}

func createProxy(ctx context.Context, client zabbix.ZabbixClient) (proxyId string, err error) {
	params := zabbix.ProxyCreateParameters{
		Proxy: zabbix.Proxy{
			Name:          "my-proxy",
			OperatingMode: 0,
		},
	}

	createResp, err := client.ProxyCreate(ctx, params)
	if err != nil {
		return "", err
	}

	return createResp.ProxyIDs[0], nil
}

func deleteProxy(ctx context.Context, client zabbix.ZabbixClient, proxyId string) error {

	_, err := client.ProxyDelete(ctx, []string{
		proxyId,
	})
	if err != nil {
		return err
	}

	return nil
}
