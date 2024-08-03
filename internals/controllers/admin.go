package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type admin struct {
	client *mongo.Client
}

func NewAdminController(client *mongo.Client) *admin {
	return &admin{
		client: client,
	}
}
func (ad *admin) GetApplicant(c *fiber.Ctx) error {
	return c.SendString("specific applicant")
}

func (ad *admin) Applicants(c *fiber.Ctx) error {
	return c.SendString("specific applicant")
}
