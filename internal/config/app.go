package config

import (
	"slido-clone-backend/internal/delivery/http"
	"slido-clone-backend/internal/delivery/http/middleware"
	"slido-clone-backend/internal/delivery/http/route"
	"slido-clone-backend/internal/delivery/websocket"
	"slido-clone-backend/internal/repository"
	"slido-clone-backend/internal/sfu"
	"slido-clone-backend/internal/usecase"
	"slido-clone-backend/internal/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB        *gorm.DB
	App       *fiber.App
	Redis     *redis.Client
	Log       *logrus.Logger
	Validator *validator.Validate
	Config    *viper.Viper
	WSHub     *websocket.Hub
}

func Bootstrap(config *BootstrapConfig) {
	// setup repositories
	userRepository := repository.NewUserRepository(config.Log)
	roomRepository := repository.NewRoomRepository(config.Log)
	participantRepository := repository.NewParticipantRepository(config.Log)
	messageRepository := repository.NewMessageRepository(config.Log)
	xpTransactionRepository := repository.NewXPTransactionRepository(config.Log)
	questionRepository := repository.NewQuestionRepository(config.Log)
	voteRepository := repository.NewVoteRepository(config.Log)
	pollRepository := repository.NewPollRepository(config.Log)

	// setup utils
	tokenUtil := util.NewTokenUtil(config.Config.GetString("JWT_SECRET"), config.Redis)

	// setup use cases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validator, userRepository, participantRepository, roomRepository, tokenUtil)
	roomUseCase := usecase.NewRoomUseCase(config.DB, config.Log, config.Validator, roomRepository, participantRepository)
	participantUseCase := usecase.NewParticipantUseCase(config.DB, config.Log, config.Validator, participantRepository, roomRepository, userRepository, tokenUtil)
	xpTransactionUseCase := usecase.NewXPTransactionUseCase(config.DB, config.Validator, config.Log, xpTransactionRepository)
	messageUseCase := usecase.NewMessageUseCase(config.DB, config.Validator, config.Log, messageRepository, roomRepository, participantRepository, xpTransactionUseCase)
	questionUseCase := usecase.NewQuestionUseCase(config.DB, config.Log, config.Validator, questionRepository, voteRepository, roomRepository, participantRepository, xpTransactionRepository)
	pollUseCase := usecase.NewPollUseCase(config.DB, config.Log, config.Validator, pollRepository, roomRepository, participantRepository, xpTransactionRepository)

	// configuration websocket hub (sebelum controller yang membutuhkan hub)
	hub := websocket.NewHub(config.Log)
	go hub.Run() // start hub run goroutine

	// setup HTTP controllers
	userController := http.NewUserController(config.Log, userUseCase)
	roomController := http.NewRoomController(config.Log, roomUseCase, tokenUtil)
	participantController := http.NewParticipantController(config.Log, participantUseCase)
	messageController := http.NewMessageController(config.Log, messageUseCase)
	questionController := http.NewQuestionController(config.Log, questionUseCase, hub)
	pollController := http.NewPollController(config.Log, pollUseCase, participantUseCase, hub)

	// setup HTTP middleware
	authMiddleware := middleware.NewAuth(userUseCase, tokenUtil)

	// websocket handler
	sfuManager := sfu.NewSFUManager(config.Log)
	eventHandler := websocket.NewEventHandler(messageUseCase, participantUseCase, questionUseCase, sfuManager)
	wsHandler := websocket.NewWebSocketHandler(hub, config.Log, tokenUtil, eventHandler)

	// setup HTTP routes
	routeConfig := route.RouteConfig{
		App:                   config.App,
		UserController:        userController,
		RoomController:        roomController,
		ParticipantController: participantController,
		MessageController:     messageController,
		QuestionController:    questionController,
		PollController:        pollController,
		AuthMiddleware:        authMiddleware,
		WSHandler:             wsHandler,
	}
	routeConfig.Setup()
}
