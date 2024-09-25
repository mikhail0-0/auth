package test

import (
	"auth/src/common"
	"auth/src/config"
	"auth/src/controller"
	"auth/src/refreshSession"
	"auth/src/router"
	"auth/src/token"
	"auth/src/user"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var gormDB *gorm.DB

var accessToken string = ""
var refreshToken string

func TestMain(m *testing.M) {
	createTestDb()

	setEnv()

	runRouter()

	code := m.Run()
	os.Exit(code)
}

func TestAuthentification(t *testing.T) {
	url := fmt.Sprintf("http://localhost:%v/auth", config.ServerPort)
	client := resty.New()

	resp, _ := client.R().Post(url)

	checkError(t, resp, common.ErrBadRequestFormat)

	resp, _ = client.R().
		SetBody(controller.AuthRequest{
			GUID:     uuid.New().String(),
			Password: "test",
		}).
		Post(url)

	checkError(t, resp, common.ErrNotFound)

	var usr user.User
	result := gormDB.First(&usr)
	if result.Error != nil {
		t.Error("test user not found")
	}

	resp, _ = client.R().
		SetBody(controller.AuthRequest{
			GUID:     usr.ID.String(),
			Password: "test1",
		}).
		Post(url)

	checkError(t, resp, common.ErrWrongPassword)

	client.R().
		SetBody(controller.AuthRequest{
			GUID:     usr.ID.String(),
			Password: "test",
		}).
		Post(url)

	resp, _ = client.R().
		SetBody(controller.AuthRequest{
			GUID:     usr.ID.String(),
			Password: "test",
		}).
		Post(url)

	var refreshSessions []refreshSession.RefreshSession
	gormDB.Find(&refreshSessions)
	if len(refreshSessions) != 1 {
		t.Error("old refresh session not deleted")
	}

	check(t, resp.StatusCode(), http.StatusCreated, "wrong status")

	var respBody token.AccessData
	err := json.Unmarshal(resp.Body(), &respBody)
	if err != nil {
		t.Error("wrong response body")
	}

	accessToken = respBody.AccessToken

	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == config.REFRESH_TOKEN_COOKIE {
			refreshToken = cookie.Value
		}
	}
	if refreshToken == "" {
		t.Error("refresh hash cookie is not set")
	}
}

func TestRefresh(t *testing.T) {
	url := fmt.Sprintf("http://localhost:%v/auth/refresh", config.ServerPort)

	client := resty.New()
	resp, _ := client.R().Post(url)
	checkError(t, resp, common.ErrCannotGetAccessToken)

	resp, _ = client.R().
		SetAuthToken(accessToken).
		Post(url)

	checkError(t, resp, common.ErrCannotGetRefreshToken)

	randomToken := make([]byte, 48)
	rand.Read(randomToken)

	stringToken := base64.RawURLEncoding.EncodeToString(randomToken)

	resp, _ = client.R().
		SetAuthToken(accessToken).
		SetCookie(&http.Cookie{
			Name:  config.REFRESH_TOKEN_COOKIE,
			Value: stringToken,
		}).
		Post(url)

	checkError(t, resp, common.ErrWrongRefreshToken)

	config.AccessExpiresSeconds = 0

	resp, _ = client.R().
		SetAuthToken(accessToken).
		SetCookie(&http.Cookie{
			Name:  config.REFRESH_TOKEN_COOKIE,
			Value: refreshToken,
		}).
		Post(url)

	check(t, resp.StatusCode(), http.StatusCreated, "wrong status")

	var respBody token.AccessData
	err := json.Unmarshal(resp.Body(), &respBody)
	if err != nil {
		t.Error("wrong response body")
	}

	resp, _ = client.R().
		SetAuthToken(accessToken).
		Post(url)

	checkError(t, resp, common.ErrWrongAccessToken)

	resp, _ = client.R().
		SetAuthToken(respBody.AccessToken).
		Post(url)

	checkError(t, resp, common.ErrAccessInvalidOrExpired)
}

func TestProtected(t *testing.T) {
	url := fmt.Sprintf("http://localhost:%v/protected", config.ServerPort)
	client := resty.New()

	resp, _ := client.R().Get(url)

	checkError(t, resp, common.ErrCannotGetAccessToken)

	resp, _ = client.R().
		SetAuthToken(accessToken).
		Get(url)

	check(t, resp.StatusCode(), http.StatusOK, "wrong status")
	check(t, string(resp.Body()), "\""+config.PROTECTED_DATA+"\"", "wrong data")
}

func check(t *testing.T, received, expected any, messageOnFail string) {
	if received != expected {
		t.Errorf("%v %v, got %v instead", messageOnFail, expected, received)
	}
}
func checkError(t *testing.T, resp *resty.Response, err common.RequestError) {
	check(t, resp.StatusCode(), err.StatusCode, "Unexpected status code")
	check(t, string(resp.Body()), "\""+err.Error()+"\"", "Unexpected body")
}

func createTestDb() {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	pg, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15",
		Env: []string{
			"POSTGRES_DB=example",
			"POSTGRES_HOST_AUTH_METHOD=trust",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	pg.Expire(30)

	postgresPort := pg.GetPort("5432/tcp")
	os.Setenv("POSTGRES_PORT", postgresPort)

	if err := pool.Retry(func() error {
		var connErr error
		gormDB, connErr = gorm.Open(postgres.Open(
			fmt.Sprintf("postgresql://postgres@localhost:%s/example", postgresPort)),
			&gorm.Config{},
		)
		if connErr != nil {
			return connErr
		}

		return nil
	}); err != nil {
		panic("Could not connect to postgres: " + err.Error())
	}

	user.Init(gormDB)
	refreshSession.Init(gormDB)
}

func runRouter() {
	err := config.Load(".test.env")
	if err != nil {
		panic(err)
	}

	router := router.GetRouter()
	go router.Run(fmt.Sprintf(":%v", config.ServerPort))
}

func setEnv() {
	os.Setenv("SERVER_PORT", "3000")

	os.Setenv("POSTGRES_HOST", "")
	os.Setenv("POSTGRES_PORT", "")
	os.Setenv("POSTGRES_DATABASE", "")
	os.Setenv("POSTGRES_USERNAME", "")
	os.Setenv("POSTGRES_PASSWORD", "")

	os.Setenv("DB_RETRIES", "0")

	os.Setenv(
		"SECRET_KEY",
		"MIqSERB0Zoww6bhTj2Wzjn1Znuqwmj3/m0ktqfN3Zz/mcMOMMtbU9zegSSepbxCcb1qbVcnWXlqVYO88",
	)
	os.Setenv("ACCESS_EXPIRES_SECONDS", "100")
	os.Setenv("REFRESH_EXPIRES_SECONDS", "100")

	os.Setenv("SMTP_EMAIL", "")
	os.Setenv("SMTP_PASSWORD", "")
	os.Setenv("SMTP_HOST", "")
	os.Setenv("SMTP_PORT", "0")

}
