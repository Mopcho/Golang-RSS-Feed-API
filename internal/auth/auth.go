package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Example:
// Authroziation: ApiKey{insert_api_ley_here}
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")

	if val == "" {
		return "", errors.New("No authentication key found")
	}

	vals := strings.Split(val, " ")

	if len(vals) != 2 {
		return "", errors.New("Malformed header value")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("Malformed first part of header")
	}

	return vals[1], nil
}
