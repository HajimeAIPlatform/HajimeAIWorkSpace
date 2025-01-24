package service

import (
	"context"
	"hajime/golangp/apps/AgentTown/agent"
	"hajime/golangp/apps/AgentTown/config"
	"hajime/golangp/apps/AgentTown/runtime"
	"hajime/golangp/apps/AgentTown/task"
	"hajime/golangp/apps/AgentTown/telemetry"
	"hajime/golangp/common/logging"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	wg sync.WaitGroup
)

func setupRouter() *gin.Engine {
	router := gin.Default()

	// Setup CORS
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	router.Use(cors.New(corsConfig))

	// Define routes
	router.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Service is running"})
	})

	router.POST("/assign-task/:agentID", func(ctx *gin.Context) {
		agentID := ctx.Param("agentID")
		var req struct {
			TaskName string `json:"task_name"`
		}
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		runtime.AssignTaskByAgentID(agentID, task.NewTask(req.TaskName))
		ctx.JSON(http.StatusOK, gin.H{"status": "task assigned"})
	})

	return router
}

func StartService() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.NewConfig("Config_Empty")

	// Create all agents using the same configuration
	agentA := agent.NewAgent(cfg)
	agentB := agent.NewAgent(cfg)
	agentC := agent.NewAgent(cfg)

	// Create and add agents
	runtime.AddAgent(agentA)
	runtime.AddAgent(agentB)
	runtime.AddAgent(agentC)

	// Start agents
	go runtime.StartAgents(ctx)

	// Start telemetry monitoring
	go telemetry.Monitor(5 * time.Second)

	// Setup Gin server
	router := setupRouter()

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		logging.Info("Starting server on port %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Danger("HTTP server error: %v", err)
		}
	}()

	// Signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // Wait for a signal

	logging.Info("Shutting down server...")

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logging.Danger("HTTP server shutdown error: %v", err)
	}

	wg.Wait()
	logging.Info("Server exited gracefully.")
}
