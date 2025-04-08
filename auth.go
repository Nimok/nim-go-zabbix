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

func (client *zabbixClient) Authenticate() error {
	authReq := AuthRequest{
		JSONRPC: "2.0",
		Method:  "user.login",
		Params: map[string]string{
			"username": client.username,
			"password": client.password,
		},
		ID: 1,
	}

	reqBody, err := json.Marshal(authReq)
	if err != nil {
		return err
	}

	resp, err := http.Post(client.url, "application/json-rpc", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("auth request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var authResp authResponse
	err = json.Unmarshal(body, &authResp)
	if err != nil {
		return err
	}

	if authResp.Error != nil {
		return fmt.Errorf("auth failed: %s", authResp.Error.Data)
	}

	client.bearerTokenLock.Lock()
	client.bearerToken = authResp.Result
	client.bearerTokenLock.Unlock()
	return nil
}
