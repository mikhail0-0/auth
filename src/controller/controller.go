package controller

import (
	"auth/src/common"
	"auth/src/config"
	"auth/src/token"
	"auth/src/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthRequest struct {
	GUID     string `json:"id"`
	Password string `json:"password"`
}

var Authenticate gin.HandlerFunc = func(c *gin.Context) {
	var dto AuthRequest
	if err := c.BindJSON(&dto); err != nil {
		c.IndentedJSON(common.GetErrorAndStatus(common.ErrBadRequestFormat))
		return
	}

	usr, err := user.Verify(dto.GUID, dto.Password)
	if err != nil {
		c.IndentedJSON(common.GetErrorAndStatus(err))
		return
	}

	ip := c.ClientIP()
	accessData, refreshData, err := token.Authenticate(dto.GUID, usr.Email, ip)
	if err != nil {
		c.IndentedJSON(common.GetErrorAndStatus(err))
		return
	}

	setRefreshCookie(c, refreshData)

	c.IndentedJSON(http.StatusCreated, accessData)
}

var Refresh gin.HandlerFunc = func(c *gin.Context) {
	refreshToken, err := c.Cookie(config.REFRESH_TOKEN_COOKIE)
	if err != nil {
		c.IndentedJSON(common.GetErrorAndStatus(common.ErrCannotGetRefreshToken))
		return
	}

	guid, ok := c.Get("guid")
	if !ok {
		c.IndentedJSON(common.GetErrorAndStatus(common.ErrBadJwtPayload))
		return
	}

	email, ok := c.Get("email")
	if !ok {
		c.IndentedJSON(common.GetErrorAndStatus(common.ErrBadJwtPayload))
		return
	}

	refreshSessionId, ok := c.Get("refresh_session_id")
	if !ok {
		c.IndentedJSON(common.GetErrorAndStatus(common.ErrBadJwtPayload))
		return
	}

	ip := c.ClientIP()

	accessData, refreshData, err := token.RefreshTokens(
		fmt.Sprintf("%v", guid),
		fmt.Sprintf("%v", email),
		fmt.Sprintf("%v", refreshSessionId),
		refreshToken,
		ip,
	)
	if err != nil {
		c.IndentedJSON(common.GetErrorAndStatus(err))
		return
	}

	setRefreshCookie(c, refreshData)
	c.IndentedJSON(http.StatusCreated, accessData)
}

var GetProtected gin.HandlerFunc = func(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, config.PROTECTED_DATA)
}

func setRefreshCookie(c *gin.Context, refreshData *token.RefreshData) {
	c.SetCookie(
		config.REFRESH_TOKEN_COOKIE,
		refreshData.RefreshToken,
		config.RefreshExpiresSeconds,
		config.AUTH_PATH,
		c.Request.Host,
		false,
		true,
	)
}
