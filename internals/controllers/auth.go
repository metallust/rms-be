package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/metallust/rms-be/internals/helper"
	"github.com/metallust/rms-be/internals/models"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	db *mongo.Database
}

func NewAuthController(db *mongo.Database) *Auth {
	return &Auth{db: db}
}

type RegisterIn struct {
	Name            string `json:"name" form:"name" validate:"required,omitempty,min=5"`
	Email           string `json:"email" form:"email" validate:"required,email"`
	Password        string `json:"password" form:"password" validate:"required,min=8"`
	UserType        string `json:"user_type" form:"user_type" validate:"required,oneof=Admin Applicant"`
	ProfileHeadline string `json:"profile_headline" form:"profile_headline"`
	Address         string `json:"address" form:"address"`
}

func (a *Auth) Signup(c *fiber.Ctx) error {
	in, err := helper.ReadBody[RegisterIn](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.NewHTTPResponse("Invalid payload", err.Error()))
	}
	//hash password
	hashedPassword, errGen := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if errGen != nil {
		log.Error("Error generating password", errGen.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helper.NewHTTPResponse("Some Thing went wrong", err.Error()))
	}
	//create user
	user := models.User{}
	user.Name = in.Name
	user.Email = in.Email
	user.Address = in.Address
	user.PasswordHash = string(hashedPassword)
	user.UserType = in.UserType

	if err := user.Save(a.db); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.NewHTTPResponse("User already exist", nil))
	}
	log.Info("User created", user)

	//generate token
	token, err := helper.CreateToken(user.ID.Hex(), user.UserType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helper.NewHTTPResponse("Some Thing went wrong", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(helper.NewHTTPResponse("Successfully added User", map[string]interface{}{"user": user, "token": token}))
}

type LoginIn struct {
	Email    string `json:"email" form:"email" validate:"email,required"`
	Password string `json:"password" form:"password" validate:"required,min=8"`
}

func (a *Auth) Login(c *fiber.Ctx) error {
	//get data
	in, err := helper.ReadBody[LoginIn](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helper.NewHTTPResponse("Invalid payload", err.Error()))
	}

	user := models.User{}
	err = user.Find(a.db, in.Email)

	//check if user exists
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helper.NewHTTPResponse("Internal server error", err.Error()))
	}
	if user.Email == `` {
		return c.Status(fiber.StatusBadRequest).JSON(helper.NewHTTPResponse("Invalid Crediatials", err.Error()))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password))
	if err != nil {
		log.Error("Invalid password ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(helper.NewHTTPResponse("Invalid Creadiantails", nil))
	}

	//generate token
	token, err := helper.CreateToken(user.ID.Hex(), user.UserType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helper.NewHTTPResponse("Some Thing went wrong", err.Error()))
	}

	//return user
	return c.Status(fiber.StatusOK).JSON(helper.NewHTTPResponse("Successfully logged in", map[string]interface{}{"user": user, "token": token}))
}
