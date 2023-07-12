package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/kdrkrgz/socalize/conf"
	"github.com/kdrkrgz/socalize/handler"
	"github.com/kdrkrgz/socalize/middleware"
	log "github.com/kdrkrgz/socalize/pkg/logger"
	"github.com/kdrkrgz/socalize/pkg/seed"
	"github.com/kdrkrgz/socalize/repository"
	_ "github.com/kdrkrgz/socalize/swagger"
)

type Application struct {
	app  *fiber.App
	repo *repository.Repository
}

func (a *Application) Register() {
	a.app.Get("/", handler.RedirectSwagger)
	a.app.Get("/healthcheck", handler.HealthCheck)
	a.app.Get("/readiness", handler.Readiness)
	a.app.Get("/users", middleware.DeserializeUser, handler.GetUsers(a.repo))
	route := a.app.Group("/swagger")
	authRoute := a.app.Group("/auth")
	authRoute.Post("/signup", handler.SignUp(a.repo))
	authRoute.Post("/signin", handler.SignIn(a.repo))
	route.Get("*", swagger.HandlerDefault)
}

// @title						Socalize API
// @version					    1.0
// @description				    Swagger for Socalize app
// @host						localhost:8000
// @BasePath					/
// @schemes					    http
// @license.name				Apache License, Version 2.0 (the "License")
// @license.url				    https://github.com/acikkaynak/deprem-yardim-backend-go/blob/main/LICENSE
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	seed.InitialDataSeed()
	repo := repository.New()
	defer repo.Close()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS",
		AllowCredentials: true,
	}))
	application := &Application{app: app, repo: repo}
	application.Register()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		_ = <-c
		log.Logger().Info("application gracefully shutting down..")
		_ = app.Shutdown()
	}()

	if err := app.Listen(fmt.Sprintf(":%v", conf.Get("AppPort"))); err != nil {
		log.Logger().Panic(fmt.Sprintf("App Err: %s", err))
	}

}
