package proxy

import (
	"errors"
	"log"
	"net/http"
	"sync"
)

// AuthMiddleware adds authentication headers to the request
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add authentication header
		r.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYmU2YmVhNzctOGIwYS00NGU5LTliY2QtZjFmYzI4OGVmNDUzIiwiZXhwIjoxNzMyMjYzMTU0LCJpc3MiOiJTRUxGX0hPU1RFRCIsInN1YiI6IkNvbnNvbGUgQVBJIFBhc3Nwb3J0In0.M8dW_vVL1DbHQRu01cPHR-Gok5bodnDK_l0uQRPzJQg")

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
