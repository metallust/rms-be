package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// get skills, education, experience, name, email, phone
type Collage struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
type Ex struct {
	Dates []string `json:"dates"`
	Name  string   `json:"name"`
	Url   string   `json:"url"`
}

// Profile contains additional information for an applicant's user profile.
type Profile struct {
	Applicant         primitive.ObjectID  `bson:"_id" json:"applicantId"`
	ResumeFileAddress string `bson:"resume_file_address" json:"resumeFileAddress"`
    Skills     []string  `json:"skills"`
	Education  []Collage `json:"education"`
	Experience []Ex      `json:"experience"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
}

const COLLECTION_PROFILE = "profiles"

// Insert inserts a new profile into the database.
func (p *Profile) Insert(db *mongo.Database) (*mongo.InsertOneResult, error) {
	collection := db.Collection(COLLECTION_PROFILE)
	return collection.InsertOne(context.Background(), p)
}
