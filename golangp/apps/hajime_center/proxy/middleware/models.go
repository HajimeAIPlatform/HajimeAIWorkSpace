package middleware

import "hajime/golangp/apps/hajime_center/models"

// OriginalResponse defines the structure for the original response
type OriginalResponse struct {
	Data    []map[string]interface{} `json:"data"`
	Limit   int                      `json:"limit"`
	Total   int                      `json:"total"`
	Page    int                      `json:"page"`
	HasMore bool                     `json:"has_more"`
}

type InstallAppResponse struct {
	InstalledApps []InstalledApps `json:"installed_apps"`
}

type InstalledApps struct {
	ID               string            `json:"id"`
	App              models.HajimeApps `json:"app"`
	AppOwnerTenantID string            `json:"app_owner_tenant_id"`
	IsPinned         bool              `json:"is_pinned"`
	LastUsedAt       int64             `json:"last_used_at"`
	Editable         bool              `json:"editable"`
	Uninstallable    bool              `json:"uninstallable"`
}

type RecommendedAPP struct {
	App struct {
		Icon           string `json:"icon"`
		IconBackground string `json:"icon_background"`
		ID             string `json:"id"`
		Mode           string `json:"mode"`
		Name           string `json:"name"`
	} `json:"app"`
	ID               string      `json:"id"`
	AppID            string      `json:"app_id"`
	Category         string      `json:"category"`
	Copyright        interface{} `json:"copyright"`
	CustomDisclaimer interface{} `json:"custom_disclaimer"`
	Description      interface{} `json:"description"`
	IsListed         bool        `json:"is_listed"`
	Position         int64       `json:"position"`
	PrivacyPolicy    interface{} `json:"privacy_policy"`
}

type NoAuthApp struct {
	Categories     []string         `json:"categories"`
	RecommendedAPP []RecommendedAPP `json:"recommended_apps"`
}
