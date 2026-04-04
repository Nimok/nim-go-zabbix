package zabbix

import (
	"sync"
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
