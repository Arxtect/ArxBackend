package api

import (
	"log"
	"net/http"

	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/controllers"
	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/routes"
	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/service/ws"
	"github.com/Arxtect/ArxBackend/golangp/common/initializers"
	"github.com/Arxtect/ArxBackend/golangp/common/logger"
	"github.com/Arxtect/ArxBackend/golangp/config"
	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/migrate/motest"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var (
	server *gin.Engine
)

var (
	configYml string
	StartCmd  = &cobra.Command{
		Use:          "server",
		Short:        "Start API server",
		Example:      "ArxBackend server -c config/settings-dev.yml",
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings-dev.yml", "Start server with provided configuration file")
}

func setup() {
	//1. 读取配置

	log.Println("🚗 Load configuration file ...")
	err := config.LoadEnv(configYml)
	if err != nil {
		log.Println("🚀 Load failed", err)
		return
	}
	log.Println(`🚗 Loading env is success....`, config.Env.Mode)
	initializers.ConnectDB(&config.Env)
	// initializers.InitRedisClient(&config.Env)
	// initializers.InitMeiliClient(&config.Env)
	if err != nil {
		return
	}
	log.Println("🚗 Connect DB is success....", config.Env.Mode)

}

func run() error {

	server = gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	CreditSystem := controllers.NewCreditSystem(initializers.DB)

	AuthController := controllers.NewAuthController(initializers.DB)
	AuthRouteController := routes.NewAuthRouteController(AuthController)

	UserController := controllers.NewUserController(initializers.DB, CreditSystem)
	UserRouteController := routes.NewRouteUserController(UserController)

	PostController := controllers.NewPostController(initializers.DB)
	PostRouteController := routes.NewRoutePostController(PostController)

	ProjectController := controllers.NewProjectController(initializers.DB)
	ProjectRouteController := routes.NewProjectRouteController(ProjectController)

	ChatController := controllers.NewChatController(CreditSystem)
	ChatRouteController := routes.NewChatRouteController(ChatController)

	DocumentController := controllers.NewDocumentController(initializers.DB, logger.Logger, initializers.Rdb, initializers.MeiliClient)
	DocumentRouteController := routes.NewDocumentRouteController(DocumentController)

	PromptController := controllers.NewPromptController(initializers.DB, logger.Logger, initializers.Rdb)
	PromptRouteController := routes.NewPromptRouteController(PromptController)

	YRedisController := controllers.NewYRedisController(initializers.DB, logger.Logger) // y-redis 回调
	YRedisRouteController := routes.NewYRedisRouteController(YRedisController)

	ConversationsController := controllers.NewConversationsController(initializers.DB) // y-redis 回调
	ConversationsRouteController := routes.NewConversationsRouteController(ConversationsController)

	S3FileController := controllers.NewS3FileController(initializers.DB)
	S3FileRouteController := routes.NewS3FileRouteController(*S3FileController)
	// /api/healthcheck
	router := server.Group("/api/v1")
	router.GET("/healthcheck", func(ctx *gin.Context) {
		message := "Welcome to ChatGPT!"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	PostRouteController.PostRoute(router)
	ProjectRouteController.ProjectkRoute(router)
	ChatRouteController.ChatRoute(router)

	DocumentRouteController.DocumentRoute(router)
	PromptRouteController.PromptRoute(router)
	YRedisRouteController.YRediskRoute(router) // y-redis 支持

	ConversationsRouteController.ConversationsRoute(router) // AI Conversations
	S3FileRouteController.S3FileRoute(router)
	if config.Env.Mode == "test" {
		// 先迁移所有的表
		motest.TestModeMigrate()
		log.Println("🚗 Initialize data creation is success....", config.Env.Mode)
	}

	// 2. 启动weosocket服务
	router.Handle("GET", "/ws", ws.HandlerWs)
	log.Println("🚗 api websocket is starting....", config.Env.Mode)

	router.GET("/collaborative_editing_demo", func(c *gin.Context) {
		ws.HandleGetStaticResource("static/index.html")(c.Writer, c.Request)
	})

	// TODO 转发房间,让其他用户可以订阅，邀请码 , 让其订阅某个房间   订阅和取消订阅, 在长连接里发送
	router.POST("/room/subscribe", func(c *gin.Context) {
		// 加入订阅,需要房间号 , 用户info等
	})

	log.Println("🚗 api server is starting....", config.Env.Mode)
	log.Fatal(server.Run(":" + config.Env.ServerPort))
	return nil
}
