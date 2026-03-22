package auth

import(
	"strings"
	"errors"
	"net/http"
)

func GetAPIKey(headers http.Header) (string, error) {
	authheader := headers.Get("Authorization")
	if authheader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitHeader := strings.Split(authheader, " ")
	if len(splitHeader) < 2 || splitHeader[0] != "ApiKey" {
		return "", errors.New("Wrong authorization header format")
	}
	
	return splitHeader[1], nil
}