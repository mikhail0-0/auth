package common

import (
	"errors"
	"log"
	"net/http"
)

type RequestError struct {
	StatusCode int
	Err        error
}

func (re RequestError) Error() string {
	log.Println(re.Err.Error())
	return re.Err.Error()
}

var ErrWrongPassword = RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("wrong password"),
}

var ErrAccessInvalidOrExpired = RequestError{
	StatusCode: http.StatusUnauthorized,
	Err:        errors.New("access token is invalid or expired"),
}

var ErrCannotGetAccessToken = RequestError{
	StatusCode: http.StatusUnauthorized,
	Err:        errors.New("cannot get access token"),
}

var ErrWrongAccessToken = RequestError{
	StatusCode: http.StatusUnauthorized,
	Err:        errors.New("access token is wrong for refresh"),
}

var ErrBadJwtPayload = RequestError{
	StatusCode: http.StatusUnauthorized,
	Err:        errors.New("cannot get payload from jwt"),
}

var ErrRefreshTokenExpired = RequestError{
	StatusCode: http.StatusUnauthorized,
	Err:        errors.New("refresh token expired"),
}

var ErrCannotGetRefreshToken = RequestError{
	StatusCode: http.StatusUnauthorized,
	Err:        errors.New("cannot get refresh token"),
}

var ErrWrongRefreshToken = RequestError{
	StatusCode: http.StatusUnauthorized,
	Err:        errors.New("wrong refresh token"),
}

var ErrNotFound = RequestError{
	StatusCode: http.StatusNotFound,
	Err:        errors.New("not found"),
}

var ErrBadRequestFormat = RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("bad request format"),
}

func GetErrorAndStatus(err error) (int, string) {
	switch re := err.(type) {
	case RequestError:
		return re.StatusCode, re.Error()
	default:
		message := "internal server error"
		log.Println(message)
		return http.StatusInternalServerError, message
	}
}
