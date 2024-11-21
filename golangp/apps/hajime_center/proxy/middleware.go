package proxy

import (
	"hajime/golangp/apps/hajime_center/initializers"
	"hajime/golangp/apps/hajime_center/models"
	"hajime/golangp/apps/hajime_center/proxy/middleware"
	"net/http"
)

func ModifyResponse(w *http.Response, r *http.Request, user models.User) error {

	db := initializers.DB
	if r.URL.Path == "/console/api/apps" || middleware.IsAppIDPath(r.URL.Path) {
		switch r.Method {
		case http.MethodGet:
			return middleware.HandleGetApp(w, r, db, user)
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
