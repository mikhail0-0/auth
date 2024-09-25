package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var ServerPort int

var PgConnStr string
var DbRetries int

var SecretKey []byte
var RefreshExpiresSeconds int
var AccessExpiresSeconds int

var SmtpPassword string
var SmtpHost string
var SmtpPort int
var SmtpEmail string

var errorsArray = make([]error, 0)

func Load(envPath string) error {
	err := godotenv.Load(envPath)
	if err != nil {
		log.Println(".env file not found")
	}

	ServerPort = getEnvInt("SERVER_PORT")

	PgConnStr = fmt.Sprintf(
		"postgresql://%v:%v@%v:%v/%v",
		getEnv("POSTGRES_USERNAME"),
		getEnv("POSTGRES_PASSWORD"),
		getEnv("POSTGRES_HOST"),
		getEnv("POSTGRES_PORT"),
		getEnv("POSTGRES_DATABASE"),
	)

	DbRetries = getEnvInt("DB_RETRIES")

	SecretKey = []byte(getEnv("SECRET_KEY"))

	RefreshExpiresSeconds = getEnvInt("REFRESH_EXPIRES_SECONDS")
	AccessExpiresSeconds = getEnvInt("ACCESS_EXPIRES_SECONDS")

	SmtpEmail = getEnv("SMTP_EMAIL")
	SmtpPassword = getEnv("SMTP_PASSWORD")
	SmtpHost = getEnv("SMTP_HOST")
	SmtpPort = getEnvInt("SMTP_PORT")

	if len(errorsArray) > 0 {
		return errors.Join(errorsArray...)
	}

	return nil
}

func getEnv(envName string) string {
	val, ok := os.LookupEnv(envName)
	if !ok {
		errorsArray = append(errorsArray, fmt.Errorf("env %v was not found", envName))
		return ""
	}
	return val
}

func getEnvInt(envName string) int {
	str := getEnv(envName)
	if str == "" {
		errorsArray = append(errorsArray, fmt.Errorf("env %v was not found", envName))
		return 0
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		errorsArray = append(errorsArray, err)
		return 0
	}
	return val
}
