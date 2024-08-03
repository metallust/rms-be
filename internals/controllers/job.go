package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type Jobs struct {
    db *mongo.Database
}

func NewJobsController(db *mongo.Database) *Jobs {
    return &Jobs{db: db}
}

func (j *Jobs) GetJobs(c *fiber.Ctx) error {
    return c.SendString("GetJobs")
}

func (j *Jobs) GetJob(c *fiber.Ctx) error {
    return c.SendString("GetJob")
}

func (j *Jobs) CreateJob(c *fiber.Ctx) error {
    return c.SendString("CreateJob")
}

func (j *Jobs) ApplyJob(c *fiber.Ctx) error {
    return c.SendString("ApplyJob")
}
