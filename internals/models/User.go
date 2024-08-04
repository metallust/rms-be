package models

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// User represents a user in the system (Admin or Applicant).
type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name            string             `bson:"name" json:"name"`
	Email           string             `bson:"email" json:"email"`
	Address         string             `bson:"address" json:"address"`
	UserType        string             `bson:"user_type" json:"userType"` // "Applicant" or "Admin"
	PasswordHash    string             `bson:"password_hash" json:"-"`
	ProfileHeadline string             `bson:"profile_headline" json:"profileHeadline"`
	CreatedAt       time.Time          `bson:"created_at" json:"createdAt"`
	// UpdatedAt       time.Time          `bson:"updated_at" json:"updatedAt"`
}

const COLLECTION_USER = "users"

func (u *User) Save(db *mongo.Database) error {

	//check email exists
	collection := db.Collection(COLLECTION_USER)
	filter := bson.M{"email": u.Email}

	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		log.Error("Error checking email", err.Error())
		return err
	}
	if count > 0 {
		log.Error("Email already exists")
		return errors.New("Email already exists")
	}
	result, err := collection.InsertOne(context.Background(), u)
	if err != nil {
		log.Error("Error saving user", err.Error())
		return err
	}
	u.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (u *User) Find(db *mongo.Database, email string) (err error) {
	collection := db.Collection(COLLECTION_USER)

	filter := bson.M{"email": email}
	result := collection.FindOne(context.Background(), filter)
	err = result.Decode(u)
	if err != nil {
		log.Error("can't decode", err.Error())
		return
	}
	return
}
