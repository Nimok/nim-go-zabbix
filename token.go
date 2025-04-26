package zabbix

import "context"

type Token struct {
	TokenID    string `json:"tokenid,omitempty"`
	Name       string `json:"name,omitempty"`
	UserID     string `json:"userid,omitempty"`
	Token      string `json:"token,omitempty"`
	Status     int    `json:"status,omitempty"`
	LastAccess int64  `json:"lastaccess,omitempty"`
	ExpiresAt  int64  `json:"expires_at,omitempty"`
}

type tokenCreateResponse struct {
	TokenIDs []string `json:"tokenids"`
}

type tokenGenerateParameters []string

type tokenGenerateResponse struct {
	TokenId string `json:"tokenid"`
	Token   string `json:"token"`
}

func (z *zabbixClient) TokenCreate(ctx context.Context, params Token) (*tokenCreateResponse, error) {

	var result tokenCreateResponse

	err := z.makeRequest(ctx, "token.create", params, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (z *zabbixClient) TokenGenerate(ctx context.Context, params tokenGenerateParameters) ([]tokenGenerateResponse, error) {

	var result []tokenGenerateResponse

	err := z.makeRequest(ctx, "token.generate", params, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
