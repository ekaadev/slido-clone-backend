package main

import (
	"fmt"
	"slido-clone-backend/internal/config"
)

func main() {
	// initialize configurations and dependencies
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	validate := config.NewValidator(viperConfig)
	redis := config.NewRedisClient(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	app := config.NewFiber(viperConfig)

	// bootstrap application, assign configurations and dependencies
	config.Bootstrap(&config.BootstrapConfig{
		DB:        db,
		App:       app,
		Redis:     redis,
		Log:       log,
		Validator: validate,
		Config:    viperConfig,
	})

	// start web server
	webPort := viperConfig.GetInt("web.port")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Error starting web server: %s", err)
	}
}
