package main

import (
	"auth/src/config"
	"auth/src/postgres"
	"auth/src/router"
	"fmt"
)

func main() {
	err := config.Load(".env")
	if err != nil {
		panic(err)
	}

	_, err = postgres.GetDB()
	if err != nil {
		panic(err)
	}

	router := router.GetRouter()
	router.Run(fmt.Sprintf(":%v", config.ServerPort))
}
