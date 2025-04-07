package zabbix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"errors"

	"github.com/mitchellh/mapstructure"
)

type ZabbixAPI interface {
	Authenticate() error

	HostCreate(ctx context.Context, params Host) (*hostCreateResponse, error)
	HostDelete(ctx context.Context, params []string) (*hostDeleteResponse, error)
	HostUpdate(ctx context.Context, params Host) (*hostUpdateResponse, error)
}

type zabbixClient struct {
	url      string
	username string
	password string
	apiToken string

	bearerToken        string
	bearerTokenRefresh bool
	bearerTokenTTL     time.Duration
	bearerTokenLock    sync.RWMutex
	stopChan           chan struct{}
}

type ZabbixClientOption func(*zabbixClient)

func WithUserPass(username, password string) ZabbixClientOption {
	return func(c *zabbixClient) {
		c.username = username
		c.password = password
	}
}

func WithBearerTokenTTL(ttl time.Duration) ZabbixClientOption {
	return func(c *zabbixClient) {
		c.bearerTokenTTL = ttl
	}
}

func WithAPIToken(apiToken string) ZabbixClientOption {
	return func(c *zabbixClient) {
		c.apiToken = apiToken
	}
}

func WithBearerTokenRefresh() ZabbixClientOption {
	return func(c *zabbixClient) {
		c.bearerTokenRefresh = true
	}
}

func NewZabbixClient(url string, opts ...ZabbixClientOption) (ZabbixAPI, error) {
	client := &zabbixClient{
		url:      url,
		stopChan: make(chan struct{}),
	}
	for _, opt := range opts {
		opt(client)
	}

	err := validateClient(client)
	if err != nil {
		return nil, err
	}

	if client.bearerTokenRefresh {
		client.StartTokenRefresher()
		defer client.StopTokenRefresher()
	}

	return client, nil
}

func validateClient(c *zabbixClient) error {
	if c.url == "" {
		return errors.New("url cant be empty")
	}

	if c.apiToken == "" {
		if c.username == "" || c.password == "" {
			return errors.New("you need to supply an api token or a user/password login")
		}
	}

	if c.apiToken != "" {
		if c.username != "" || c.password != "" {
			return errors.New("you cant supply both an api token and a user/password login")
		}
	}

	return nil
}

func (c *zabbixClient) StartTokenRefresher() {
	go func() {
		ticker := time.NewTicker(c.bearerTokenTTL - 5*time.Minute)
		defer ticker.Stop()
		fmt.Println("[INFO] Starting token refresher...")
		for {
			select {
			case <-ticker.C:
				fmt.Println("[INFO] Refreshing token...")
				if err := c.Authenticate(); err != nil {
					fmt.Println("[ERROR] Token refresh failed:", err)
				} else {
					fmt.Println("[INFO] Token refreshed successfully.")
				}

			case <-c.stopChan:
				fmt.Println("[INFO] Token refresher stopped.")
				return
			}
		}
	}()
}

func (client *zabbixClient) StopTokenRefresher() {
	close(client.stopChan)
}

type apiResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  any    `json:"result"`
	ID      int    `json:"id"`
	Error   *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	} `json:"error,omitempty"`
}

func (c *zabbixClient) makeRequest(ctx context.Context, method string, params any, result any) error {
	c.bearerTokenLock.RLock()
	token := c.bearerToken
	c.bearerTokenLock.RUnlock()

	res := &apiResponse{
		Result: result,
	}

	request := map[string]any{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      1,
	}

	reqBody, _ := json.Marshal(request)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json-rpc")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	rawConf := make(map[string]any)

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&rawConf); err != nil {
		return fmt.Errorf("json decode error: %v", err)
	}

	mapper, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           res,
		TagName:          "json",
	})
	if err != nil {
		return fmt.Errorf("mapstructure create decoder error: %v", err)
	}

	if err := mapper.Decode(rawConf); err != nil {
		return fmt.Errorf("mapstructure decode error: %v", err)
	}

	if res.Error != nil {
		return fmt.Errorf("API error: %d - %s %s", res.Error.Code, res.Error.Message, res.Error.Data)
	}

	return nil

}
