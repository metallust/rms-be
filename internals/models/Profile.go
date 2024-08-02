package models

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// Profile contains additional information for an applicant's user profile.
type Profile struct {
	Applicant         User   `bson:"applicant" json:"applicantId"`
	ResumeFileAddress string `bson:"resume_file_address" json:"resumeFileAddress"`
	Skills            string `bson:"skills" json:"skills"`
	Education         string `bson:"education" json:"education"`
	Experience        string `bson:"experience" json:"experience"`
	Name              string `bson:"name" json:"name"`
	Email             string `bson:"email" json:"email"`
	Phone             string `bson:"phone" json:"phone"`
}

const COLLECTION_PROFILE = "profiles"

// Insert inserts a new profile into the database.
func (p *Profile) Insert(db *mongo.Database) (*mongo.InsertOneResult, error) {
    collection := db.Collection(COLLECTION_PROFILE)
    return collection.InsertOne(context.Background(), p)
}
