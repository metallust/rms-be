package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Job represents a job opening created by an admin user.
type Job struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title             string             `bson:"title" json:"title"`
	Description       string             `bson:"description" json:"description"`
	PostedOn          time.Time          `bson:"posted_on" json:"postedOn"`
	TotalApplications int                `bson:"total_applications" json:"totalApplications"`
	CompanyName       string             `bson:"company_name" json:"companyName"`
	PostedBy          primitive.ObjectID `bson:"posted_by" json:"postedBy"`
}

const COLLECTION_JOB = "jobs"


func (j *Job) Insert(db *mongo.Database)  (*mongo.InsertOneResult, error) {
    collection := db.Collection(COLLECTION_JOB)
    return collection.InsertOne(context.Background(), j)
}
