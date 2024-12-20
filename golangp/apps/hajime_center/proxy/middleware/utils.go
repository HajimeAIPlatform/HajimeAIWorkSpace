package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hajime/golangp/apps/hajime_center/models"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func IsAppIDPath(path string) bool {
	// 匹配 "/console/api/apps/{app_id}"，确保后面没有其他路径
	matched, _ := regexp.MatchString(`^/console/api/apps/[a-fA-F0-9\-]+/?$`, path)
	return matched
}

func IsDatasetIDPath(path string) bool {
	// 匹配 "/console/api/datasets/{dataset_id}"，确保后面没有其他路径
	matched, _ := regexp.MatchString(`^/console/api/datasets/[a-fA-F0-9\-]+/?$`, path)
	return matched
}

func ReadResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(resp.Body)

	return buf.Bytes(), err
}

func WriteResponseBody(resp *http.Response, body []byte) {
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))
	resp.ContentLength = int64(len(body))
	resp.Header.Set("Content-Length", fmt.Sprint(len(body)))
}

func StructToMap(obj interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &result)
	return result, err
}

func MapToStruct(data map[string]interface{}, result interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, result)
}

func MapToStructApps(data models.HajimeApps, result interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, result)
}

func WriteErrorResponse(w http.ResponseWriter, code, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := map[string]interface{}{
		"code":    code,
		"message": message,
		"status":  status,
	}
	json.NewEncoder(w).Encode(response)
}

func IsPathExcluded(path string, excludedPaths []string, prefix string) bool {
	for _, excludedPath := range excludedPaths {
		if path == prefix+excludedPath {
			return true
		}
	}
	return false
}

func IsPathPrefix(path string, excludedPaths []string, prefix string) bool {
	for _, excludedPath := range excludedPaths {
		if strings.HasPrefix(path, prefix+excludedPath) {
			return true
		}
	}
	return false
}
