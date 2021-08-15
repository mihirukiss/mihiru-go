package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mihiru-go/models"
	"time"
)

const collectionNameLive = "live"

type LiveDatabase interface {
	InsertLive(live *models.LiveWithObjectId) error
	UpdateLive(live *models.LiveWithObjectId) error
	QueryLiveByTimestamp(startTimestamp int64, endTimeStamp int64) ([]*models.LiveWithObjectId, error)
	CountLiveByDay() ([]*models.DayCount, error)
}

func (d *MongoDatabase) InsertLive(live *models.LiveWithObjectId) error {
	collection := d.DB.Collection(collectionNameLive)
	insertResult, err := collection.InsertOne(context.Background(), live.LiveWithLastModified)
	if err != nil {
		return err
	}
	live.ID = insertResult.InsertedID.(primitive.ObjectID)
	return nil
}

func (d *MongoDatabase) UpdateLive(live *models.LiveWithObjectId) error {
	collection := d.DB.Collection(collectionNameLive)
	_, err := collection.UpdateByID(context.Background(), live.ID, bson.M{"$set": live.LiveWithLastModified})
	return err
}

func (d *MongoDatabase) QueryLiveByTimestamp(startTimestamp int64, endTimeStamp int64) ([]*models.LiveWithObjectId, error) {
	collection := d.DB.Collection(collectionNameLive)
	cursor, err := collection.Find(context.Background(),
		bson.D{{Key: "timestamp", Value: bson.M{"$gte": startTimestamp}}, {Key: "timestamp", Value: bson.M{"$lt": endTimeStamp}}},
		&options.FindOptions{Sort: bson.D{{Key: "timestamp", Value: 1}}},
	)
	if err != nil {
		return nil, err
	}
	defer CloseCursor(cursor, context.Background())
	var data []*models.LiveWithObjectId
	for cursor.Next(context.Background()) {
		var live *models.LiveWithObjectId
		if err = cursor.Decode(&live); err != nil {
			return nil, err
		}
		data = append(data, live)
	}

	return data, nil
}

func (d *MongoDatabase) CountLiveByDay() ([]*models.DayCount, error) {
	collection := d.DB.Collection(collectionNameLive)
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
