package zabbix_test

import (
	"context"
	"testing"

	zabbix "github.com/nimok/nim-go-zabbix"
)

func TestGetProblems(t *testing.T) {

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
	recent := false
	_, err = client.ProblemGet(ctx, zabbix.ProblemGetParams{
		GetParameters: zabbix.GetParameters{
			Output:    "extend",
			Limit:     10,
			Sortfield: []string{"eventid"},
			Sortorder: "DESC",
		},
		Recent: &recent,
		HostIDs: []string{
			"10084",
		},
	})

	if err != nil {
		t.Log("Problem get failed:", err)
		t.FailNow()
	}
}
