package auth

import "errors"

var AccessTokenExpiredError = errors.New("cannot connect to mongodb")
var InvalidTokenError = errors.New("invalid token")