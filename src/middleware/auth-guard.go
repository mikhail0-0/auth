package middleware

import (
	"auth/src/common"
	"auth/src/token"

	"github.com/gin-gonic/gin"
)

var AuthGuard gin.HandlerFunc = func(c *gin.Context) {
	accessToken, err := getAccessToken(c)
	if err != nil {
		c.AbortWithStatusJSON(common.GetErrorAndStatus(err))
		return
	}

	claims, err := token.Authorize(accessToken)
	if err != nil {
		c.AbortWithStatusJSON(common.GetErrorAndStatus(err))
		return
	}

	for key, val := range *claims {
		c.Set(key, val)
	}
}

func getAccessToken(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	lenBearer := len("Bearer ")
	if len(authHeader) < lenBearer {
		return "", common.ErrCannotGetAccessToken
	}

	return authHeader[lenBearer:], nil
}
