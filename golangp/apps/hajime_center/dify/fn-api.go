package dify

import (
	"net/http"
)

func setAPIAuthorization(dc *DifyClient, req *http.Request) {
	//req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", dc.Key))
	req.Header.Set("Content-Type", "application/json")
}

func SendGetRequestToAPI(dc *DifyClient, api string, header map[string]string) (httpCode int, bodyText []byte, err error) {
	return SendGetRequest(false, dc, api, header)
}

func SendPostRequestToAPI(dc *DifyClient, api string, postBody interface{}, header map[string]string) (httpCode int, bodyText []byte, err error) {
	return SendPostRequest(false, dc, api, postBody, header)
}

func SendPutRequestToAPI(dc *DifyClient, api string, putBody interface{}, header map[string]string) (httpCode int, bodyText []byte, err error) {
	return SendPutRequest(false, dc, api, putBody, header)
}
