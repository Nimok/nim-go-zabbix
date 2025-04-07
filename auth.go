package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AuthRequest struct {
	JSONRPC string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  map[string]string `json:"params"`
	ID      int               `json:"id"`
}

type authResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  string `json:"result"`
	ID      int    `json:"id"`
	Error   *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	} `json:"error,omitempty"`
}

func (api *zabbixClient) Authenticate() error {
	authReq := AuthRequest{
		JSONRPC: "2.0",
		Method:  "user.login",
		Params: map[string]string{
			"username": api.username,
			"password": api.password,
		},
		ID: 1,
	}

	reqBody, _ := json.Marshal(authReq)
	resp, err := http.Post(api.url, "application/json-rpc", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("auth request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var authResp authResponse
	json.Unmarshal(body, &authResp)

	if authResp.Error != nil {
		return fmt.Errorf("auth failed: %s", authResp.Error.Data)
	}

	api.bearerTokenLock.Lock()
	api.bearerToken = authResp.Result
	api.bearerTokenLock.Unlock()
	return nil
}
