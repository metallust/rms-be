package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
    db *mongo.Database
}

func NewUserController(db *mongo.Database) *User {
    return &User{db: db}
}

func (u *User) UploadResume(c *fiber.Ctx) error {
    return c.SendString("Upload Resume")
}
