package config

import (
	"slido-clone-backend/internal/delivery/http"
	"slido-clone-backend/internal/delivery/http/middleware"
	"slido-clone-backend/internal/delivery/http/route"
	"slido-clone-backend/internal/delivery/websocket"
	"slido-clone-backend/internal/repository"
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

	// setup utils
	tokenUtil := util.NewTokenUtil(config.Config.GetString("JWT_SECRET"), config.Redis)

	// setup use cases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validator, userRepository, participantRepository, roomRepository, tokenUtil)
	roomUseCase := usecase.NewRoomUseCase(config.DB, config.Log, config.Validator, roomRepository)
	participantUseCase := usecase.NewParticipantUseCase(config.DB, config.Log, config.Validator, participantRepository, roomRepository, userRepository, tokenUtil)
	messageUseCase := usecase.NewMessageUseCase(config.DB, config.Validator, config.Log, messageRepository, roomRepository, participantRepository)

	// setup HTTP controllers
	userController := http.NewUserController(config.Log, userUseCase)
	roomController := http.NewRoomController(config.Log, roomUseCase)
	participantController := http.NewParticipantController(config.Log, participantUseCase)
	messageController := http.NewMessageController(config.Log, messageUseCase)

	// setup HTTP middleware
	authMiddleware := middleware.NewAuth(userUseCase, tokenUtil)

	// configuration websocket
	hub := websocket.NewHub(config.Log)
	go hub.Run() // start hub run goroutine

	// websocket handler
	eventHandler := websocket.NewEventHandler(messageUseCase)
	wsHandler := websocket.NewWebSocketHandler(hub, config.Log, tokenUtil, eventHandler)

	// setup HTTP routes
	routeConfig := route.RouteConfig{
		App:                   config.App,
		UserController:        userController,
		RoomController:        roomController,
		ParticipantController: participantController,
		MessageController:     messageController,
		AuthMiddleware:        authMiddleware,
		WSHandler:             wsHandler,
	}
	routeConfig.Setup()
}
