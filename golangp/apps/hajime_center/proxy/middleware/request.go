package middleware

import (
	"encoding/json"
	"fmt"
	"hajime/golangp/common/dify"
	"hajime/golangp/common/logging"
	"io/ioutil"
	"net/http"
)

func FetchInstalledApps(r *http.Request) ([]InstalledApps, error) {
	client := &http.Client{}

	difyClient, err := dify.GetDifyClient()
	if err != nil {
		logging.Warning("Auth Failed: " + err.Error())
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	Token, err := difyClient.GetUserToken("admin")
	if err != nil {
		logging.Warning("failed to unauthorized: " + err.Error())
		return nil, fmt.Errorf("failed to unauthorized: %w", err)
	}

	// Construct the URL using the scheme and host
	url := fmt.Sprintf("%s/installed-apps", difyClient.ConsoleHost)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add the authorization token
	req.Header.Set("Authorization", `Bearer `+Token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var installAppResponse InstallAppResponse
	if err := json.Unmarshal(body, &installAppResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return installAppResponse.InstalledApps, nil
}
