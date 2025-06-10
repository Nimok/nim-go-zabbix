package zabbix_test

import (
	"context"
	"testing"

	zabbix "github.com/nimok/nim-go-zabbix"
)

func TestTokenCreateAndGenerate(t *testing.T) {
	ctx := context.Background()

	client, err := zabbix.NewClient(url, zabbix.WithUserPass(user, passwd))
	if err != nil {
		t.Fatal(err)
	}

	// Authenticate
	if err := client.Authenticate(); err != nil {
		t.Fatal("Initial auth failed:", err)
	}

	params := zabbix.Token{
		Name:   "testing-token",
		UserID: "1",
	}

	tokenResp, err := client.TokenCreate(ctx, params)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.TokenGenerate(ctx, tokenResp.TokenIDs)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.TokenDelete(ctx, tokenResp.TokenIDs)
	if err != nil {
		t.Fatal(err)
	}

}
