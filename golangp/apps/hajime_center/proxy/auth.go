package proxy

import (
	"errors"
	"fmt"
	"hajime/golangp/apps/hajime_center/dify"
	"hajime/golangp/apps/hajime_center/initializers"
	"hajime/golangp/apps/hajime_center/models"
	"hajime/golangp/common/logging"
	"hajime/golangp/common/utils"
	"log"
	"net/http"
	"strings"
	"sync"
)

func DeserializeUser(r *http.Request) (user *models.User, err error) {
	authorizationHeader := r.Header.Get("Authorization")
	fields := strings.Fields(authorizationHeader)

	accessToken := ""

	if len(fields) != 0 && fields[0] == "Bearer" {
		accessToken = fields[1]
	}

	if accessToken == "" {
		return nil, errors.New("you are not logged in")
	}

	config, _ := initializers.LoadEnv(".")
	sub, err := utils.ValidateToken(accessToken, config.AccessTokenPublicKey)
	if err != nil {
		return nil, err
	}

	result := initializers.DB.First(&user, "id = ?", fmt.Sprint(sub))
	if result.Error != nil {
		return nil, errors.New("the user belonging to this token no logger exists")
	}
	return user, nil
}

// AuthMiddleware adds authentication headers to the request
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the path is /dify/console/api/setup
		if r.URL.Path != "/dify/console/api/setup" {
			user, _ := DeserializeUser(r)

			difyClient, err := dify.GetDifyClient()
			if err != nil {
				logging.Warning("Auth Failed: " + err.Error())
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			Token, err := difyClient.GetUserToken(user.Role)
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
