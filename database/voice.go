package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mihiru-go/models"
)

const collectionNameVoice = "voice"

type VoiceDatabase interface {
	InsertVoice(voice *models.Voice) error
	UpdateVoice(voice *models.VoiceWithObjectId) error
	DeleteVoice(id primitive.ObjectID) error
	GetVoiceById(id primitive.ObjectID) (*models.VoiceWithObjectId, error)
	ListVoiceByLiver(liver string) ([]*models.VoiceWithObjectId, error)
}

func (d *MongoDatabase) InsertVoice(voice *models.Voice) error {
	collection := d.DB.Collection(collectionNameVoice)
	_, err := collection.InsertOne(context.Background(), voice)
	return err
}

func (d *MongoDatabase) UpdateVoice(voice *models.VoiceWithObjectId) error {
	collection := d.DB.Collection(collectionNameVoice)
	_, err := collection.UpdateByID(context.Background(), voice.ID, bson.M{"$set": voice.Voice})
	return err
}

func (d *MongoDatabase) DeleteVoice(id primitive.ObjectID) error {
	collection := d.DB.Collection(collectionNameVoice)
	_, err := collection.UpdateByID(context.Background(), id, bson.M{"$set": bson.D{{Key: "deleted", Value: true}}})
	return err
}

func (d *MongoDatabase) GetVoiceById(id primitive.ObjectID) (*models.VoiceWithObjectId, error) {
	var voice *models.VoiceWithObjectId
	collection := d.DB.Collection(collectionNameVoice)
	err := collection.FindOne(context.Background(), bson.D{{Key: "_id", Value: id}}).Decode(&voice)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	return voice, nil
}

func (d *MongoDatabase) ListVoiceByLiver(liver string) ([]*models.VoiceWithObjectId, error) {
	collection := d.DB.Collection(collectionNameVoice)
	cursor, err := collection.Find(context.Background(),
		bson.D{{Key: "liver", Value: liver}, {Key: "deleted", Value: false}},
		&options.FindOptions{Sort: bson.D{{Key: "sort_no", Value: 1}}},
	)
	if err != nil {
		return nil, err
	}
	defer CloseCursor(cursor, context.Background())
	var data []*models.VoiceWithObjectId
	for cursor.Next(context.Background()) {
		var voice *models.VoiceWithObjectId
		if err = cursor.Decode(&voice); err != nil {
			return nil, err
		}
		data = append(data, voice)
	}

	return data, nil
}
