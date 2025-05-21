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

type TokenGenerateParameters []string

type tokenGenerateResponse struct {
	TokenId string `json:"tokenid"`
	Token   string `json:"token"`
}

type TokenDeleteParameters []string

type tokenDeleteResponse struct {
	TokenIDs []string `json:"tokenids"`
}

func (z *zabbixClient) TokenCreate(ctx context.Context, params Token) (*tokenCreateResponse, error) {

	var result tokenCreateResponse

	err := z.makeRequest(ctx, "token.create", params, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (z *zabbixClient) TokenGenerate(ctx context.Context, params TokenGenerateParameters) ([]tokenGenerateResponse, error) {

	var result []tokenGenerateResponse

	err := z.makeRequest(ctx, "token.generate", params, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (z *zabbixClient) TokenDelete(ctx context.Context, params TokenDeleteParameters) (*tokenDeleteResponse, error) {

	var result tokenDeleteResponse

	err := z.makeRequest(ctx, "token.delete", params, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
