package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"mihiru-go/models"
)

const collectionNameUser = "user"

type UserDatabase interface {
	InsertUser(user *models.UserWithObjectId) error
	UpdateUser(user *models.UserWithObjectId) error
	GetUserByLoginName(loginName string) (*models.UserWithObjectId, error)
	GetUserById(id primitive.ObjectID) (*models.UserWithObjectId, error)
}

func (d *MongoDatabase) InsertUser(user *models.UserWithObjectId) error {
	collection := d.DB.Collection(collectionNameUser)
	insertResult, err := collection.InsertOne(context.Background(), user.User)
	if err != nil {
		return err
	}
	user.ID = insertResult.InsertedID.(primitive.ObjectID)
	return nil
}

func (d *MongoDatabase) UpdateUser(user *models.UserWithObjectId) error {
	collection := d.DB.Collection(collectionNameUser)
	_, err := collection.UpdateByID(context.Background(), user.ID, user)
	return err
}

func (d *MongoDatabase) GetUserByLoginName(loginName string) (*models.UserWithObjectId, error) {
	var user *models.UserWithObjectId
	collection := d.DB.Collection(collectionNameUser)
	err := collection.FindOne(context.Background(), bson.D{{Key: "loginName", Value: loginName}}).Decode(&user)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	return user, nil
}

func (d *MongoDatabase) GetUserById(id primitive.ObjectID) (*models.UserWithObjectId, error) {
	var user *models.UserWithObjectId
	collection := d.DB.Collection(collectionNameUser)
	err := collection.FindOne(context.Background(), bson.D{{Key: "_id", Value: id}}).Decode(&user)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	return user, nil
}
