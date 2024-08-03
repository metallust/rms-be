package web

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

type Webserver struct {
	fiber.Config
}

func NewWebserver() *fiber.App {
	return fiber.New(fiber.Config{
		JSONEncoder:             json.Marshal,
		JSONDecoder:             json.Unmarshal,
	})
}
