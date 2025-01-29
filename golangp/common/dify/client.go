package dify

import (
	"crypto/tls"
	"fmt"
	"hajime/golangp/common/initializers"
	"hajime/golangp/common/logging"
	"net/http"
	"strings"
	"time"
)

type DifyClientConfig struct {
	Key         string
	Host        string
	HostUrl     string
	ConsoleHost string
	Timeout     int
	SkipTLS     bool
	User        string
}

type DifyClient struct {
	Key          string
	Host         string
	HostUrl      string
	ConsoleHost  string
	ConsoleToken string
	Timeout      time.Duration
	SkipTLS      bool
	Client       *http.Client
	User         string
}

func CreateDifyClient(config DifyClientConfig) (*DifyClient, error) {
	cnf, err := initializers.LoadEnv(".")
	if err != nil {
		return nil, err
	}
	key := cnf.DifyApiKey
	if key == "" {
		return nil, fmt.Errorf("dify API Key is required")
	}

	host := cnf.DifyHost
	if host == "" {
		return nil, fmt.Errorf("dify Host is required")
	}

	consoleURL := host + "/console/api"

	timeout := 0 * time.Second
	if config.Timeout <= 0 {
		if config.Timeout < 0 {
			fmt.Println("Timeout should be a positive number, reset to default value: 10s")
		}
		timeout = DEFAULT_TIMEOUT * time.Second
	}

	skipTLS := false
	if config.SkipTLS {
		skipTLS = true
	}

	config.User = strings.TrimSpace(config.User)
	if config.User == "" {
		config.User = DEFAULT_USER
	}

	var client *http.Client

	if skipTLS {
		client = &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}
	} else {
		client = &http.Client{}
	}

	if timeout > 0 {
		client.Timeout = timeout
	}

	return &DifyClient{
		Key:         key,
		Host:        host,
		HostUrl:     host + "/api",
		ConsoleHost: consoleURL,
		Timeout:     timeout,
		SkipTLS:     skipTLS,
		Client:      client,
		User:        config.User,
	}, nil
}

func GetDifyClient() (*DifyClient, error) { // 修改返回类型为 (*DifyClient, error)
	client, err := CreateDifyClient(DifyClientConfig{})
	if err != nil {
		logging.Warning("failed to create DifyClient: %v\n", err)
		return nil, err // 返回 nil 和 err
	}

	fmt.Println(client)

	return client, nil
}
