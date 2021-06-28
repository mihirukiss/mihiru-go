package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"sort"
	"time"
)

type MongoDatabase struct {
	DB            *mongo.Database
	Client        *mongo.Client
	Context       context.Context
	escapeStrings []string
}

func New(uri, username, password, dbname string) (*MongoDatabase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(uri).SetAuth(options.Credential{
		Username: username,
		Password: password,
	})
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(pingCtx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	db := client.Database(dbname)
	escapeStrings := []string{"$", "(", ")", "*", "+", ".", "[", "?", "\\", "^", "{", "|"}
	sort.Strings(escapeStrings)
	return &MongoDatabase{DB: db, Client: client, Context: ctx, escapeStrings: escapeStrings}, nil
}

func (d *MongoDatabase) Close() {
	err := d.Client.Disconnect(d.Context)
	if err != nil {
		log.Println(err.Error())
	}
}

func CloseCursor(cursor *mongo.Cursor, ctx context.Context) {
	err := cursor.Close(ctx)
	if err != nil {
		log.Println(err.Error())
	}
}
