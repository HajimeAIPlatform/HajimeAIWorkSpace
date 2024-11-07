package proxy

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

// SetupCORS configures and returns a CORS handler
func SetupCORS(handler http.Handler) http.Handler {
	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Replace with specific origins in production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	})

	return corsConfig.Handler(handler)
}

func DeserializeUser(r *http.Request) (user *models.User, err error) {
	authorizationHeader := r.Header.Get("Authorization")
	fields := strings.Fields(authorizationHeader)

	accessToken := ""

	if len(fields) != 0 && fields[0] == "Bearer" {
		accessToken = fields[1]
	}

	if accessToken == "" {
		return user, errors.New("you are not logged in")
	}

	config, _ := initializers.LoadEnv(".")
	sub, err := utils.ValidateToken(accessToken, config.AccessTokenPublicKey)
	if err != nil {
		return user, err
	}

	result := initializers.DB.First(&user, "id = ?", fmt.Sprint(sub))
	if result.Error != nil {
		err = errors.New("the user belonging to this token no longer exists")
		return
	}
	return user, nil
}

// AuthMiddleware adds authentication headers to the request
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the path is /dify/console/api/setup
		if r.URL.Path != "/dify/console/api/setup" {
			user, err := DeserializeUser(r)

			fmt.Println("user.Role")

			if err != nil {
				logging.Warning("Auth Failed: " + err.Error())
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			difyClient, err := dify.GetDifyClient()
			if err != nil {
				logging.Warning("Auth Failed: " + err.Error())
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// CreateProxiedServer sets up and starts the HTTP server with middleware
func CreateProxiedServer(wg *sync.WaitGroup) *http.Server {

	router := mux.NewRouter()

	// Register handlers with middleware
	router.Handle("/dify/console/api/apps", AuthMiddleware(http.HandlerFunc(DifyHandler)))
	router.Handle("/dify/console/api/apps/{app_id}", AuthMiddleware(http.HandlerFunc(DifyHandler)))
	router.HandleFunc("/dify/console/api/apps/publish/{app_id}", HandlePublish).Methods("POST")

	router.Handle("/dify/", AuthMiddleware(http.HandlerFunc(DifyHandler)))

	corsHandler := SetupCORS(router)

	server := &http.Server{
		Addr:    ":8001",
		Handler: corsHandler,
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
