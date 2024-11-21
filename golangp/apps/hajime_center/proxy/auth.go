package proxy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"hajime/golangp/apps/hajime_center/dify"
	"hajime/golangp/apps/hajime_center/proxy/middleware"
	"hajime/golangp/common/logging"
	"log"
	"net/http"
	"strings"
	"sync"
)

func writeErrorResponse(w http.ResponseWriter, code, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := map[string]interface{}{
		"code":    code,
		"message": message,
		"status":  status,
	}
	json.NewEncoder(w).Encode(response)
}

func isPathExcluded(path string, excludedPaths []string) bool {
	for _, excludedPath := range excludedPaths {
		if path == excludedPath {
			return true
		}
	}
	return false
}

func isPathPrefix(path string, excludedPaths []string) bool {
	for _, excludedPath := range excludedPaths {
		if strings.HasPrefix(path, excludedPath) {
			return true
		}
	}
	return false
}

// Check if the path is /dify/console/api/setup
var excludedPaths = []string{
	"/dify/console/api/setup",
	"/dify/console/api/system-features",
	"/dify/console/api/installed-apps",
	"/dify/console/api/features",
	"/dify/console/api/datasets/retrieval-setting",
}
var excludedPathsPrefix = []string{
	"/dify/api",
	"/dify/console/api/installed-apps",
	"/dify/console/api/apps/",
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

		if !isPathExcluded(r.URL.Path, excludedPaths) && !isPathPrefix(r.URL.Path, excludedPathsPrefix) {
			user, err := middleware.DeserializeUser(r)
			if err != nil {
				logging.Warning("Auth Failed: " + err.Error())
				writeErrorResponse(w, "401", err.Error(), http.StatusBadRequest)
				return
			}

			Token, err := difyClient.GetUserToken(user.Role)
			//Token, err := difyClient.GetUserToken("admin")
			fmt.Println(Token)

			if err != nil {
				logging.Warning("Token retrieval failed: " + err.Error())
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			// Add the token to the request header
			r.Header.Set("Authorization", "Bearer "+Token)

			ctx := context.WithValue(r.Context(), "user", user)
			r = r.WithContext(ctx)
		}

		if isPathExcluded(r.URL.Path, excludedPaths) || isPathPrefix(r.URL.Path, excludedPathsPrefix) {
			Token, err := difyClient.GetUserToken("admin")
			if err != nil {
				logging.Warning("Token retrieval failed: " + err.Error())
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			r.Header.Set("Authorization", "Bearer "+Token)
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// CreateProxiedServer sets up and starts the HTTP server with middleware
func CreateProxiedServer(wg *sync.WaitGroup) *http.Server {

	router := mux.NewRouter()

	// Register handlers with middleware
	router.HandleFunc("/dify/console/api/apps/no_auth", middleware.GetAllNoAuthApp).Methods("GET")
	router.HandleFunc("/dify/console/api/apps/publish/{app_id}", middleware.HandlePublish).Methods("POST")
	router.HandleFunc("/dify/console/api/apps/unpublish/{app_id}", middleware.HandleUnpublish).Methods("POST")
	router.Handle("/dify/console/api/apps/{app_id}", AuthMiddleware(http.HandlerFunc(DifyHandler)))
	router.Handle("/dify/console/api/datasets/{dataset_id}", AuthMiddleware(http.HandlerFunc(DifyHandler)))

	router.Handle("/dify/console/api/apps", AuthMiddleware(http.HandlerFunc(DifyHandler)))
	router.PathPrefix("/dify/").Handler(AuthMiddleware(http.HandlerFunc(DifyHandler)))

	server := &http.Server{
		Addr:    ":8001",
		Handler: router,
	}

	// Start server in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("Starting proxy server on port %s", server.Addr)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	return server
}
