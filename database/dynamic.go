package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mihiru-go/models"
	"time"
)

const collectionNameDynamic = "dynamic"

type DynamicDatabase interface {
	InsertDynamic(dynamic *models.DynamicWithObjectId) error
	UpdateDynamic(dynamic *models.DynamicWithObjectId) error
	QueryDynamicByTimestamp(startTimestamp int64, endTimeStamp int64) ([]*models.DynamicWithObjectId, error)
	CountDynamicByDay() ([]*models.DayCount, error)
}

func (d *MongoDatabase) InsertDynamic(dynamic *models.DynamicWithObjectId) error {
	collection := d.DB.Collection(collectionNameDynamic)
	insertResult, err := collection.InsertOne(context.Background(), dynamic.DynamicWithLastModified)
	if err != nil {
		return err
	}
	dynamic.ID = insertResult.InsertedID.(primitive.ObjectID)
	return nil
}

func (d *MongoDatabase) UpdateDynamic(dynamic *models.DynamicWithObjectId) error {
	collection := d.DB.Collection(collectionNameDynamic)
	_, err := collection.UpdateByID(context.Background(), dynamic.ID, dynamic)
	return err
}

func (d *MongoDatabase) QueryDynamicByTimestamp(startTimestamp int64, endTimeStamp int64) ([]*models.DynamicWithObjectId, error) {
	collection := d.DB.Collection(collectionNameDynamic)
	cursor, err := collection.Find(context.Background(),
		bson.D{{Key: "timestamp", Value: bson.M{"$gte": startTimestamp}}, {Key: "timestamp", Value: bson.M{"$lt": endTimeStamp}}},
		&options.FindOptions{Sort: bson.D{{Key: "timestamp", Value: 1}}},
	)
	if err != nil {
		return nil, err
	}
	defer CloseCursor(cursor, context.Background())
	var data []*models.DynamicWithObjectId
	for cursor.Next(context.Background()) {
		var dynamic *models.DynamicWithObjectId
		if err = cursor.Decode(&dynamic); err != nil {
			return nil, err
		}
		data = append(data, dynamic)
	}

	return data, nil
}

func (d *MongoDatabase) CountDynamicByDay() ([]*models.DayCount, error) {
	collection := d.DB.Collection(collectionNameDynamic)
	cursor, err := collection.Aggregate(context.Background(), bson.A{
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"$dateToString": bson.M{
						"format": "%Y.%m.%d",
						"date": bson.M{
							"$add": bson.A{primitive.NewDateTimeFromTime(time.Unix(0, 0)), bson.M{
								"$multiply": bson.A{"$timestamp", 1000},
							}},
						},
						"timezone": "+08",
					},
				},
				"count":   bson.M{"$sum": 1},
				"version": bson.M{"$max": "$last_modified"},
			},
		},
		bson.M{"$sort": bson.M{"_id": 1}},
	})
	if err != nil {
		return nil, err
	}
	defer CloseCursor(cursor, context.Background())
	var data []*models.DayCount
	for cursor.Next(context.Background()) {
		var dayCount *models.DayCount
		if err = cursor.Decode(&dayCount); err != nil {
			return nil, err
		}
		data = append(data, dayCount)
	}

	return data, nil
}
