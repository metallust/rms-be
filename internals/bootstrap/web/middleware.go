package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/metallust/rms-be/internals/helper"
	"go.mongodb.org/mongo-driver/mongo"
)

type Middlewares struct {
	app *fiber.App
	db  *mongo.Database
}

func NewMiddlewares(app *fiber.App, db *mongo.Database) *Middlewares {
	return &Middlewares{app, db}
}

func (m *Middlewares) Init() {
	m.Cors()
	m.app.Use(m.OPT_Auth)
}

func (m *Middlewares) Cors() {
	m.app.Use(cors.New())
}

func (m *Middlewares) OPT_Auth(c *fiber.Ctx) error {

    excludePaths := []string{"/login", "/signup"}
    reqpath := c.Path()
    for _, path := range excludePaths {
        if path == reqpath{
            return c.Next()
        }
    }
    
	token := c.Get("Token", ``)

	if token == `` {
        log.Error("Token not found")
		return c.Status(fiber.StatusBadRequest).JSON(helper.NewHTTPResponse("Token not found", nil))
	}

	tokenData, err := helper.VerifyToken(token)
	if err != nil {
        log.Error("token invalid ", err.Error())
        return c.Status(fiber.StatusUnauthorized).JSON(helper.NewHTTPResponse("Unauthorized", nil))
	}

	c.Locals("tokendata", tokenData)

	return c.Next()
}
