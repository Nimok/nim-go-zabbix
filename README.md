# nim-go-zabbix


## Overview

GO package that enables easy integration with Zabbix API v7.
Made to be useful in both long running applications and scripts.

## Install

```
go get github.com/nimok/nim-go-zabbix
```

## Features

Supports different authentication methods using functional options on the client.
Example:

```go
// Create client with username/password as authentication method
client, err := zabbix.NewClient("http://<your-zabbix-server>/api_jsonrpc.php",
    zabbix.WithUserPass("bestUsername", "excellentPassword"))
if err != nil {
    log.Fatal(err)
}
```

```go
// Create client with an api token
client, err := zabbix.NewClient("http://<your-zabbix-server>/api_jsonrpc.php",
    zabbix.WithAPIToken("someapitoken"))
if err != nil {
    log.Fatal(err)
}
```

Supports refreshing on underlaying Bearar Token on specified interval.

```go
client, err := zabbix.NewClient("http://<your-zabbix-server>/api_jsonrpc.php",
    zabbix.WithAPIToken("someapitoken"),
)
if err != nil {
    log.Fatal(err)
}

// Login
err = client.Authenticate()
if err != nil {
    log.Fatal(err)
}

err = client.StartTokenRefresher(time.Minute * 10) //Refresh the bearer token every 10 minutes
if err != nil {
    log.Fatal(err)
}
```

Use a custom callback for errors:

```go
client, err := zabbix.NewClient("http://<your-zabbix-server>/api_jsonrpc.php",
    zabbix.WithAPIToken("someapitoken"),
    zabbix.WithErrorCallback(func(err error) { // Setup custom error handling
        // Do what you want with the errors
        log.Println(err)
    }),
)
```

## Quickstart

```go 
package main

import (
	"context"
	"log"

	zabbix "github.com/nimok/nim-go-zabbix"
)

func main() {
	// Create context for usage in requests
	ctx := context.Background()

	// Create client with username/password as authentication method
	client, err := zabbix.NewClient("http://<your-zabbix-server>/api_jsonrpc.php",
		zabbix.WithUserPass("bestUsername", "excellentPassword"))
	if err != nil {
		log.Fatal(err)
	}

	// Login
	err = client.Authenticate()
	if err != nil {
		log.Fatal(err)
	}

	// Setup filter based GET
	filter := make(map[string]any, 0)
	filter["host"] = "Zabbix server"

	hosts, err := client.HostGet(ctx, zabbix.HostGetParameters{
		GetParameters: zabbix.GetParameters{
			Filter: filter,
			Output: "extend",
		},
	})
	if err != nil {
		log.Println(err)
	}

	// Do something with your host/hosts
	for _, host := range hosts {
		log.Println(host.HostID)
	}

	_, err := client.Logout(ctx)
	if err != nil {
		log.Println(err)
	}
}

```

## License

nim-go-zabbix is released under the [MIT License](LICENSE).