package zabbix

import (
	"encoding/json"
	"fmt"
)

type apiVersionRespone struct {
	JSONRPC string `json:"jsonrpc"`
	Result  string `json:"result"`
	ID      int    `json:"id"`
}

func (z *ZabbixClient) GetApiVersion() (*apiVersionRespone, error) {
	data, err := z.makeRequest("apiinfo.version", map[string]any{})
	if err != nil {
		fmt.Println("API call failed:", err)
	}

	var resp apiVersionRespone

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, err
}
