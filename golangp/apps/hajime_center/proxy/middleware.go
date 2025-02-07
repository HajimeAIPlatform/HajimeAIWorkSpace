package proxy

import (
	"hajime/golangp/apps/hajime_center/models"
	"hajime/golangp/apps/hajime_center/proxy/middleware"
	"hajime/golangp/common/initializers"
	"net/http"
	"regexp"
)

func ModifyResponse(w *http.Response, r *http.Request, user models.User) error {

	db := initializers.DB

	// Define regex pattern to match the app_id in the URL
	appIDPattern := `/console/api/installed-apps/([0-9a-fA-F-]+)/conversations`
	re := regexp.MustCompile(appIDPattern)

	// Check if the URL matches the app_id pattern
	matches := re.FindStringSubmatch(r.URL.Path)
	if len(matches) == 2 {
		switch r.Method {
		case http.MethodGet:
			return middleware.HandleGetConversation(w, r, db, user)
		// Add other methods if needed
		default:
			return nil
		}
	}

	if r.URL.Path == "/console/api/apps" || middleware.IsAppIDPath(r.URL.Path) {
		switch r.Method {
		case http.MethodGet:
			return middleware.HandleGetApp(w, r, user)
		case http.MethodPost:
			return middleware.HandlePostApp(w, r, db, user)
		case http.MethodPut:
			return middleware.HandlePutApp(w, r, db, user)
		case http.MethodDelete:
			return middleware.HandleDeleteApp(w, r, db, user)
		default:
			return nil
		}
	}

	if r.URL.Path == "/console/api/installed-apps" {
		switch r.Method {
		case http.MethodGet:
			return middleware.HandleGetInstallApp(w, r, db, user)
		default:
			return nil
		}
	}

	if r.URL.Path == "/console/api/datasets/init" {
		switch r.Method {
		case http.MethodPost:
			return middleware.HandlePostDatasetInit(w, r, db, user)
		default:
			return nil
		}
	}
	if r.URL.Path == "/console/api/datasets" {
		switch r.Method {
		case http.MethodGet:
			return middleware.HandleGetAllDatasets(w, r, user)
		default:
			return nil
		}
	}

	if middleware.IsDatasetIDPath(r.URL.Path) {
		switch r.Method {
		case http.MethodDelete:
			return middleware.HandleDeleteDatasets(w, r, user)
		default:
			return nil
		}
	}

	return nil
}
