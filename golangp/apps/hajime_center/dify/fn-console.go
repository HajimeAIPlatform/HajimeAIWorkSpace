package dify

import (
	"fmt"
	"net/http"
)

func setConsoleAuthorization(dc *DifyClient, req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", dc.ConsoleToken))
	req.Header.Set("Content-Type", "application/json")
}

func SendGetRequestToConsole(dc *DifyClient, api string) (httpCode int, bodyText []byte, err error) {
	return SendGetRequest(true, dc, api, nil)
}

func SendPostRequestToConsole(dc *DifyClient, api string, postBody interface{}) (httpCode int, bodyText []byte, err error) {
	return SendPostRequest(true, dc, api, postBody, nil)
}

func SendPutRequestToConsole(dc *DifyClient, api string, putBody interface{}) (httpCode int, bodyText []byte, err error) {
	return SendPutRequest(true, dc, api, putBody, nil)
}

func SendDeleteRequestToConsole(dc *DifyClient, api string) (httpCode int, bodyText []byte, err error) {
	return SendDeleteRequest(true, dc, api)
}
