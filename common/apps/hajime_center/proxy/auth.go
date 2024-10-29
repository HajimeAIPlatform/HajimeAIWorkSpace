package proxy

import (
	"HajimeAIWorkSpace/common/apps/hajime_center/dify"
	"HajimeAIWorkSpace/common/apps/hajime_center/logger"
	"errors"
	"log"
	"net/http"
	"sync"
)

// AuthMiddleware adds authentication headers to the request
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add authentication header
		difyCleint, err := dify.GetDifyClient()

		if err != nil {
			logger.Warning("Auth Failed: " + err.Error())
		}

		Token, _ := difyCleint.GetUserToken()
		r.Header.Add("Authorization", "Bearer "+Token)

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
