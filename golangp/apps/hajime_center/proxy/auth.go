package proxy

import (
	"context"
	"errors"
	"hajime/golangp/apps/hajime_center/constants"
	"hajime/golangp/apps/hajime_center/dify"
	"hajime/golangp/apps/hajime_center/proxy/middleware"
	"hajime/golangp/common/logging"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// Check if the path is /dify/console/api/setup
var excludedPaths = []string{
	"/console/api/setup",
	"/console/api/system-features",
	"/console/api/installed-apps",
	"/console/api/features",
	"/console/api/datasets/retrieval-setting",
}
var excludedPathsPrefix = []string{
	"/api",
	"/console/api/installed-apps",
	"/console/api/apps/",
}

// AuthMiddleware adds authentication headers to the request
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		difyClient, err := dify.GetDifyClient()
		if err != nil {
			logging.Warning("Auth Failed: " + err.Error())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		user, err := middleware.DeserializeUser(r)

		if !middleware.IsPathExcluded(r.URL.Path, excludedPaths, constants.DifyServerPrefix) && !middleware.IsPathPrefix(r.URL.Path, excludedPathsPrefix, constants.DifyServerPrefix) {

			if err != nil {
				logging.Warning("Auth Failed: " + err.Error())
				middleware.WriteErrorResponse(w, "401", err.Error(), http.StatusBadRequest)
				return
			}

			Token, err := difyClient.GetUserToken(user.Role)
			//Token, err := difyClient.GetUserToken("admin")

			if err != nil {
				logging.Warning("Token retrieval failed: " + err.Error())
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			// Add the token to the request header
			r.Header.Set("Authorization", "Bearer "+Token)
		}

		if middleware.IsPathExcluded(r.URL.Path, excludedPaths, constants.DifyServerPrefix) || middleware.IsPathPrefix(r.URL.Path, excludedPathsPrefix, constants.DifyServerPrefix) {
			Token, err := difyClient.GetUserToken("admin")
			if err != nil {
				logging.Warning("Token retrieval failed: " + err.Error())
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			r.Header.Set("Authorization", "Bearer "+Token)
		}

		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

var apiPaths = []string{
	"/console/api/apps/{app_id}/chat-messages",
	"/console/api/installed-apps/{app_id}/chat-messages",
	"/console/api/apps/{app_id}/workflows/draft/run",
	"/console/api/installed-apps/{app_id}/workflows/run",
}

// CreateProxiedServer sets up and starts the HTTP server with middleware
func CreateProxiedServer(wg *sync.WaitGroup) *http.Server {

	mux_router := mux.NewRouter()
	router := mux_router.PathPrefix(constants.DifyServerPrefix).Subrouter()

	// Register handlers with middleware
	router.HandleFunc("/console/api/apps/no_auth", middleware.GetAllNoAuthApp).Methods("GET")
	router.HandleFunc("/console/api/apps/publish/{app_id}", middleware.HandlePublish).Methods("POST")
	router.HandleFunc("/console/api/apps/unpublish/{app_id}", middleware.HandleUnpublish).Methods("POST")
	router.Handle("/console/api/apps/{app_id}", AuthMiddleware(http.HandlerFunc(DifyHandler)))
	router.Handle("/console/api/datasets/{dataset_id}", AuthMiddleware(http.HandlerFunc(DifyHandler)))

	router.Handle("/console/api/apps/{app_id}/model-config", middleware.ModelUpdateMiddleware(http.HandlerFunc(DifyHandler)))

	router.Handle("/console/api/apps/{app_id}/workflows/draft", middleware.WorkflowDraftMiddleware(http.HandlerFunc(DifyHandler)))

	//chat
	for _, path := range apiPaths {
		router.Handle(path, middleware.ChatMessageMiddleware(http.HandlerFunc(DifyHandler)))
	}
	router.Handle("/console/api/apps", AuthMiddleware(http.HandlerFunc(DifyHandler)))

	router.PathPrefix("/").Handler(AuthMiddleware(http.HandlerFunc(DifyHandler)))

	server := &http.Server{
		Addr:    ":8001",
		Handler: router,
	}

	// Start server in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		logging.Info("Starting proxy server on port %s", server.Addr)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logging.Danger("HTTP server error: %v", err)
		}
		logging.Info("Stopped serving new connections.")
	}()

	return server
}
