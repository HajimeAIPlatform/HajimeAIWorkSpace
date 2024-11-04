package proxy

import (
	"errors"
	"fmt"
	"hajime/golangp/common/logging"
	"hajime/golangp/hajime_center/dify"
	"log"
	"net/http"
	"sync"
)

// AuthMiddleware adds authentication headers to the request
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the path is /dify/console/api/setup
		fmt.Println("r.URL.Path", r.URL.Path)
		if r.URL.Path != "/dify/console/api/setup" {
			difyClient, err := dify.GetDifyClient()
			if err != nil {
				logging.Warning("Auth Failed: " + err.Error())
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			Token, err := difyClient.GetUserToken()
			if err != nil {
				logging.Warning("Token retrieval failed: " + err.Error())
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			// Add the token to the request header
			r.Header.Set("Authorization", "Bearer "+Token)
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// CreateProxiedServer sets up and starts the HTTP server with middleware
func CreateProxiedServer(wg *sync.WaitGroup) *http.Server {
	mux := http.NewServeMux()

	// Register handlers with middleware
	mux.Handle("/dify/", AuthMiddleware(http.HandlerFunc(DifyHandler)))

	server := &http.Server{
		Addr:    ":8001",
		Handler: mux,
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

// DifyHandler forwards requests after removing the "dify" prefix
