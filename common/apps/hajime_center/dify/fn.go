package dify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func SendGetRequest(forConsole bool, dc *DifyClient, api string, header map[string]string) (httpCode int, bodyText []byte, err error) {
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return -1, nil, err
	}

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	if forConsole {
		setConsoleAuthorization(dc, req)
	} else {
		setAPIAuthorization(dc, req)
	}

	resp, err := dc.Client.Do(req)
	if err != nil {
		return -1, nil, err
	}
	defer resp.Body.Close()

	bodyText, err = io.ReadAll(resp.Body)
	return resp.StatusCode, bodyText, err
}

func SendPostRequest(forConsole bool, dc *DifyClient, api string, postBody interface{}, header map[string]string) (httpCode int, bodyText []byte, err error) {
	var payload *strings.Reader
	if postBody != nil {
		buf, err := json.Marshal(postBody)
		if err != nil {
			return -1, nil, err
		}
		payload = strings.NewReader(string(buf))
	} else {
		payload = nil
	}

	req, err := http.NewRequest("POST", api, payload)
	if err != nil {
		return -1, nil, err
	}

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	if forConsole {
		setConsoleAuthorization(dc, req)
	} else {
		setAPIAuthorization(dc, req)
	}

	resp, err := dc.Client.Do(req)
	if err != nil {
		return -1, nil, err
	}
	defer resp.Body.Close()

	bodyText, err = io.ReadAll(resp.Body)
	return resp.StatusCode, bodyText, err
}

func SendPutRequest(forConsole bool, dc *DifyClient, api string, putBody interface{}, header map[string]string) (httpCode int, bodyText []byte, err error) {
	var payload *strings.Reader
	if putBody != nil {
		buf, err := json.Marshal(putBody)
		if err != nil {
			return -1, nil, err
		}
		payload = strings.NewReader(string(buf))
	} else {
		payload = nil
	}

	req, err := http.NewRequest("PUT", api, payload)
	if err != nil {
		return -1, nil, err
	}

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	if forConsole {
		setConsoleAuthorization(dc, req)
	} else {
		setAPIAuthorization(dc, req)
	}

	resp, err := dc.Client.Do(req)
	if err != nil {
		return -1, nil, err
	}
	defer resp.Body.Close()

	bodyText, err = io.ReadAll(resp.Body)
	return resp.StatusCode, bodyText, err
}

func SendDeleteRequest(forConsole bool, dc *DifyClient, api string) (httpCode int, bodyText []byte, err error) {
	req, err := http.NewRequest("DELETE", api, nil)
	if err != nil {
		return -1, nil, err
	}

	if forConsole {
		setConsoleAuthorization(dc, req)
	} else {
		setAPIAuthorization(dc, req)
	}

	resp, err := dc.Client.Do(req)
	if err != nil {

		return -1, nil, err
	}
	defer resp.Body.Close()

	bodyText, err = io.ReadAll(resp.Body)
	return resp.StatusCode, bodyText, err
}

// CommonRiskForSendRequest 检查状态码是否在 2xx 范围内，如果不是则返回错误

func CommonRiskForSendRequest(code int, err error) error {
	if err != nil {
		return err
	}

	// 检查状态码是否在 2xx 范围内
	if code < 200 || code >= 300 {
		return fmt.Errorf("status code: %d", code)
	}
	return nil
}

func CommonRiskForSendRequestWithCode(code int, err error, targetCode int) error {
	if err != nil {
		return err
	}

	if code != targetCode {
		return fmt.Errorf("status code: %d", code)
	}

	return nil
}
