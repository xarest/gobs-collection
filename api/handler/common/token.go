package common

import (
	"errors"
	"strings"
)

const BEARER_AUTH = "Bearer"

func ExtractToken(tokenString string) (token string, err error) {
	if tokenString == "" {
		return "", errors.New("missing token")
	}
	bearerIdx := strings.Index(tokenString, BEARER_AUTH)
	if bearerIdx >= 0 {
		tokenString = tokenString[bearerIdx+len(BEARER_AUTH)+1:]
	}
	return tokenString, nil
}
