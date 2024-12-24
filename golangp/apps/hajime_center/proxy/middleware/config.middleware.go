package middleware

import (
	"bytes"
	"encoding/json"
	"hajime/golangp/apps/hajime_center/dify"
	"hajime/golangp/common/logging"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type ListDatasetsModelConfigDetail struct {
	Provider         string `json:"provider"`
	Name             string `json:"name"`
	Mode             string `json:"mode"`
	CompletionParams struct {
	} `json:"completion_params"`
}

type UpdateAppsModelConfigPayload struct {
	PrePrompt        string `json:"pre_prompt"`
	PromptType       string `json:"prompt_type"`
	ChatPromptConfig struct {
	} `json:"chat_prompt_config"`
	CompletionPromptConfig struct {
	} `json:"completion_prompt_config"`
	UserInputForm        []any  `json:"user_input_form"`
	DatasetQueryVariable string `json:"dataset_query_variable"`
	OpeningStatement     string `json:"opening_statement"`
	SuggestedQuestions   []any  `json:"suggested_questions"`
	MoreLikeThis         struct {
		Enabled bool `json:"enabled"`
	} `json:"more_like_this"`
	SuggestedQuestionsAfterAnswer struct {
		Enabled bool `json:"enabled"`
	} `json:"suggested_questions_after_answer"`
	SpeechToText struct {
		Enabled bool `json:"enabled"`
	} `json:"speech_to_text"`
	TextToSpeech struct {
		Enabled  bool   `json:"enabled"`
		Voice    string `json:"voice"`
		Language string `json:"language"`
	} `json:"text_to_speech"`
	RetrieverResource struct {
		Enabled bool `json:"enabled"`
	} `json:"retriever_resource"`
	SensitiveWordAvoidance struct {
		Enabled bool   `json:"enabled"`
		Type    string `json:"type"`
		Configs []any  `json:"configs"`
	} `json:"sensitive_word_avoidance"`
	AgentMode struct {
		Enabled      bool   `json:"enabled"`
		MaxIteration int    `json:"max_iteration"`
		Strategy     string `json:"strategy"`
		Tools        []any  `json:"tools"`
	} `json:"agent_mode"`
	Model          ListDatasetsModelConfigDetail `json:"model"`
	DatasetConfigs DatasetConfigs                `json:"dataset_configs"`
	FileUpload     FileUploadPayload             `json:"file_upload"`
}

type DatasetArray struct {
	Enabled bool   `json:"enabled"`
	ID      string `json:"id"`
}

type ImagePayload struct {
	Enabled         bool     `json:"enabled"`
	NumberLimits    int      `json:"number_limits"`
	Detail          string   `json:"detail"`
	TransferMethods []string `json:"transfer_methods"`
}

type FileUploadPayload struct {
	Image ImagePayload `json:"image"`
}

type DatasetConfigs struct {
	RetrievalModel string `json:"retrieval_model"`
	Datasets       struct {
		Datasets []struct {
			Dataset DatasetArray `json:"dataset"`
		} `json:"datasets"`
	} `json:"datasets"`
}

func ModelUpdateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		appID := vars["app_id"]
		if appID == "" {
			WriteErrorResponse(w, "400", "app_id is required", http.StatusBadRequest)
			logging.Warning("app_id is required")
			return
		}

		user, err := DeserializeUser(r)
		if err != nil {
			logging.Warning("Auth Failed: " + err.Error())
			WriteErrorResponse(w, "401", err.Error(), http.StatusUnauthorized)
			return
		}

		difyClient, err := dify.GetDifyClient()
		if err != nil {
			logging.Warning("Auth Failed: " + err.Error())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		Token, err := difyClient.GetUserToken(user.Role)
		if err != nil {
			logging.Warning("Failed to get user token: " + err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		r.Header.Set("Authorization", "Bearer "+Token)

		// 读取请求体
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logging.Warning("Unable to read the request body: " + err.Error())
			http.Error(w, "Unable to read the request body", http.StatusBadRequest)
			return
		}

		// 检查请求体是否为空
		if len(bodyBytes) == 0 {
			logging.Warning("The request body is empty")
			http.Error(w, "The request body is empty", http.StatusBadRequest)
			return
		}

		// 重置请求体
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// 解码请求体
		var payload UpdateAppsModelConfigPayload
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			logging.Warning("Failed to decode request body: " + err.Error())
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// 重置请求体以供后续中间件或处理器使用
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		knowledge := false
		variables := false

		// 检查 dataset_configs.datasets.datasets 长度
		if len(payload.DatasetConfigs.Datasets.Datasets) > 0 {
			knowledge = true
		}

		// 检查 user_input_form 长度
		if len(payload.UserInputForm) > 0 {
			variables = true
		}

		err = user.UpdateConfigUsage(appID, knowledge, variables)
		if err != nil {
			logging.Warning("Failed to update config usage: " + err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// 调用下一个处理器
		next.ServeHTTP(w, r)
	})
}
