package zabbix_test

import (
	"context"
	"testing"

	zabbix "github.com/nimok/nim-go-zabbix"
)

func TestTemplateGet(t *testing.T) {
	ctx := context.Background()

	client, err := zabbix.NewZabbixClient(url, zabbix.WithUserPass(user, passwd))
	if err != nil {
		t.Fatal(err)
	}

	// Authenticate
	if err := client.Authenticate(); err != nil {
		t.Fatal("Initial auth failed:", err)
	}

	filter := make(map[string]any, 0)
	filter["host"] = "Zabbix server health"

	template, err := client.TemplateGet(ctx, zabbix.TemplateGetParameters{
		GetParameters: zabbix.GetParameters{
			Filter: filter,
			Output: "extend",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(template) == 0 {
		t.Fatal("No template found")
	}

	if template[0].Host != "Zabbix server health" {
		t.Fatal("Template host does not match")
	}

}
