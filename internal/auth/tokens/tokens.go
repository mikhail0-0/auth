package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"test-task/config"
	"test-task/internal/auth"
	sessionsService "test-task/internal/auth/sessions"
	usersService "test-task/internal/auth/users"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func createRefreshToken(guid string) (string, string, error) {
	refreshToken := make([]byte, 64)
	rand.Read(refreshToken)

	hash, err := getHash(refreshToken)

	if err != nil {	return "", "", err }

	refreshString := base64ToString(refreshToken)

	return refreshString, hash, nil	
}

func createAccessToken(guid string) (auth.AccessResponse, error) {
	exp := time.Now().Unix() + config.AccessExpiresSeconds

	token := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.MapClaims{
			"guid": guid,
			"exp":  exp,
		})

	tokenString, err := token.SignedString(config.SecretKey)
	if err != nil {	return auth.AccessResponse{}, err }

	return auth.AccessResponse{
		GUID: guid,
		Exp: exp,
		AccessToken: tokenString,
	}, nil
}

func base64ToString(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

func stringToBase64(str string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(str)
}

func getHash(b []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(hash), nil
}

func getTokensJson(tokens auth.TokensData) ([]byte, error) {
	json, err := json.Marshal(tokens)

	if err != nil {	return []byte{}, err }

	return json, nil
}

func verifyRefreshToken(hash, token string) error {
	byteHash, err := stringToBase64(hash)
	if err != nil {	return err }

	byteToken, err := stringToBase64(token)
	if err != nil { return	err }

	err = bcrypt.CompareHashAndPassword(byteHash, byteToken)
	if err != nil { return	err }

	return nil
}

func verifyAccessToken(hash, token string) error {
	byteHash, err := stringToBase64(hash)
	if err != nil {	return err }

	byteToken := []byte(token)[(len(token) - 64):]

	err = bcrypt.CompareHashAndPassword(byteHash, byteToken)
	if err != nil {	return err }

	return nil
}

func checkRefresh(a auth.RefreshData) error {
	session, err := sessionsService.FindSession(a.GUID)
	if err != nil {	return err }

	err = verifyRefreshToken(session.RefreshHash, a.RefreshToken)
	if err != nil {	return err }

	return nil
}

func CheckGuid(guid string) error {
	_, err := usersService.FindByGuid(guid)
	return err
}

func CreateTokenPair(guid string) ([]byte, error) {
	refreshToken, refreshHash, err := createRefreshToken(guid)
	if err != nil { return []byte{}, err }

	accessResponse, err := createAccessToken(guid)
	if err != nil { return []byte{}, err }

	access := accessResponse.AccessToken
	accessByte := []byte(access)[(len(access) - 64):]

	accessHash, err := getHash(accessByte)
	if err != nil { return []byte{}, err }

	refreshExp := time.Now().Unix() + config.AccessExpiresSeconds 

	_, err = sessionsService.DeleteByGuid(guid)
	if err != nil { return []byte{}, err }

	sessionsService.SaveSession(auth.Session{
		GUID:         guid,
		RefreshHash:  refreshHash,
		AccessHash:   accessHash,
		Exp:          refreshExp,
	})

	tokensData := auth.TokensData{
		GUID:         accessResponse.GUID,
		Exp:          accessResponse.Exp,
		AccessToken:  accessResponse.AccessToken,
		RefreshToken: refreshToken,
	}

	tokensJson, err := getTokensJson(tokensData)
	if err != nil { return []byte{}, err }

	return tokensJson, nil
}

func VerifyAccess(accessPayload auth.AccessPayload, accessToken string) error {
	session, err := sessionsService.FindSession(accessPayload.GUID)
	if err != nil {	return err }

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(
		accessToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.SecretKey), nil
		},
	)

	if err != nil {	return err }

	if !token.Valid { return auth.InvalidTokenError  }

	var exp int64 = 0
	for key, val := range claims {
		if key == "exp" {
			exp = int64(val.(float64))
		}
	}

	if exp == 0 { return auth.InvalidTokenError }
	if time.Now().Unix() > exp { return auth.AccessTokenExpiredError }

	err = verifyAccessToken(session.AccessHash, accessToken)
	if err != nil { return auth.InvalidTokenError }

	return nil
}

func VerifyRefresh(refreshData auth.RefreshData) ([]byte , error) {

	err := checkRefresh(refreshData)
	if err != nil {	return []byte{}, err }

	return CreateTokenPair(refreshData.GUID)
}

func CompareHashAndValue(hash, value string) error {
	byteHash, err := stringToBase64(hash)
	if err != nil { return err }

	byteValue, err := stringToBase64(value)
	if err != nil {	return err }

	err = bcrypt.CompareHashAndPassword(byteHash, byteValue)
	if err != nil { return err }

	return nil
}