package zabbix

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
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

	// Run tests
	code := m.Run()

	// Teardown if necessary
	os.Exit(code)
}

func TestApiVersion(t *testing.T) {
	client := NewZabbixClient(url, user, passwd, time.Hour)

	//if err := api.Authenticate(); err != nil {
	//	fmt.Println("Initial auth failed:", err)
	//	return
	//}

	//api.StartTokenRefresher()
	//defer api.StopTokenRefresher()

	resp, err := client.GetApiVersion()
	if err != nil {
		t.Log("smoketest failed")
		t.FailNow()
	}

	if resp.Result != "7.2.3" {
		t.Log("invalid api version")
		t.Fail()
	}

}
