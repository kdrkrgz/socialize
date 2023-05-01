package main

import (
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/kdrkrgz/socalize/conf"
	"github.com/kdrkrgz/socalize/handler"
	log "github.com/kdrkrgz/socalize/pkg/logger"
	"github.com/kdrkrgz/socalize/pkg/seed"
	"github.com/kdrkrgz/socalize/repository"
	_ "github.com/kdrkrgz/socalize/swagger"
	"os"
	"os/signal"
	"syscall"
)

type Application struct {
	app  *fiber.App
	repo *repository.Repository
}

func (a *Application) Register() {
	a.app.Get("/", handler.RedirectSwagger)
	a.app.Get("/users", handler.GetUsers(a.repo))
	route := a.app.Group("/swagger")
	route.Get("*", swagger.HandlerDefault)
}

// @title						Socalize API
// @version					    1.0
// @description				    Swagger for Socalize app
// @host						http://localhost:8000
// @BasePath					/
// @schemes					    https http
// @license.name				Apache License, Version 2.0 (the "License")
// @license.url				    https://github.com/acikkaynak/deprem-yardim-backend-go/blob/main/LICENSE
// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						X-Api-Key
func main() {
	seed.InitialDataSeed()
	repo := repository.New()
	defer repo.Close()
	app := fiber.New()
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
