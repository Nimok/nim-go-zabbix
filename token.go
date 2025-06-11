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

type TokenCreateResponse struct {
	TokenIDs []string `json:"tokenids"`
}

type TokenGenerateParameters []string

type TokenGenerateResponse struct {
	TokenId string `json:"tokenid"`
	Token   string `json:"token"`
}

type TokenDeleteParameters []string

type TokenDeleteResponse struct {
	TokenIDs []string `json:"tokenids"`
}

func (z *zabbixClient) TokenCreate(ctx context.Context, params Token) (*TokenCreateResponse, error) {

	var result TokenCreateResponse

	err := z.makeRequest(ctx, "token.create", params, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (z *zabbixClient) TokenGenerate(ctx context.Context, params TokenGenerateParameters) ([]TokenGenerateResponse, error) {

	var result []TokenGenerateResponse

	err := z.makeRequest(ctx, "token.generate", params, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (z *zabbixClient) TokenDelete(ctx context.Context, params TokenDeleteParameters) (*TokenDeleteResponse, error) {

	var result TokenDeleteResponse

	err := z.makeRequest(ctx, "token.delete", params, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
