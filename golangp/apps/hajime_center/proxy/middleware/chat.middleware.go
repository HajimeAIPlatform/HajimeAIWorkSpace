package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hajime/golangp/apps/hajime_center/constants"
	"hajime/golangp/apps/hajime_center/dify"
	"hajime/golangp/apps/hajime_center/models"
	"hajime/golangp/common/logging"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type ResponseWriterInterceptor struct {
	http.ResponseWriter
	Body *bytes.Buffer
}

func (rwi *ResponseWriterInterceptor) Write(b []byte) (int, error) {
	rwi.Body.Write(b) // 捕获响应体
	return rwi.ResponseWriter.Write(b)
}

func ChatMessageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r) // Assuming you're using Gorilla Mux
		appID := vars["app_id"]
		if appID == "" {
			WriteErrorResponse(w, "400", "app_id is required", http.StatusBadRequest)
		}
		user, err := DeserializeUser(r)
		if err != nil {
			logging.Warning("Auth Failed: " + err.Error())
			WriteErrorResponse(w, "401", err.Error(), http.StatusBadRequest)
			return
		}
		exceedsLimit, err := user.UpdateAppUsage(appID)

		if err != nil {
			logging.Warning("Failed to update app usage: " + err.Error())
			WriteErrorResponse(w, "200", err.Error(), http.StatusBadRequest)
			return
		}

		difyClient, err := dify.GetDifyClient()
		if err != nil {
			logging.Warning("Auth Failed: " + err.Error())
			WriteErrorResponse(w, "401", err.Error(), http.StatusBadRequest)
			return
		}

		Token, err := difyClient.GetUserToken(user.Role)
		if err != nil {
			logging.Warning("Token retrieval failed: " + err.Error())
			WriteErrorResponse(w, "401", err.Error(), http.StatusBadRequest)
			return
		}

		// Create a new ResponseWriterInterceptor
		rwi := &ResponseWriterInterceptor{
			ResponseWriter: w,
			Body:           &bytes.Buffer{},
		}

		r.Header.Set("Authorization", "Bearer "+Token)
		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(rwi, r)
		// Intercept the response here
		processResponse(rwi.Body.Bytes(), exceedsLimit, user, appID)

	})

}
func processResponse(body []byte, exceedsLimit bool, user *models.User, appID string) {
	// 定义结构体来匹配 JSON 数据
	type Usage struct {
		TotalTokens int `json:"total_tokens"`
	}

	type Metadata struct {
		Usage Usage `json:"usage"`
	}

	type Response struct {
		Event          string   `json:"event"`
		Metadata       Metadata `json:"metadata"`
		ConversationID string   `json:"conversation_id"`
	}

	data := string(body)

	// 分割数据行
	lines := strings.Split(data, "\n")

	// 找到最后一个以 "data:" 开头的行
	var lastDataLine string
	for _, line := range lines {
		if strings.HasPrefix(line, "data:") {
			lastDataLine = strings.TrimPrefix(line, "data:")
		}
	}

	// 解析 JSON 数据
	var response Response
	err := json.Unmarshal([]byte(lastDataLine), &response)
	if err != nil {
		logging.Warning("JSON parsing error: %v\nRaw data: %s", err, lastDataLine)
		return
	}

	// 只处理 "event": "message_end" 的情况
	if response.Event == "message_end" {
		tokenCost := float64(response.Metadata.Usage.TotalTokens) * constants.ChatCostPerToken

		if exceedsLimit {
			err := user.UpdateBalance(tokenCost, "ChatCostPerToken")
			if err != nil {
				logging.Warning("Failed to update user balance: " + err.Error())
			}
		}

		tokenEarn := float64(response.Metadata.Usage.TotalTokens) * constants.UsageEarnPerToken

		ownerUser, err := models.GetUserByAppID(appID)
		if err != nil {
			logging.Warning("Failed to get owner user: " + err.Error())
			return
		}

		err = ownerUser.UpdateBalance(tokenEarn, "UsageEarnPerToken")
		if err != nil {
			logging.Warning("Failed to update user balance: " + err.Error())
		}

		_, err = models.CreateConversation(response.ConversationID, user.ID.String())
		if err != nil {
			logging.Warning("Failed to create conversation: " + err.Error())
		}

		logging.Info(fmt.Sprintf(
			"User ID: %s, App ID: %s, Total Tokens: %d, Balance Deducted: %.2f, More than the first three free: %t",
			user.ID, appID, response.Metadata.Usage.TotalTokens, tokenCost, exceedsLimit,
		))
	}
}
