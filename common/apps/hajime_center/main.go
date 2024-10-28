package main

import (
	"HajimeAIWorkSpace/common/apps/hajime_center/chat-config"
	"HajimeAIWorkSpace/common/apps/hajime_center/controllers"
	"HajimeAIWorkSpace/common/apps/hajime_center/initializers"
	"HajimeAIWorkSpace/common/apps/hajime_center/routes"
	"github.com/alecthomas/kong"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	PostController      controllers.PostController
	PostRouteController routes.PostRouteController

	ChatController      controllers.ChatController
	ChatRouteController routes.ChatRouteController

	DocumentController      controllers.DocumentController
	DocumentRouteController routes.DocumentRouteController

	AppsController      controllers.AppsController
	AppsRouteController routes.AppsRouteController

	ModelController      controllers.ModelController
	ModelRouteController routes.ModelRouteController

	ConversationsController      controllers.ConversationsController
	ConversationsRouteController routes.ConversationsRouteController


	CreditSystem *controllers.CreditSystem

	Channel chan struct{} // 定义 sem 变量
)

func init() {
	conf, err := initializers.LoadEnv(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	initializers.ConnectDB(&conf)
	Channel = make(chan struct{}, conf.ThreadNumber)
	//initializers.ConnectDBDify(&conf)

	CreditSystem = controllers.NewCreditSystem(initializers.DB)

	AuthController = controllers.NewAuthController(initializers.DB, CreditSystem)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	PostController = controllers.NewPostController(initializers.DB)
	PostRouteController = routes.NewRoutePostController(PostController)

	ChatController = controllers.NewChatController(CreditSystem)
	ChatRouteController = routes.NewChatRouteController(ChatController)

	DocumentController = controllers.NewDocumentController(initializers.DB, initializers.DBDify)
	DocumentRouteController = routes.NewDocumentRouteController(DocumentController)

	AppsController = controllers.NewAppsController(initializers.DB)
	AppsRouteController = routes.NewAppsRouteController(AppsController)

	ModelController = controllers.NewModelController(initializers.DB)
	ModelRouteController = routes.NewModelRouteController(ModelController)


	ConversationsController = controllers.NewConversationsController(initializers.DB)
	ConversationsRouteController = routes.NewConversationsRouteController(ConversationsController)


	server = gin.Default()
}

func main() {
	kong.Parse(&chat_config.CLI)

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
			"Content-Type",
			"Accept",
			"Authorization",
			"App-Authorization",
			"Access-Control-Allow-Origin",
			"app-authorization",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	server.Use(cors.New(corsConfig))

	router := server.Group("/api/v1")
	routerV2 := server.Group("/api")
	router.GET("/healthcheck", func(ctx *gin.Context) {
		message := "Welcome to ChatGPT!"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	routerV2.GET("/healthcheck", func(ctx *gin.Context) {
		message := "Welcome to ChatGPT!"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	AuthRouteController.AuthRoute(router)
	PostRouteController.PostRoute(router)
	ChatRouteController.ChatRoute(router)
	DocumentRouteController.DocumentRoute(router)
	AppsRouteController.AppsRoute(router)
	ModelRouteController.ModelRoute(router)
	ConversationsRouteController.ConversationsRoute(router)


	log.Fatal(server.Run(":" + conf.ServerPort))
}
