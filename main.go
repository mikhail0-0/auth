package main

import (
	"context"
	"fmt"
	"test-task/config"
	"test-task/internal/auth"
	sessionsService "test-task/internal/auth/sessions"
	usersService "test-task/internal/auth/users"
	"test-task/internal/router"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



func createCollections(connStr string) error {
	clientOpts := options.Client().ApplyURI(connStr)

	secondCTx, _ := context.WithTimeout(context.TODO(), time.Second)

	client, err := mongo.Connect(secondCTx, clientOpts)
	if err != nil {
		return err
	}

	err = client.Ping(secondCTx, nil)
	if err != nil {
		return err
	}

	database := client.Database("auth")

	usersService.Create(database.Collection("users"))
	sessionsService.Create(database.Collection("sessions"))

	return nil
}

func fakeMigration() error {
	guid := uuid.New().String()
	_, err := usersService.SaveUser(auth.User{GUID: guid})
	if err != nil { return err }

	fmt.Printf("Users GUID: %v", guid)
	fmt.Println()
	return nil
}

func main() {
	err := config.Load(".env")
	if err != nil {
		fmt.Println("bad environment")
	}

	err = createCollections(config.MongoConnStr)
	if err != nil {
		fmt.Println("cannot connect to mongodb")
	}

	err = fakeMigration()
	if err != nil {
		fmt.Println("cannot connect to mongodb")
	}

	err = router.Run()
	if err != nil {
		fmt.Println("cannot start server")
	}

}
