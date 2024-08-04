package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/metallust/rms-be/internals/helper"
	"github.com/metallust/rms-be/internals/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	db *mongo.Database
}

func NewUserController(db *mongo.Database) *User {
	return &User{db: db}
}

func (u *User) UploadResume(c *fiber.Ctx) error {
	//check user == applicant
	if c.Locals("tokendata").(map[string]string)["role"] != "Applicant" {
		log.Error("Invalid role ", c.Locals("tokendata"))
		return c.Status(fiber.StatusUnauthorized).JSON(helper.NewHTTPResponse("Unauthorized access only for applicants", nil))
	}
	//get the file
	file, err := c.FormFile("resume")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helper.NewHTTPResponse("Internal error", nil))
	}

	//create a filepath
	newFileName := fmt.Sprintf("%s_%s", time.Now().Format("20060102150405"), file.Filename)
	destination := fmt.Sprintf("./uploads/%s", newFileName)
	//save the file locally
	if err := c.SaveFile(file, destination); err != nil {
		log.Error("Error saving file", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helper.NewHTTPResponse("Internal error", nil))
	}

	//api call to 3rd party
	f, err := os.Open(destination)
	stat, _ := f.Stat()
	buf := make([]byte, stat.Size())
	if err != nil {
        log.Error("Error opening file", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helper.NewHTTPResponse("Internal error", nil))
	}
	defer f.Close()
	_, err = f.Read(buf)
	if err != nil {
        log.Error("Error reading file", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helper.NewHTTPResponse("Internal error", nil))
	}
	req, err := http.NewRequest("POST", "https://api.apilayer.com/resume_parser/upload", bytes.NewBuffer(buf))
    if err != nil {
        log.Error("Error creating request", err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(helper.NewHTTPResponse("Internal error", nil))
    }
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("apikey", os.Getenv("KEY"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
        log.Error("Error sending request", err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(helper.NewHTTPResponse("Internal error", nil))
	}
	defer resp.Body.Close()

	// Read and print the response from the API
	body, err := io.ReadAll(resp.Body)
	if err != nil {
        log.Error("Error reading response", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helper.NewHTTPResponse("Internal error", nil))
	}

    profile := models.Profile{}
    json.Unmarshal(body, &profile)
    profile.ResumeFileAddress = destination

    //get the applicant id
	//save the profile
    obj, err := primitive.ObjectIDFromHex(c.Locals("tokendata").(map[string]string)["id"])
    if err != nil {
        log.Error("Error converting id", err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(helper.NewHTTPResponse("Internal error", nil))
    }
    profile.Applicant = obj
    profile.Insert(u.db)

    log.Info("Profile created successfully")
	return c.Status(fiber.StatusOK).JSON(helper.NewHTTPResponse("Successfully created your profile", profile))
}
