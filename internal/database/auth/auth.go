package auth

import (
	"errors"
	"net/http"
	"strings"
)

//GetAPIKey returns the API key from the request headers
//Authorization : ApiKey {apiKey}
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "",errors.New("no Authentication info found")
	}

	vals := strings.Split(val," ")
	if len(vals) != 2 {
		return "",errors.New("invalid Authentication info")
	}

	if vals[0] != "ApiKey" {
		return "",errors.New("malformed Authentication Header")
	}

	return vals[1],nil
}