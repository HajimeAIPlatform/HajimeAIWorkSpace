package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hajime/golangp/apps/hajime_center/dify"
	"hajime/golangp/common/logging"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// Graph represents the entire graph structure
type Graph struct {
	Nodes    []Node   `json:"nodes"`
	Edges    []Edge   `json:"edges"`
	Viewport Viewport `json:"viewport"`
}

// Node represents a single node in the graph
type Node struct {
	ID               string   `json:"id"`
	Type             string   `json:"type"`
	Data             NodeData `json:"data"`
	Position         Position `json:"position"`
	TargetPosition   string   `json:"targetPosition"`
	SourcePosition   string   `json:"sourcePosition"`
	PositionAbsolute Position `json:"positionAbsolute"`
	Width            int      `json:"width"`
	Height           int      `json:"height"`
	Selected         bool     `json:"selected"`
}

// NodeData represents the data contained within a node
type NodeData struct {
	Type               string              `json:"type"`
	Title              string              `json:"title"`
	Desc               string              `json:"desc"`
	Variables          []interface{}       `json:"variables"`
	Selected           bool                `json:"selected"`
	ToolParameters     *ToolParameters     `json:"tool_parameters,omitempty"`
	ToolConfigurations *ToolConfigurations `json:"tool_configurations,omitempty"`
	ProviderID         string              `json:"provider_id,omitempty"`
	ProviderType       string              `json:"provider_type,omitempty"`
	ProviderName       string              `json:"provider_name,omitempty"`
	ToolName           string              `json:"tool_name,omitempty"`
	ToolLabel          string              `json:"tool_label,omitempty"`
	Outputs            []Output            `json:"outputs,omitempty"`
}

// ToolParameters represents the parameters for a tool node
type ToolParameters struct {
	Content Content `json:"content"`
}

// Content represents the content of a tool parameter
type Content struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// ToolConfigurations represents the configurations for a tool node
type ToolConfigurations struct {
	ErrorCorrection string `json:"error_correction"`
	Border          int    `json:"border"`
}

// Output represents an output variable from a node
type Output struct {
	Variable      string   `json:"variable"`
	ValueSelector []string `json:"value_selector"`
}

// Position represents the position of a node
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Edge represents a connection between two nodes
type Edge struct {
	ID           string   `json:"id"`
	Type         string   `json:"type"`
	Source       string   `json:"source"`
	SourceHandle string   `json:"sourceHandle"`
	Target       string   `json:"target"`
	TargetHandle string   `json:"targetHandle"`
	Data         EdgeData `json:"data"`
	ZIndex       int      `json:"zIndex"`
}

// EdgeData represents the data contained within an edge
type EdgeData struct {
	SourceType    string `json:"sourceType"`
	TargetType    string `json:"targetType"`
	IsInIteration bool   `json:"isInIteration"`
}

// Viewport represents the viewport settings of the graph
type Viewport struct {
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Zoom float64 `json:"zoom"`
}

// Features represents the features configuration
type Features struct {
	OpeningStatement              string                        `json:"opening_statement"`
	SuggestedQuestions            []interface{}                 `json:"suggested_questions"`
	SuggestedQuestionsAfterAnswer SuggestedQuestionsAfterAnswer `json:"suggested_questions_after_answer"`
	TextToSpeech                  TextToSpeech                  `json:"text_to_speech"`
	SpeechToText                  SpeechToText                  `json:"speech_to_text"`
	RetrieverResource             RetrieverResource             `json:"retriever_resource"`
	SensitiveWordAvoidance        SensitiveWordAvoidance        `json:"sensitive_word_avoidance"`
	FileUpload                    FileUpload                    `json:"file_upload"`
}

// SuggestedQuestionsAfterAnswer represents the configuration for suggested questions after an answer
type SuggestedQuestionsAfterAnswer struct {
	Enabled bool `json:"enabled"`
}

// TextToSpeech represents the text-to-speech configuration
type TextToSpeech struct {
	Enabled  bool   `json:"enabled"`
	Voice    string `json:"voice"`
	Language string `json:"language"`
}

// SpeechToText represents the speech-to-text configuration
type SpeechToText struct {
	Enabled bool `json:"enabled"`
}

// RetrieverResource represents the retriever resource configuration
type RetrieverResource struct {
	Enabled bool `json:"enabled"`
}

// SensitiveWordAvoidance represents the sensitive word avoidance configuration
type SensitiveWordAvoidance struct {
	Enabled bool `json:"enabled"`
}

// FileUpload represents the file upload configuration
type FileUpload struct {
	Image                    Image            `json:"image"`
	Enabled                  bool             `json:"enabled"`
	AllowedFileTypes         []string         `json:"allowed_file_types"`
	AllowedFileExtensions    []string         `json:"allowed_file_extensions"`
	AllowedFileUploadMethods []string         `json:"allowed_file_upload_methods"`
	NumberLimits             int              `json:"number_limits"`
	FileUploadConfig         FileUploadConfig `json:"fileUploadConfig"`
}

// Image represents the image upload configuration
type Image struct {
	Enabled         bool     `json:"enabled"`
	NumberLimits    int      `json:"number_limits"`
	TransferMethods []string `json:"transfer_methods"`
}

// FileUploadConfig represents the file upload configuration limits
type FileUploadConfig struct {
	FileSizeLimit      int `json:"file_size_limit"`
	BatchCountLimit    int `json:"batch_count_limit"`
	ImageFileSizeLimit int `json:"image_file_size_limit"`
	VideoFileSizeLimit int `json:"video_file_size_limit"`
	AudioFileSizeLimit int `json:"audio_file_size_limit"`
}

// Main data structure containing all components
type MainData struct {
	Graph                 Graph         `json:"graph"`
	Features              Features      `json:"features"`
	EnvironmentVariables  []interface{} `json:"environment_variables"`
	ConversationVariables []interface{} `json:"conversation_variables"`
	Hash                  string        `json:"hash"`
}

func containsToolNode(nodes []Node) bool {
	for _, node := range nodes {
		if node.Data.Type == "tool" {
			return true
		}
	}
	return false
}

func WorkflowDraftMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		appID := vars["app_id"]
		fmt.Println("appID: ", appID)
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
		var payload MainData
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			logging.Warning("Failed to decode request body: " + err.Error())
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if containsToolNode(payload.Graph.Nodes) {
			err := user.UpdateWorkflowDraftUsage(appID, true)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			fmt.Println("The graph does not contain any nodes of type 'tool'.")
		}

		// 重置请求体以供后续中间件或处理器使用
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// 调用下一个处理器
		next.ServeHTTP(w, r)
	})
}
