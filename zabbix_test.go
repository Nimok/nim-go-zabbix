package zabbix_test

import (
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

	code := m.Run()

	os.Exit(code)
}

func TestClientWithoutAnyAuthMethod(t *testing.T) {
	_, err := zabbix.NewZabbixClient("any url")
	if err == nil {
		t.Log("client should not be allowed to be created without auth")
		t.FailNow()
	}

}

func TestClientWithApiToken(t *testing.T) {
	_, err := zabbix.NewZabbixClient(url, zabbix.WithAPIToken("some-token"))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

}

func TestClientWithUserPass(t *testing.T) {
	client, err := zabbix.NewZabbixClient(url, zabbix.WithUserPass(user, passwd),
		zabbix.WithBearerTokenTTL(10*time.Second))
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
	time.Sleep(2 * time.Second)
	client.StopTokenRefresher()
}
