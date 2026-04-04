package zabbix

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func newTokenRefresherTestClient(t *testing.T) *zabbixClient {
	t.Helper()

	client, err := NewClient("http://localhost", WithAPIToken("test-token"))
	if err != nil {
		t.Fatalf("failed to build test client: %v", err)
	}

	zClient, ok := client.(*zabbixClient)
	if !ok {
		t.Fatal("failed to cast client to zabbixClient")
	}

	return zClient
}

func TestStartTokenRefresherInvalidInterval(t *testing.T) {
	client := newTokenRefresherTestClient(t)

	if err := client.StartTokenRefresher(0); err == nil {
		t.Fatal("expected error for zero interval")
	}

	if err := client.StartTokenRefresher(-1 * time.Second); err == nil {
		t.Fatal("expected error for negative interval")
	}
}

func TestStopTokenRefresherIdempotent(t *testing.T) {
	client := newTokenRefresherTestClient(t)

	if err := client.StartTokenRefresher(time.Hour); err != nil {
		t.Fatalf("failed to start refresher: %v", err)
	}

	client.StopTokenRefresher()
	client.StopTokenRefresher()
}

func TestStartTokenRefresherWhileRunningReturnsError(t *testing.T) {
	client := newTokenRefresherTestClient(t)
	defer client.StopTokenRefresher()

	if err := client.StartTokenRefresher(time.Hour); err != nil {
		t.Fatalf("failed to start refresher: %v", err)
	}

	if err := client.StartTokenRefresher(time.Hour); err == nil {
		t.Fatal("expected error when starting refresher twice")
	}
}

func TestStartStopStartTokenRefresher(t *testing.T) {
	client := newTokenRefresherTestClient(t)

	if err := client.StartTokenRefresher(time.Hour); err != nil {
		t.Fatalf("failed to start refresher: %v", err)
	}
	client.StopTokenRefresher()

	if err := client.StartTokenRefresher(time.Hour); err != nil {
		t.Fatalf("failed to restart refresher: %v", err)
	}
	client.StopTokenRefresher()
}

func TestStopTokenRefresherConcurrentCalls(t *testing.T) {
	client := newTokenRefresherTestClient(t)

	if err := client.StartTokenRefresher(time.Hour); err != nil {
		t.Fatalf("failed to start refresher: %v", err)
	}

	const callers = 10
	var wg sync.WaitGroup
	wg.Add(callers)

	for i := 0; i < callers; i++ {
		go func() {
			defer wg.Done()
			client.StopTokenRefresher()
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("concurrent StopTokenRefresher calls timed out")
	}
}

func TestStopTokenRefresherFromErrorCallbackDoesNotDeadlock(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json-rpc")
		_, _ = fmt.Fprint(w, `{"jsonrpc":"2.0","error":{"code":-32602,"message":"Invalid params.","data":"auth failed"},"id":1}`)
	}))
	defer server.Close()

	var callbackHits atomic.Int32
	callbackDone := make(chan struct{})
	var clientIntf Client
	var err error

	clientIntf, err = NewClient(
		server.URL,
		WithUserPass("user", "pass"),
		WithErrorCallback(func(err error) {
			callbackHits.Add(1)
			client, ok := clientIntf.(*zabbixClient)
			if ok {
				client.StopTokenRefresher()
			}
			select {
			case <-callbackDone:
			default:
				close(callbackDone)
			}
		}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	client := clientIntf.(*zabbixClient)
	if err := client.StartTokenRefresher(10 * time.Millisecond); err != nil {
		t.Fatalf("failed to start refresher: %v", err)
	}

	select {
	case <-callbackDone:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for error callback to stop refresher")
	}

	if callbackHits.Load() == 0 {
		t.Fatal("expected error callback to be called at least once")
	}
}
