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

const (
	SelectExtendedOutput = "extend"
	SelectCount          = "count"
)

const (
	GetParametersSortOrderASC  = "ASC"
	GetParametersSortOrderDESC = "DESC"
)

type Client interface {
	Authenticate() error
	StartTokenRefresher(refreshInterval time.Duration) error
	StopTokenRefresher()

	HostGet(ctx context.Context, params HostGetParameters) ([]Host, error)
	HostCreate(ctx context.Context, params []Host) (*HostCreateResponse, error)
	HostUpdate(ctx context.Context, params Host) (*HostUpdateResponse, error)
	HostDelete(ctx context.Context, params []string) (*HostDeleteResponse, error)
	HostMassAdd(ctx context.Context, params HostMassAddParams) (*HostMassAddResponse, error)

	HostInterfaceGet(ctx context.Context, params HostInterfaceGetParams) ([]HostInterface, error)
	HostInterfaceCreate(ctx context.Context, params HostInterface) (*HostInterfaceCreateResponse, error)
	HostInterfaceUpdate(ctx context.Context, params HostInterface) (*HostInterfaceUpdateResponse, error)
	HostInterfaceDelete(ctx context.Context, params []string) (*HostInterfaceDeleteResponse, error)

	HostgroupGet(ctx context.Context, params HostGroupGetParameters) ([]HostGroup, error)

	ProblemGet(ctx context.Context, params ProblemGetParams) (*[]Problem, error)

	ProxyGet(ctx context.Context, params ProxyGetParameters) ([]Proxy, error)
	ProxyCreate(ctx context.Context, params ProxyCreateParameters) (*ProxyCreateResponse, error)
	ProxyDelete(ctx context.Context, params []string) (*ProxyDeleteResponse, error)

	TemplateGet(ctx context.Context, params TemplateGetParameters) ([]Template, error)

	TokenCreate(ctx context.Context, params Token) (*TokenCreateResponse, error)
	TokenGenerate(ctx context.Context, params TokenGenerateParameters) ([]TokenGenerateResponse, error)
	TokenDelete(ctx context.Context, params TokenDeleteParameters) (*TokenDeleteResponse, error)

	MaintenanceGet(ctx context.Context, params MaintenanceGetParams) (*[]Maintenance, error)
	MaintenanceCreate(ctx context.Context, params MaintenanceCreateParams) (*MaintenanceCreateResponse, error)
	MaintenanceUpdate(ctx context.Context, params MaintenanceUpdateParams) (*MaintenanceUpdateResponse, error)
	MaintenanceDelete(ctx context.Context, params MaintenanceDeleteParams) (*MaintenanceDeleteResponse, error)

	Logout(ctx context.Context) (LogoutSuccess, error)
}

var _ Client = (*zabbixClient)(nil)

type zabbixClient struct {
	url      string
	username string
	password string
	apiToken string

	bearerToken     string
	bearerTokenLock sync.RWMutex

	refresherLock    sync.Mutex
	refresherRunning bool
	stopChan         chan struct{}
	refresherDone    chan struct{}
	stopOnce         *sync.Once

	errorCallback func(error)
}

type ClientOption func(*zabbixClient)

func WithUserPass(username, password string) ClientOption {
	return func(c *zabbixClient) {
		c.username = username
		c.password = password
	}
}

func WithAPIToken(apiToken string) ClientOption {
	return func(c *zabbixClient) {
		c.apiToken = apiToken
	}
}

func WithErrorCallback(callback func(error)) ClientOption {
	return func(c *zabbixClient) {
		c.errorCallback = callback
	}
}

func NewClient(url string, opts ...ClientOption) (Client, error) {
	client := &zabbixClient{
		url:           url,
		errorCallback: func(err error) {},
	}
	for _, opt := range opts {
		opt(client)
	}

	err := validateClient(client)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func validateClient(c *zabbixClient) error {
	if c.url == "" {
		return errors.New("url can't be empty")
	}

	if c.apiToken == "" {
		if c.username == "" || c.password == "" {
			return errors.New("you need to supply an api token or a user/password login")
		}
	}

	if c.apiToken != "" {
		if c.username != "" || c.password != "" {
			return errors.New("you can't supply both an api token and a user/password login")
		}
	}

	return nil
}

func (c *zabbixClient) StartTokenRefresher(refreshInterval time.Duration) error {
	if refreshInterval <= 0 {
		return errors.New("refresh interval must be greater than 0")
	}

	c.refresherLock.Lock()
	if c.refresherRunning {
		c.refresherLock.Unlock()
		return errors.New("token refresher is already running")
	}

	stopChan := make(chan struct{})
	doneChan := make(chan struct{})
	stopOnce := &sync.Once{}

	c.stopChan = stopChan
	c.refresherDone = doneChan
	c.stopOnce = stopOnce
	c.refresherRunning = true
	c.refresherLock.Unlock()

	go func() {
		ticker := time.NewTicker(refreshInterval)

		defer func() {
			ticker.Stop()
			c.refresherLock.Lock()
			c.refresherRunning = false
			c.stopChan = nil
			c.refresherDone = nil
			c.stopOnce = nil
			c.refresherLock.Unlock()
			close(doneChan)
		}()

		for {
			select {
			case <-ticker.C:
				if err := c.Authenticate(); err != nil {
					go c.errorCallback(err)
				}
			case <-stopChan:
				return
			}
		}
	}()
	return nil
}

func (client *zabbixClient) StopTokenRefresher() {
	client.refresherLock.Lock()
	if !client.refresherRunning {
		client.refresherLock.Unlock()
		return
	}

	stopChan := client.stopChan
	doneChan := client.refresherDone
	stopOnce := client.stopOnce
	client.refresherLock.Unlock()

	stopOnce.Do(func() {
		close(stopChan)
	})
	<-doneChan
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

	reqBody, err := json.Marshal(request)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json-rpc")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

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
