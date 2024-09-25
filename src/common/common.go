package common

import (
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

func GetHash(str string) (string, error) {
	byteStr, err := stringToBase64(str)
	if err != nil {
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword(byteStr, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return Base64ToString(hash), nil
}

func CompareHashAndString(hash, password string) error {
	byteHash, err := stringToBase64(hash)
	if err != nil {
		return err
	}

	byteStr, err := stringToBase64(password)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(
		byteHash,
		byteStr,
	)

	return err
}

func Base64ToString(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

func stringToBase64(str string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(str)
}
