package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type admin struct {
	db *mongo.Database
}

func NewAdminController(db *mongo.Database) *admin {
	return &admin{
        db: db,
	}
}
func (ad *admin) GetApplicant(c *fiber.Ctx) error {
	return c.SendString("specific applicant")
}

func (ad *admin) Applicants(c *fiber.Ctx) error {
	return c.SendString("specific applicant")
}
