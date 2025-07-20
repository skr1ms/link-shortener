package main

import (
	"fmt"
	"linkshortener/configs"
	"linkshortener/internal/auth"
	"linkshortener/internal/link"
	"linkshortener/internal/stats"
	"linkshortener/internal/user"
	"linkshortener/migrations"
	"linkshortener/pkg/event"
	"linkshortener/pkg/jwt"
	"linkshortener/pkg/middleware"
	"net/http"
	"os"
)

func appInit() http.Handler {
	config, err := configs.LoadConfig()
	if err != nil {
		panic(err)
	}

	database := migrations.RunMigrations(config)

	linkRepository := link.NewLinkRepository(database)
	userRepository := user.NewUserRepository(database)
	statsRepository := stats.NewStatsRepository(database)
	eventBus := event.NewEventBus()

	// services
	authService := auth.NewAuthService(userRepository,
		jwt.NewJWT(config.Auth.SecretKey, config.Auth.RefreshTokenSecretKey))

	statsService := stats.NewStatsService(&stats.StatsServiceDeps{
		EventBus:        eventBus,
		StatsRepository: statsRepository,
	})

	router := http.NewServeMux()

	// handlers
	auth.NewAuthHandler(router, &auth.AuthHandlerDeps{
		Config: &configs.AuthConfig{
			SecretKey: os.Getenv("SECRET_KEY"),
		},
		AuthService: authService,
	})

	link.NewLinkHandler(router, &link.LinkHandlerDeps{
		Config:         config,
		LinkRepository: linkRepository,
		EventBus:       eventBus,
	})

	stats.NewStatsHandler(router, &stats.StatsHandlerDeps{
		Config:          config,
		StatsRepository: statsRepository,
	})

	// middlewares
	stack := middleware.Chain(
		middleware.Cors,
		middleware.LogRequest,
	)
	go statsService.AddClick()

	return stack(router)
}

func main() {
	appInit := appInit()

	server := &http.Server{
		Addr:    ":8081",
		Handler: appInit,
	}

	fmt.Println("Server is running on port 8081")
	server.ListenAndServe()
}
