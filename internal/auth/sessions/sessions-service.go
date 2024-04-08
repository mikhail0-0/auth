package sessionsService

import (
	"context"
	"errors"
	"test-task/internal/auth"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection

func Create (c *mongo.Collection) {
	collection = c
}

func SaveSession(session auth.Session) (*mongo.InsertOneResult, error) {
	return collection.InsertOne(context.TODO(), session)
}

func FindSession(guid string) (auth.Session, error) {
	filter := bson.D{
		{"guid", guid},
		{"exp", bson.D{{"$gt", time.Now().Unix()}}},	
	}

	var session auth.Session
	err := collection.FindOne(context.TODO(), filter).Decode(&session)
	if err != nil {
		return auth.Session{}, errors.New("GUID not found")
	}

	return session, nil
}

func DeleteByGuid(guid string) (*mongo.DeleteResult, error) {
	filter := bson.D{{"guid", guid}}
	return collection.DeleteMany(context.TODO(), filter)
}