package zabbix_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	zabbix "github.com/nimok/nim-go-zabbix"
)

var url string
var user string
var passwd string
var token string

func TestMain(m *testing.M) {

	_ = godotenv.Load()

	url = os.Getenv("TESTING_ZABBIX_URL")
	user = os.Getenv("TESTING_ZABBIX_USER")
	passwd = os.Getenv("TESTING_ZABBIX_PASS")

	if url == "" {
		fmt.Println("url not set")
		os.Exit(1)
	}

	if user == "" {
		fmt.Println("user not set")
		os.Exit(1)
	}
	if passwd == "" {
		fmt.Println("passwd not set")
		os.Exit(1)
	}

	bootstrapTokenId := setup()

	code := m.Run()

	teardown(bootstrapTokenId)

	os.Exit(code)
}

func setup() (tokenId string) {
	fmt.Println("setting up tests")
	ctx := context.Background()

	client, err := zabbix.NewClient(url, zabbix.WithUserPass(user, passwd))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Authenticate
	if err := client.Authenticate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	params := zabbix.Token{
		Name:   "bootstrap-token",
		UserID: "1",
	}

	tokenResp, err := client.TokenCreate(ctx, params)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	genResp, err := client.TokenGenerate(ctx, tokenResp.TokenIDs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	token = genResp[0].Token

	return genResp[0].TokenId
}

func teardown(bootstrapTokenId string) {
	fmt.Println("tearing down tests")
	ctx := context.Background()

	client, err := zabbix.NewClient(url, zabbix.WithUserPass(user, passwd))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Authenticate
	if err := client.Authenticate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = client.TokenDelete(ctx, zabbix.TokenDeleteParameters{
		bootstrapTokenId,
	})
	if err != nil {
		fmt.Println("cleanup failed, you might need to do some manual cleanup")
	}
}

func TestClientWithoutAnyAuthMethod(t *testing.T) {
	_, err := zabbix.NewClient("any url")
	if err == nil {
		t.Log("client should not be allowed to be created without auth")
		t.FailNow()
	}

}

func TestClientWithUserPass(t *testing.T) {
	client, err := zabbix.NewClient(url, zabbix.WithUserPass(user, passwd),
		zabbix.WithErrorCallback(func(err error) {
			t.Error(err)
		}))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	client.Authenticate()
	err = client.StartTokenRefresher(3 * time.Second)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	time.Sleep(15 * time.Second)

	filter := make(map[string]any, 0)
	filter["host"] = "Zabbix server"

	hosts, err := client.HostGet(context.TODO(), zabbix.HostGetParameters{
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
	client.StopTokenRefresher()
}

func TestClientWithAPIToken(t *testing.T) {
	client, err := zabbix.NewClient(url, zabbix.WithAPIToken(token),
		zabbix.WithErrorCallback(func(err error) {
			t.Log(err)
		}))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = client.Authenticate()
	if err != nil {
		t.Fatal(err)
	}

	filter := make(map[string]any, 0)
	filter["host"] = "Zabbix server"

	hosts, err := client.HostGet(context.TODO(), zabbix.HostGetParameters{
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

}
