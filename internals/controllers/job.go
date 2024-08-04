package controllers

import (
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/metallust/rms-be/internals/helper"
	"github.com/metallust/rms-be/internals/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var COMPANYNAME string = os.Getenv("COMPANY_NAME")

type Jobs struct {
	db *mongo.Database
}

func NewJobsController(db *mongo.Database) *Jobs {
	return &Jobs{db: db}
}

type Jobin struct {
	Title       string `json:"title" form:"tile" validate:"required,omitempty,min=5"`
	Description string `json:"description" form:"description" validate:"required,omitempty,min=5"`
}

func (j *Jobs) CreateJob(c *fiber.Ctx) error {
	data := c.Locals("tokendata").(map[string]string)
	if role := data["role"]; role != "Admin" {
		log.Error("not an admin", data)
		return c.Status(fiber.StatusUnauthorized).JSON(helper.NewHTTPResponse("Unauthorized", nil))
	}
	in, err := helper.ReadBody[Jobin](c)
	if err != nil {
		log.Error("Error reading body", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(helper.NewHTTPResponse("Bad request", err.Error()))
	}

	postedby, err := primitive.ObjectIDFromHex(data["id"])
	if err != nil {
		log.Error("Error converting id", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(helper.NewHTTPResponse("Bad request", err.Error()))
	}
	job := models.Job{
		ID:                [12]byte{},
		Title:             in.Title,
		Description:       in.Description,
		PostedOn:          time.Now(),
		TotalApplications: 0,
		CompanyName:       COMPANYNAME,
		PostedBy:          postedby,
	}
	job.Insert(j.db)

	log.Info("Job created successfully", job)
	return c.Status(fiber.StatusOK).JSON(helper.NewHTTPResponse("Job created successfully", job))
}

type Jobout struct {
	ID          string `json:"id" form:"id" validate:"required"`
	Title       string `json:"title" form:"tile" validate:"required,omitempty,min=5"`
	Description string `json:"description" form:"description" validate:"required,omitempty,min=5"`
	CompanyName string `json:"company_name" form:"company_name" validate:"required"`
	PostedBy    string `json:"posted_by" form:"posted_by" validate:"required"`
}

func (j *Jobs) GetJobs(c *fiber.Ctx) error {
	collections := j.db.Collection(models.COLLECTION_JOB)
	results, err := collections.Find(context.Background(), bson.M{})
	if err != nil {
		log.Error("Error fetching all job", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helper.NewHTTPResponse("Internal error", nil))
	}
	var jobs []models.Job
	err = results.All(context.Background(), &jobs)
	if err != nil {
		log.Error("Error decoding all users", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helper.NewHTTPResponse("Internal error", nil))
	}

	jobsout := make([]Jobout, len(jobs))
	for i, job := range jobs {
		jobsout[i] = Jobout{
			ID:          job.ID.Hex(),
			Title:       job.Title,
			Description: job.Description,
			CompanyName: job.CompanyName,
			PostedBy:    job.PostedBy.Hex(),
		}
	}

	return c.Status(fiber.StatusOK).JSON(helper.NewHTTPResponse("Successfully fetched all the jobs", jobsout))
}

func (j *Jobs) ApplyJob(c *fiber.Ctx) error {
	data := c.Locals("tokendata").(map[string]string)
	if role := data["role"]; role != "Applicant" {
		log.Error("not an admin", data)
		return c.Status(fiber.StatusUnauthorized).JSON(helper.NewHTTPResponse("Unauthorized", nil))
	}

	job_id := c.Query("job_id")
	jobID, err := primitive.ObjectIDFromHex(job_id)
	if err != nil {
		log.Error("Error converting id", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(helper.NewHTTPResponse("Bad request", err.Error()))
	}

	applicantID, err := primitive.ObjectIDFromHex(data["id"])
	if err != nil {
		log.Error("Error converting id", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(helper.NewHTTPResponse("Bad request", err.Error()))
	}

	collections := j.db.Collection(models.COLLECTION_JOB)
	filter := bson.M{"_id": jobID}
	result := collections.FindOne(context.Background(), filter)
	job := models.Job{}
	err = result.Decode(&job)
	if err != nil {
		log.Error("Error decoding job", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helper.NewHTTPResponse("Internal error", nil))
	}

	for _, application := range job.Applications {
		if application == applicantID {
			return c.Status(fiber.StatusBadRequest).JSON(helper.NewHTTPResponse("You have already applied for this job", nil))
		}
	}

	job.Applications = append(job.Applications, applicantID)
	job.TotalApplications += 1
	_, err = collections.ReplaceOne(context.Background(), filter, job)
	if err != nil {
		log.Error("Error updating job", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helper.NewHTTPResponse("Internal error", nil))
	}

	log.Info("Job updated successfully", job)
	return c.Status(fiber.StatusOK).JSON(helper.NewHTTPResponse("Successfully applied for the job", nil))
}

func (j *Jobs) GetJob(c *fiber.Ctx) error {
	data := c.Locals("tokendata").(map[string]string)
	if role := data["role"]; role != "Admin" {
		log.Error("not an admin", data)
		return c.Status(fiber.StatusUnauthorized).JSON(helper.NewHTTPResponse("Unauthorized", nil))
	}

	job_id := c.Params("id")
    jobID, err := primitive.ObjectIDFromHex(job_id)
    if err != nil {
        log.Error("Error converting id", err.Error())
        return c.Status(fiber.StatusBadRequest).JSON(helper.NewHTTPResponse("Bad request", err.Error()))
    }

	collections := j.db.Collection(models.COLLECTION_JOB)
	results := collections.FindOne(context.Background(), bson.M{"_id": jobID})

	var job models.Job
	err = results.Decode(&job)
	if err != nil {
		log.Error("Error decoding all users", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helper.NewHTTPResponse("Internal error", nil))
	}
	log.Info("job fetched successfully", job)
	return c.Status(fiber.StatusOK).JSON(helper.NewHTTPResponse("Successfully fetched job", job))
}
