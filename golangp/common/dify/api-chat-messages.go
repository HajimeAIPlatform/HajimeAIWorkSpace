package dify

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ChatMessagesPayload struct {
	Inputs         map[string]interface{}    `json:"inputs"`
	Query          string                    `json:"query"`
	ResponseMode   string                    `json:"response_mode,omitempty"`
	ConversationID string                    `json:"conversation_id,omitempty"`
	User           string                    `json:"user ,omitempty"`
	Files          []ChatMessagesPayloadFile `json:"files,omitempty"`
}

type ChatMessagesPayloadFile struct {
	Type           string `json:"type"`
	TransferMethod string `json:"transfer_method"`
	URL            string `json:"url,omitempty"`
	UploadFileID   string `json:"upload_file_id,omitempty"`
}

type ChatMessagesResponse struct {
	Event          string `json:"event"`
	MessageID      string `json:"message_id"`
	ConversationID string `json:"conversation_id"`
	Mode           string `json:"mode"`
	Answer         string `json:"answer"`
	Metadata       any    `json:"metadata"`
	CreatedAt      int    `json:"created_at"`
}

type ChatMessagesStopResponse struct {
	Result string `json:"result"`
}

func PrepareChatPayload(payload map[string]interface{}) (string, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (dc *DifyClient) ChatMessages(query string, inputs map[string]interface{}, conversation_id string, files []ChatMessagesPayloadFile, appAuthorization string) (result ChatMessagesResponse, err error) {
	var payload ChatMessagesPayload

	//if len(inputs) == 0 {
	//	return result, fmt.Errorf("inputs is required")
	//} else {
	//	var tryDecode map[string]interface{}
	//	err := json.Unmarshal([]byte(inputs), &tryDecode)
	//	if err != nil {
	//		return result, fmt.Errorf("inputs should be a valid JSON string")
	//	}
	//	payload.Inputs = tryDecode
	//}
	if inputs != nil {
		payload.Inputs = inputs
	} else {
		payload.Inputs = make(map[string]interface{}) // 将 Inputs 设置为空的映射
	}
	if query == "" {
		return result, fmt.Errorf("query should be a valid JSON string")
	} else {
		payload.Query = query
	}

	payload.ResponseMode = RESPONSE_MODE_BLOCKING
	payload.User = dc.User

	if conversation_id != "" {
		payload.ConversationID = conversation_id
	}

	if len(files) > 0 {
		payload.Files = files
	}

	api := dc.GetAPI(API_CHAT_MESSAGES)

	code, body, err := SetAppsPostAuthorization(dc, api, payload, appAuthorization)

	err = CommonRiskForSendRequest(code, err)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal the response: %v", err)
	}
	return result, nil
}

func (dc *DifyClient) ChatMessagesStreaming(query string, inputs map[string]interface{}, conversationID string, files []ChatMessagesPayloadFile, appAuthorization string) (<-chan string, error) {
	var payload ChatMessagesPayload

	if inputs != nil {
		payload.Inputs = inputs
	} else {
		payload.Inputs = make(map[string]interface{})
	}
	if query == "" {
		return nil, fmt.Errorf("query should be a valid JSON string")
	} else {
		payload.Query = query
	}

	payload.ResponseMode = "streaming"

	if conversationID != "" {
		payload.ConversationID = conversationID
	}

	if len(files) > 0 {
		payload.Files = files
	}

	api := dc.GetAPI(API_CHAT_MESSAGES)

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", api, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", appAuthorization)

	client := &http.Client{
		Timeout: time.Minute,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received non-200 response code: %d", resp.StatusCode)
	}

	dataChan := make(chan string)
	go func() {
		defer close(dataChan)
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.TrimSpace(line) == "" {
				continue
			} else {
				jsonStr := strings.TrimPrefix(line, "data: ")

				var message ChatMessagesResponse
				err := json.Unmarshal([]byte(jsonStr), &message)
				if err != nil {
					fmt.Println("Error decoding JSON:", err)
					// 发送错误消息到通道
					errorMessage, _ := json.Marshal(map[string]string{
						"error": fmt.Sprintf("error decoding JSON: %v", err),
					})
					dataChan <- string(errorMessage)
					continue
				}

				fmt.Println("Decoded Message:", message)
				// 将解析后的 message 重新编码为 JSON 字符串
				messageJson, err := json.Marshal(message)
				if err != nil {
					fmt.Println("Error encoding JSON:", err)
					// 发送错误消息到通道
					errorMessage, _ := json.Marshal(map[string]string{
						"error": fmt.Sprintf("error encoding JSON: %v", err),
					})
					dataChan <- string(errorMessage)
					continue
				}
				dataChan <- string(messageJson)
			}
			fmt.Println("Received Line:", line)
		}

		if err := scanner.Err(); err != nil {
			errorMessage, _ := json.Marshal(map[string]string{
				"error": fmt.Sprintf("scanner error: %v", err),
			})
			dataChan <- string(errorMessage)
		}
	}()

	return dataChan, nil
}

func (dc *DifyClient) ChatMessagesStop(task_id string, Authorization string) (result ChatMessagesStopResponse, err error) {
	if task_id == "" {
		return result, fmt.Errorf("task_id is required")
	}

	if Authorization == "" {
		return result, fmt.Errorf("Authorization is required")
	}

	api := dc.GetAPI(API_CHAT_MESSAGES_STOP)
	api = UpdateAPIParam(api, API_PARAM_TASK_ID, task_id)

	// Ensure payload is not nil if required by SetAppsPostAuthorization
	payload := map[string]string{} // Add necessary payload data if needed

	code, body, err := SetAppsPostAuthorization(dc, api, payload, Authorization)
	if err != nil {
		return result, fmt.Errorf("failed to send request: %v", err)
	}

	err = CommonRiskForSendRequest(code, err)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal the response: %v", err)
	}
	return result, nil
}
