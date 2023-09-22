package service

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

func ValidateApiKey(key string, c echo.Context) (bool, error) {
	if key == "" {
		return false, errors.New("Missing API key")
	}
	return key == "valid-key", nil
}

// // // apiKeyIsValid checks if the given API key is valid and returns the principal if it is.
// func apiKeyIsValid(rawKey string, availableKeys map[string][]byte) (string, bool) {
// 	hash := sha256.Sum256([]byte(rawKey))
// 	key := string(hash[:])
// 	reverseKeyIndex := make(map[string]string)
// 	for name,key := availableKeys {
// 		reverseKeyIndex[key] = name
// 	}
// 	name, found := reverseKeyIndex[apiKey]
// 	return name, found
// }

// bearerToken extracts the content from the header, striping the Bearer prefix
func bearerToken(r *http.Request, header string) (string, error) {
	rawToken := r.Header.Get(header)
	pieces := strings.SplitN(rawToken, " ", 2)

	if len(pieces) < 2 {
		return "", errors.New("token with incorrect bearer format")
	}

	token := strings.TrimSpace(pieces[1])

	return token, nil
}
