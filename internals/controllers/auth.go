package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type Auth struct {
    client *mongo.Client
}

func NewAuthController(client *mongo.Client) *Auth {
    return &Auth{client: client}
}

func (a *Auth)Signup(c *fiber.Ctx) error {
    return c.SendString("Signup")
}

func (a *Auth)Login(c *fiber.Ctx) error {
    return c.SendString("Login")
}
