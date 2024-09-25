package token

import (
	"auth/src/common"
	"auth/src/config"
	"auth/src/mail"
	"auth/src/refreshSession"
	"crypto/rand"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessData struct {
	GUID        string `json:"guid"`
	Ip          string `json:"ip"`
	Exp         int64  `json:"exp"`
	AccessToken string `json:"access_token"`
}

type RefreshData struct {
	Exp          int64  `json:"exp"`
	RefreshToken string `json:"refresh_token"`
}

func Authenticate(guid, email, ip string) (*AccessData, *RefreshData, error) {
	targetSession, err := refreshSession.FindByUserId(guid)
	if err == nil {
		refreshSession.Delete(targetSession)
	}

	accessData, refreshData, err := createTokenPair(guid, email, ip)
	if err != nil {
		return nil, nil, err
	}

	return accessData, refreshData, nil
}

func Authorize(accessToken string) (*jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(
		accessToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return config.SecretKey, nil
		},
	)

	if err != nil {
		return nil, common.ErrAccessInvalidOrExpired
	}

	return &claims, nil
}

func RefreshTokens(guid, email, refreshSessionId, refreshToken, ip string) (*AccessData, *RefreshData, error) {
	rs, err := refreshSession.CheckRefresh(guid, refreshToken)
	if err != nil {
		return nil, nil, err
	}

	if rs.ID.String() != refreshSessionId {
		return nil, nil, common.ErrWrongAccessToken
	}

	if rs.Ip != ip {
		err = mail.SendMessage(email, "trying to connect from another ip address")
		log.Println(err)
	}

	refreshSession.Delete(rs)

	accessData, refreshData, err := createTokenPair(guid, email, ip)
	if err != nil {
		return nil, nil, err
	}

	return accessData, refreshData, nil
}

func createTokenPair(guid, email, ip string) (*AccessData, *RefreshData, error) {
	refreshData, err := createRefreshToken()
	if err != nil {
		return nil, nil, err
	}

	rs, err := refreshSession.Create(
		guid,
		refreshData.RefreshToken,
		ip,
		refreshData.Exp,
	)

	if err != nil {
		return nil, nil, err
	}

	accessData, err := createAccessToken(guid, email, rs.ID.String(), ip)
	if err != nil {
		return nil, nil, err
	}

	return &accessData, &refreshData, nil
}

func createAccessToken(guid, email, refreshSessionId, ip string) (AccessData, error) {
	exp := time.Now().Unix() + int64(config.AccessExpiresSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"guid":               guid,
		"email":              email,
		"refresh_session_id": refreshSessionId,
		"ip":                 ip,
		"exp":                exp,
	})
	tokenString, err := token.SignedString(config.SecretKey)
	if err != nil {
		return AccessData{}, err
	}

	return AccessData{
		GUID:        guid,
		Ip:          ip,
		Exp:         exp,
		AccessToken: tokenString,
	}, nil
}

func createRefreshToken() (RefreshData, error) {
	exp := time.Now().Unix() + int64(config.RefreshExpiresSeconds)

	token := make([]byte, 48)
	rand.Read(token)

	return RefreshData{
		Exp:          exp,
		RefreshToken: common.Base64ToString(token),
	}, nil
}
