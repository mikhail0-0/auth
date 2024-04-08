package config

import (
	"encoding/base64"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var ServerHost   string
var ServerPort   string
var MongoConnStr string
var SecretKey []byte
var RefreshExpiresSeconds int64
var AccessExpiresSeconds int64

func Load(envPath string) error {
	err := godotenv.Load(envPath)
	if err != nil { return err }

	ServerHost = os.Getenv("SERVER_HOST")
	ServerPort = os.Getenv("SERVER_PORT")
	MongoConnStr = os.Getenv("MONGODB_CONN_STR")

	SecretKey, err = base64.RawURLEncoding.DecodeString(
		os.Getenv("SECRET_KEY"),
	)
	if err != nil { return err }

	RefreshExpiresSeconds, err = strconv.ParseInt(
		os.Getenv("REFRESH_EXPIRES_SECONDS"),
		10, 
		64,
	)
	if err != nil { return err }
	
	AccessExpiresSeconds, err = strconv.ParseInt(
		os.Getenv("REFRESH_EXPIRES_SECONDS"),
		10, 
		64,
	)
	if err != nil { return err }

	return nil
}
