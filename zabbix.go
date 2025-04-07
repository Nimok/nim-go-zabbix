package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type ZabbixAPI interface {
	GetApiVersion() (*apiVersionRespone, error)
}

type ZabbixClient struct {
	URL      string
	Username string
	Password string

	token     string
	tokenTTL  time.Duration
	tokenLock sync.RWMutex

	stopChan chan struct{}
}

func NewZabbixClient(url, username, password string, tokenTTL time.Duration) ZabbixAPI {
	return &ZabbixClient{
		URL:      url,
		Username: username,
		Password: password,
		tokenTTL: tokenTTL,
		stopChan: make(chan struct{}),
	}
}

func (api *ZabbixClient) StartTokenRefresher() {
	go func() {
		ticker := time.NewTicker(api.tokenTTL - 5*time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Println("[INFO] Refreshing token...")
				if err := api.Authenticate(); err != nil {
					fmt.Println("[ERROR] Token refresh failed:", err)
				} else {
					fmt.Println("[INFO] Token refreshed successfully.")
				}

			case <-api.stopChan:
				fmt.Println("[INFO] Token refresher stopped.")
				return
			}
		}
	}()
}

func (api *ZabbixClient) StopTokenRefresher() {
	close(api.stopChan)
}

func (api *ZabbixClient) makeRequest(method string, params any) ([]byte, error) {
	api.tokenLock.RLock()
	token := api.token
	api.tokenLock.RUnlock()

	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      1,
	}

	reqBody, _ := json.Marshal(request)
	req, _ := http.NewRequest(http.MethodPost, api.URL, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json-rpc")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
