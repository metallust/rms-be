package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/metallust/rms-be/internals/helper"
	"github.com/metallust/rms-be/internals/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
func (ad *admin) GetApplicants(c *fiber.Ctx) error {
	//check if the user is admin
    data := c.Locals("tokendata").(map[string]string)
    if role := data["role"]; role != "Admin" {
        log.Error("not an admin", data)
        return c.Status(fiber.StatusUnauthorized).JSON(helper.NewHTTPResponse("Unauthorized", nil))
    }

    collections := ad.db.Collection(models.COLLECTION_USER)
    filter := bson.M{"user_type": "Applicant"}
    results, err := collections.Find(context.Background(), filter)
    if err != nil {
        log.Error("Error fetching all users", err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(helper.NewHTTPResponse("Internal error", nil))
    }
    var allusers []models.User
    err = results.All(context.Background(), &allusers)
    if err != nil {
        log.Error("Error decoding all users", err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(helper.NewHTTPResponse("Internal error", nil))
    }
    log.Info("All users", allusers)
    return c.Status(fiber.StatusOK).JSON(helper.NewHTTPResponse("Successfully fetched all the applicants", allusers))
}

func (ad *admin) GetApplicant(c *fiber.Ctx) error {
	//check if the user is admin
    data := c.Locals("tokendata").(map[string]string)
    if role := data["role"]; role != "Admin" {
        log.Error("not an admin", data)
        return c.Status(fiber.StatusUnauthorized).JSON(helper.NewHTTPResponse("Unauthorized", nil))
    }
    
    collections := ad.db.Collection(models.COLLECTION_USER)
    applicantID, err := primitive.ObjectIDFromHex(c.Params("id"))
    if err != nil {
        log.Error("Error converting id", err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(helper.NewHTTPResponse("Internal error", nil))
    }
    filter := bson.M{"_id": applicantID, "user_type": "Applicant"}
    result := collections.FindOne(context.Background(), filter)
    user := models.User{}
    err = result.Decode(&user)
    if err != nil {
        log.Error("Error decoding user", err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(helper.NewHTTPResponse("Internal error", nil))
    }

    collections = ad.db.Collection(models.COLLECTION_PROFILE)
    filter = bson.M{"_id": applicantID}
    result = collections.FindOne(context.Background(), filter)
    profile := models.Profile{}
    err = result.Decode(&profile)
    if err != nil {
        log.Error("Error decoding user", err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(helper.NewHTTPResponse("Internal error", nil))
    }

    log.Info("User", user, "Profile", profile)
    return c.Status(fiber.StatusOK).JSON(helper.NewHTTPResponse("Successfully fetched the applicant", map[string]interface{}{"user": user, "profile": profile}))
}
