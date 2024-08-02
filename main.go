package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

    "github.com/metallust/rms-be/internals/bootstrap/web"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	// Common
	httpServer *fiber.App
	client    *mongo.Client
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() {
	waits := make(chan int)

	a.setupEnv()
	a.setupDatabases()
	a.setupHTTP()
	go a.shutdown()

	port := ":" + os.Getenv("WEB_PORT")
	if port == ":" {
		port = ":3000"
	}

	if err := a.httpServer.Listen(port); err != nil {
		a.closeServices()
		log.Fatal("Error starting the server : ", err)
	}

	<-waits
}

func (a *App) setupHTTP() {
	app := web.NewWebserver()

	// setup middlewares
	// setup static files
	// setup routes
	a.httpServer = app
}

func (a *App) setupDatabases() {

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. ")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

    a.client = client
}

func (a *App) setupEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found", err)
	}
}

func (a *App) shutdown() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	signal.Notify(s, syscall.SIGTERM)
	go func() {
		<-s

		log.Println("shutting down...")
		a.closeServices()
		os.Exit(0)
	}()
}

func (a *App) closeServices() {
	if err := a.httpServer.Shutdown(); err != nil {
		log.Fatal("failed to shutdown [httpserver]", err)
	} else {
		log.Println("cleaned up [httpserver]")
	}

	if err := a.client.Disconnect(context.TODO()); err != nil {
        log.Fatal("failed to shutdown [mongodb]", err)
	} else {
        log.Println("cleaned up [mongodb]")
	}
}

func main() {
    NewApp().Run()
}

