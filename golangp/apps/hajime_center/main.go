package main

import (
	"context"
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"hajime/golangp/apps/hajime_center/controllers"
	"hajime/golangp/apps/hajime_center/initializers"
	"hajime/golangp/apps/hajime_center/proxy"
	"hajime/golangp/apps/hajime_center/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	AppsController      controllers.AppsController
	AppsRouteController routes.AppsRouteController
	wg                  sync.WaitGroup
)

func init() {
	conf, err := initializers.LoadEnv(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	initializers.ConnectDB(&conf)
	//initializers.ConnectDBDify(&conf)

	AuthController = controllers.NewAuthController(initializers.DB)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	AppsController = controllers.NewAppsController(initializers.DB)
	AppsRouteController = routes.NewAppsRouteController(AppsController)

	server = gin.Default()
}

func main() {
	conf, err := initializers.LoadEnv(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	// 将 DOMAIN 值解析为切片
	allowedOrigins := strings.Split(conf.ClientOrigin, ",")

	// 如果 allowedOrigins 只有一个空字符串元素，设置为允许所有来源
	if len(allowedOrigins) == 1 && allowedOrigins[0] == "" {
		allowedOrigins = []string{"*"}
	}

	// 使用 gin-contrib/cors 中间件
	corsConfig := cors.Config{
		AllowOrigins: allowedOrigins, // 允许所有来源
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Description-Type",
			"Accept",
			"Authorization",
			"App-Authorization",
			"Access-Control-Allow-Origin",
			"app-authorization",
		},
		ExposeHeaders:    []string{"Description-Length"},
		AllowCredentials: true,
	}

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthcheck", func(ctx *gin.Context) {
		message := "Welcome to ChatGPT!"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	AuthRouteController.AuthRoute(router)
	AppsRouteController.AppsRoute(router)

	// Start the main server in a new goroutine
	httpServer := &http.Server{
		Addr:    ":" + conf.ServerPort,
		Handler: server,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		log.Printf("Starting server on port %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	//proxy serve
	proxiedServer := proxy.CreateProxiedServer(&wg)

	// Signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // Wait for a signal

	log.Println("Shutting down server...")

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP server shutdown error: %v", err)
	}

	if err := proxiedServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("ProxiedServer shutdown error: %v", err)
	}

	wg.Wait()
	log.Println("Server exited gracefully.")
}
