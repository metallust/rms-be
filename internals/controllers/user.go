package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
    client *mongo.Client
}

func NewUserController(client *mongo.Client) *User {
    return &User{client: client}
}

func (u *User) UploadResume(c *fiber.Ctx) error {
    return c.SendString("Upload Resume")
}
