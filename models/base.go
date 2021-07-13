package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ObjectIdFields struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`
}

type LastModifiedFields struct {
	LastModified int64 `bson:"last_modified" json:"last_modified"`
}

type DayCount struct {
	Day     string `bson:"_id" json:"day"`
	Count   int64  `bson:"count" json:"count"`
	Version int64  `bson:"version" json:"version"`
}
