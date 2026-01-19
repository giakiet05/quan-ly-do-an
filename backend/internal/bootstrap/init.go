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
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repos struct {
	repo.UserRepo
	repo.ClassroomRepo
	repo.ClassPostRepo
	repo.GroupRepo
	repo.NotificationRepo
	repo.ChannelRepo
	repo.MessageRepo
	repo.EmailVerificationRepo
	repo.PasswordResetRepo
	repo.ClassroomJoinRequestRepo
	repo.ClassroomInvitationRepo
	repo.ProjectRepo
}

type Services struct {
	service.AuthService
	service.UserService
	service.GroupService
	service.NotificationService
	service.ChannelService
	service.MessageService
	service.ClassroomJoinService
	service.ClassroomService
	service.ClassroomInvitationService
	service.ClassPostService
	service.ProjectService
}

type Controllers struct {
	controller.AuthController
	controller.UserController
	controller.NotificationController
	controller.WebSocketController
	controller.ChannelController
	controller.MessageController
	controller.GroupController
	controller.ClassroomJoinController
	controller.ClassroomController
	controller.ClassroomInvitationController
	controller.ClassPostController
	controller.ProjectController
}

func initRepos(client *mongo.Client, db *mongo.Database) *Repos {
	return &Repos{
		UserRepo:                 repo.NewUserRepo(db),
		NotificationRepo:         repo.NewNotificationRepo(db),
		ChannelRepo:              repo.NewChannelRepo(db),
		MessageRepo:              repo.NewMessageRepo(db),
		GroupRepo:                repo.NewGroupRepo(db),
		EmailVerificationRepo:    repo.NewEmailVerificationRepo(db),
		PasswordResetRepo:        repo.NewPasswordResetRepo(db),
		ClassroomRepo:            repo.NewClassroomRepo(db),
		ClassroomJoinRequestRepo: repo.NewClassroomJoinRequestRepo(db),
		ClassroomInvitationRepo:  repo.NewClassroomInvitationRepo(db),
		ClassPostRepo:            repo.NewClassPostRepo(db),
		ProjectRepo:              repo.NewProjectRepo(db),
	}
}

func initServices(
	repos *Repos,
	redisClient *redis.Client,
	emailSender email.Sender,
	eventBus *bus.EventBus,
	cron *cron.Cron,
	tokenService *auth.TokenService,
) *Services {
	return &Services{
		GroupService:                service.NewGroupService(repos.GroupRepo, repos.ClassroomRepo, repos.ChannelRepo, repos.UserRepo, repos.ProjectRepo),
		AuthService:                 service.NewAuthService(repos.UserRepo, repos.EmailVerificationRepo, repos.PasswordResetRepo, emailSender, redisClient, tokenService),
		UserService:                 service.NewUserService(repos.UserRepo, eventBus, redisClient),
		NotificationService:         service.NewNotificationService(repos.NotificationRepo, repos.UserRepo, eventBus, redisClient),
		ChannelService:              service.NewChannelService(repos.ChannelRepo, repos.MessageRepo, eventBus),
		MessageService:              service.NewMessageService(repos.MessageRepo, repos.ChannelRepo, eventBus, redisClient),
		ClassroomJoinService:        service.NewClassroomJoinService(repos.ClassroomJoinRequestRepo, repos.ClassroomRepo, repos.UserRepo),
		ClassroomService:            service.NewClassroomService(repos.ClassroomRepo, repos.UserRepo, repos.ChannelRepo),
		ClassroomInvitationService:  service.NewClassroomInvitationService(repos.ClassroomInvitationRepo, repos.ClassroomRepo, repos.UserRepo, eventBus),
		ClassPostService:            service.NewClassPostService(repos.ClassPostRepo, repos.ClassroomRepo, repos.UserRepo),
		ProjectService:              service.NewProjectService(repos.ProjectRepo, repos.ClassroomRepo),
	}
}

func initControllers(services *Services, wsHub *ws.Hub) *Controllers {
	return &Controllers{
		GroupController:                *controller.NewGroupController(services.GroupService),
		AuthController:                 *controller.NewAuthController(services.AuthService),
		UserController:                 *controller.NewUserController(services.UserService),
		NotificationController:         *controller.NewNotificationController(services.NotificationService),
		WebSocketController:            *controller.NewWebSocketController(wsHub),
		ChannelController:              *controller.NewChannelController(services.ChannelService),
		MessageController:              *controller.NewMessageController(services.MessageService),
		ClassroomJoinController:        *controller.NewClassroomJoinController(services.ClassroomJoinService),
		ClassroomController:            *controller.NewClassroomController(services.ClassroomService),
		ClassroomInvitationController:  *controller.NewClassroomInvitationController(services.ClassroomInvitationService),
		ClassPostController:            *controller.NewClassPostController(services.ClassPostService),
		ProjectController:              *controller.NewProjectController(services.ProjectService),
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
	route.RegisterClassroomRoutes(api, &controllers.ClassroomController, &controllers.ClassroomInvitationController)
	route.RegisterClassroomJoinRoutes(api, &controllers.ClassroomJoinController)
	route.RegisterClassPostRoutes(api, &controllers.ClassPostController)
	route.RegisterGroupRoutes(api, &controllers.GroupController)
	route.RegisterProjectRoutes(api, &controllers.ProjectController)
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
	cronService := cron.New()

	repos := initRepos(mongoClient, db)
	services := initServices(repos, redisClient, emailSender, eventBus, cronService, tokenService)
	controllers := initControllers(services, wsHub)
	initRoutes(controllers, router)

	// Start background services
	go wsHub.Start()
	services.NotificationService.Start()
	services.MessageService.Start()
	services.ChannelService.Start()
	services.GroupService.Start()
	services.ProjectService.Start()
	cronService.Start()

	return router, nil
}
