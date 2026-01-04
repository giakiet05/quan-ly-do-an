package bootstrap

import (
	"log"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/auth"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/config"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/controller"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/platform/bus"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/platform/email"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/platform/ws"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/repo"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/route"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repos struct {
	repo.UserRepo
	repo.ClassroomRepo
	repo.GroupRepo
	repo.NotificationRepo
	repo.ChannelRepo
	repo.MessageRepo
	repo.EmailVerificationRepo
	repo.PasswordResetRepo
	repo.ClassroomInvitationRepo
}

type Services struct {
	service.AuthService
	service.UserService
	service.GroupService
	service.NotificationService
	service.ChannelService
	service.MessageService
	service.ClassroomInvitationService
}

type Controllers struct {
	controller.AuthController
	controller.UserController
	controller.NotificationController
	controller.WebSocketController
	controller.ChannelController
	controller.MessageController
	controller.GroupController
	controller.ClassroomInvitationController
}

func initRepos(client *mongo.Client, db *mongo.Database) *Repos {
	return &Repos{
		UserRepo:                repo.NewUserRepo(db),
		NotificationRepo:        repo.NewNotificationRepo(db),
		ChannelRepo:             repo.NewChannelRepo(db),
		MessageRepo:             repo.NewMessageRepo(db),
		GroupRepo:               repo.NewGroupRepo(db),
		EmailVerificationRepo:   repo.NewEmailVerificationRepo(db),
		PasswordResetRepo:       repo.NewPasswordResetRepo(db),
		ClassroomRepo:           repo.NewClassroomRepo(db),
		ClassroomInvitationRepo: repo.NewClassroomInvitationRepo(db),
	}
}

func initServices(repos *Repos, redisClient *redis.Client, emailSender email.Sender, eventBus *bus.EventBus, tokenService *auth.TokenService) *Services {
	return &Services{
		GroupService:               service.NewGroupService(repos.GroupRepo, repos.ClassroomRepo, repos.ChannelRepo, repos.UserRepo, eventBus),
		AuthService:                service.NewAuthService(repos.UserRepo, repos.EmailVerificationRepo, repos.PasswordResetRepo, emailSender, redisClient, tokenService),
		UserService:                service.NewUserService(repos.UserRepo, eventBus, redisClient),
		NotificationService:        service.NewNotificationService(repos.NotificationRepo, repos.UserRepo, eventBus, redisClient),
		ChannelService:             service.NewChannelService(repos.ChannelRepo, repos.MessageRepo, eventBus),
		MessageService:             service.NewMessageService(repos.MessageRepo, repos.ChannelRepo, eventBus, redisClient),
		ClassroomInvitationService: service.NewClassroomInvitationService(repos.ClassroomInvitationRepo, repos.ClassroomRepo, repos.UserRepo, emailSender),
	}
}

func initControllers(services *Services, wsHub *ws.Hub) *Controllers {
	return &Controllers{
		GroupController:               *controller.NewGroupController(services.GroupService),
		AuthController:                *controller.NewAuthController(services.AuthService),
		UserController:                *controller.NewUserController(services.UserService),
		NotificationController:        *controller.NewNotificationController(services.NotificationService),
		WebSocketController:           *controller.NewWebSocketController(wsHub),
		ChannelController:             *controller.NewChannelController(services.ChannelService),
		MessageController:             *controller.NewMessageController(services.MessageService),
		ClassroomInvitationController: *controller.NewClassroomInvitationController(services.ClassroomInvitationService),
	}
}

func initRoutes(controllers *Controllers, r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	api := r.Group("/api")
	api.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to LKForum API!"})
	})

	route.RegisterAuthRoutes(api, &controllers.AuthController, &controllers.UserController)
	route.RegisterUserRoutes(api, &controllers.UserController)
	route.RegisterNotificationRoutes(api, &controllers.NotificationController)
	route.RegisterWebSocketRoutes(api, &controllers.WebSocketController)
	route.RegisterChannelRoutes(api, &controllers.ChannelController)
	route.RegisterMessageRoutes(api, &controllers.MessageController)
	route.SetupClassroomInvitationRoutes(r, &controllers.ClassroomInvitationController)
}

func Init() (*gin.Engine, error) {
	config.LoadConfig()
	auth.InitGoogleOAuthConfig()

	redisClient := config.NewRedisClient()

	tokenService, err := InitializeTokenService(redisClient)
	if err != nil {
		log.Printf("Warning: Token invalidation service not available: %v\n", err)
	}

	mongoClient := config.NewMongoClient()
	db := mongoClient.Database(config.Cfg.DBName)
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", config.Cfg.FrontendURL)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	eventBus := bus.NewEventBus()
	wsHub := ws.NewHub(eventBus)
	emailSender := email.NewSMTPSender()

	repos := initRepos(mongoClient, db)
	services := initServices(repos, redisClient, emailSender, eventBus, tokenService)
	controllers := initControllers(services, wsHub)
	initRoutes(controllers, router)

	// Start background services
	go wsHub.Start()
	services.NotificationService.Start()
	services.MessageService.Start()
	services.ChannelService.Start()

	return router, nil
}
