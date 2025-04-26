package zabbix

import "context"

type LogoutSuccess bool

func (z *zabbixClient) Logout(ctx context.Context) (LogoutSuccess, error) {

	var result LogoutSuccess
	params := []string{}

	err := z.makeRequest(ctx, "user.logout", params, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
