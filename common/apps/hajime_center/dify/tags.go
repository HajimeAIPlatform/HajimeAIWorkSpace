package dify

import (
	"encoding/json"
	"fmt"
	"strings"
)

type TagsModel struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type TagsResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	BindingCount string `json:"binding_count"` // 修改为 string 类型
}

type TagsBindingPayload struct {
	TagIds   []string `json:"tag_ids"`
	TargetID string   `json:"target_id"`  // app_id
	Type     string   `json:"type"`
}

func (dc *DifyClient) GetTagsList() (result []TagsResponse, err error) {
	api := dc.GetConsoleAPI(CONSOLE_API_APPS_TAGS_GET)

	code, body, err := SendGetRequestToConsole(dc, api)

	err = CommonRiskForSendRequest(code, err)
	if err != nil {
		fmt.Println("error: ", string(body))
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal the response: %v", err)
	}
	return result, nil
}

func (dc *DifyClient) CreateTag(name string) (result TagsResponse, err error) {
	payload := TagsModel{
		Name: name,
		Type: "app",
	}

	api := dc.GetConsoleAPI(CONSOLE_API_APPS_TAGS_CREATE)
	code, body, err := SendPostRequestToConsole(dc, api, payload)

	err = CommonRiskForSendRequest(code, err)
	if err != nil {
		fmt.Println("error: ", string(body))
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal the response: %v", err)
	}
	return result, nil
}

func (dc *DifyClient) InitTag() {
	// 获取现有的标签列表
	existingTags, err := dc.GetTagsList()
	if err != nil {
		fmt.Println("error fetching tags:", err)
		return
	}

	// 创建一个map来快速查找现有标签
	existingTagNames := make(map[string]bool)
	for _, tag := range existingTags {
		existingTagNames[tag.Name] = true
	}

	// 定义需要初始化的标签列表
	var tagList = []string{"assistant", "knowledge"}

	for _, tag := range tagList {
		// 检查标签是否已经存在
		if existingTagNames[tag] {
			fmt.Println("tag already exists:", tag)
			continue
		}

		// 创建新标签
		result, err := dc.CreateTag(tag)
		if err != nil {
			fmt.Println("error creating tag:", err)
		} else {
			fmt.Println("created tag:", result)
		}
	}
}

func (dc *DifyClient) HandleTagsBinding(appID string, appType string) (err error) {
	// 获取现有的标签列表
	existingTags, err := dc.GetTagsList()
	if err != nil {
		return fmt.Errorf("error fetching tags: %v", err)
	}

	// 查找与 appType 对应的标签 ID
	var tagID string
	for _, tag := range existingTags {
		if strings.EqualFold(tag.Name, appType) {
			tagID = tag.ID
			break
		}
	}
	if tagID == "" {
		return fmt.Errorf("tag not found for appType: %s", appType)
	}

	// 创建绑定请求的 payload
	payload := TagsBindingPayload{
		TagIds:   []string{tagID},
		TargetID: appID,
		Type:     "app", // 假设类型为 "app"
	}

	// 发送绑定请求
	api := dc.GetConsoleAPI(CONSOLE_API_APPS_TAGS_BINDINGS)
	code, body, err := SendPostRequestToConsole(dc, api, payload)
	if err != nil {
		return fmt.Errorf("error sending bind request: %v", err)
	}

	err = CommonRiskForSendRequest(code, err)
	if err != nil {
		fmt.Println("error: ", string(body))
		return err
	}

	// 检查 HTTP 状态码
	if code != 200 {
		return fmt.Errorf("failed to bind tags, received status code: %d", code)
	}

	fmt.Println("tags successfully bound")
	return nil
}

