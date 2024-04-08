package usersService

import (
	"context"
	"errors"
	"test-task/internal/auth"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection

func Create (c *mongo.Collection) {
	collection = c
}

func SaveUser(u auth.User) (*mongo.InsertOneResult, error) {
	return collection.InsertOne(context.TODO(), u)
}

func FindByGuid(guid string) (auth.User, error) {
	filter := bson.D{{"guid", guid}}

	var result auth.User
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return auth.User{}, errors.New("GUID not found")
	}

	return result, nil
}