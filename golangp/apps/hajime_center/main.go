package main

import (
	"context"
	"errors"
	"hajime/golangp/apps/hajime_center/controllers"
	"hajime/golangp/apps/hajime_center/initializers"
	"hajime/golangp/apps/hajime_center/proxy"
	"hajime/golangp/apps/hajime_center/routes"
	"hajime/golangp/common/logging"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server *gin.Engine

	CreditSystem *controllers.CreditSystem

	AuthController              controllers.AuthController
	AuthRouteController         routes.AuthRouteController
	ReferralCodeController      controllers.ReferralCodeController
	ReferralCodeRouteController routes.ReferralCodeRouteController

	BalanceHistoryController      controllers.BalanceHistoryController
	BalanceHistoryRouteController routes.BalanceHistoryRouteController

	wg sync.WaitGroup
)

func init() {
	conf, err := initializers.LoadEnv(".")
	if err != nil {
		logging.Danger("🚀 Could not load environment variables: %v", err)
	}

	initializers.ConnectDB(&conf)
	//initializers.ConnectDBDify(&conf)
	CreditSystem = controllers.NewCreditSystem(initializers.DB)

	AuthController = controllers.NewAuthController(initializers.DB, CreditSystem)
	AuthRouteController = routes.NewAuthRouteController(AuthController)
	ReferralCodeController = controllers.NewReferralCodeController(initializers.DB, CreditSystem)
	ReferralCodeRouteController = routes.NewReferralCodeRouteController(ReferralCodeController)

	BalanceHistoryController = controllers.NewBalanceHistoryController(initializers.DB)
	BalanceHistoryRouteController = routes.NewBalanceHistoryRouteController(BalanceHistoryController)

	server = gin.Default()
}

func main() {
	conf, err := initializers.LoadEnv(".")
	if err != nil {
		logging.Danger("🚀 Could not load environment variables: %v", err)
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
		message := "Welcome to HajimeAI!"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	AuthRouteController.AuthRoute(router)
	ReferralCodeRouteController.ReferralCodeRoute(router)
	BalanceHistoryRouteController.BalanceHistoryRoute(router)
	// Start the main server in a new goroutine
	httpServer := &http.Server{
		Addr:    ":" + conf.ServerPort,
		Handler: server,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		logging.Info("Starting server on port %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logging.Danger("HTTP server error:  %v", err)
		}
	}()

	//proxy serve
	proxiedServer := proxy.CreateProxiedServer(&wg)

	// Signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // Wait for a signal

	logging.Info("Shutting down server...")

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logging.Danger("HTTP server shutdown error:  %v", err)
	}

	if err := proxiedServer.Shutdown(shutdownCtx); err != nil {
		logging.Danger("ProxiedServer shutdown error:  %v", err)
	}

	wg.Wait()
	logging.Info("Server exited gracefully.")
}
