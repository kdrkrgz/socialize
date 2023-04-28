package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kdrkrgz/socalize/handler"
	log "github.com/kdrkrgz/socalize/pkg/logger"
	"github.com/kdrkrgz/socalize/pkg/seed"
	"github.com/kdrkrgz/socalize/repository"
	"os"
	"os/signal"
	"syscall"
)

type Application struct {
	app  *fiber.App
	repo *repository.Repository
}

func (a *Application) Register() {
	a.app.Get("/users", handler.GetUsers(a.repo))
	//route := a.app.Group("/swagger")
	//route.Get("*", swagger.HandlerDefault)
}

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

	if err := app.Listen(":8000"); err != nil {
		log.Logger().Panic(fmt.Sprintf("App Err: %s", err))
	}

}
