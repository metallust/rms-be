package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	// "github.com/metallust/rms-be/internals/bootstrap/web"
	"github.com/metallust/rms-be/internals/bootstrap/web"
	"github.com/metallust/rms-be/internals/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	// Common
	httpServer *fiber.App
	db     *mongo.Database
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

	port := ":" + os.Getenv("PORT")
	if port == ":" {
		log.Fatal("Set your 'PORT' environment variable. ")
	}

	if err := a.httpServer.Listen(port); err != nil {
		a.closeServices()
		log.Fatal("Error starting the server : ", err)
	}

	<-waits
}

func (a *App) setupHTTP() {
    a.httpServer = fiber.New()

	// setup middlewares
    middleware := web.NewMiddlewares(a.httpServer, a.db)
    middleware.Init()
	// setup static files
	// setup routes

	authcontroller := controllers.NewAuthController(a.db)
	jobscontroller := controllers.NewJobsController(a.db)
	admincontroller := controllers.NewAdminController(a.db)
	usercontroller := controllers.NewUserController(a.db)

	//Create a profile on the system (Name, Email, Password, UserType (Admin/Applicant), Profile Headline, Address).
	a.httpServer.Post("/signup", authcontroller.Signup)
	//Authenticate users and return a JWT token upon successful validation.
	a.httpServer.Post("/login", authcontroller.Login)

	//Authenticated API for uploading resume files (only PDF or DOCX) of the applicant. Only Applicant type users can access this API.
	a.httpServer.Post("/uploadResume", usercontroller.UploadResume)
	a.httpServer.Route("/admin", func(admin fiber.Router) {
		//Authenticated API for creating job openings. Only Admin type users can access this API.
		admin.Get("/job", jobscontroller.CreateJob)
		//Authenticated API for fetching information regarding a job opening.
		//Returns details about the job opening and a list of applicants. Only Admin type users can access this API.
		admin.Post("/job/:job_id", jobscontroller.GetJob)
		//Authenticated API for fetching a list of all users in the system. Only Admin type users can access this API
		admin.Get("/applicants", admincontroller.GetApplicants)
		//Authenticated API for fetching extracted data of an applicant. Only Admin type users can access this API.
        admin.Get("/applicant/:id", admincontroller.GetApplicant)
	})

	a.httpServer.Route("/jobs", func(router fiber.Router) {
		//Authenticated API for fetching job openings. All users can access this API.
		router.Post("/", jobscontroller.GetJobs)
		//Authenticated API for applying to a particular job. Only Applicant users are allowed to apply for jobs.
		router.Post("/apply", jobscontroller.ApplyJob)
	})

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

	a.db = client.Database(os.Getenv("DB"))
}

func (a *App) setupEnv() {
	if err := godotenv.Load(".env"); err != nil {
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

	if err := a.db.Client().Disconnect(context.TODO()); err != nil {
		log.Fatal("failed to shutdown [mongodb]", err)
	} else {
		log.Println("cleaned up [mongodb]")
	}
}

func main() {
	NewApp().Run()
}
